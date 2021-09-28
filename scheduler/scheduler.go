package scheduler

import (
    "regexp"
    "sync"
    "time"
)

import (
    "github.com/baidu/go-lib/log"
    "github.com/baidu/go-lib/queue"
)

import (
    "github.com/coopersong/mini-spider/loader"
    "github.com/coopersong/mini-spider/parser"
)

type Scheduler struct {
    // 任务队列
    TaskQue         queue.Queue
    // url去重表
    UrlTable        sync.Map
    // 任务channel
    TaskChan        chan struct{}
    // 最大爬取深度
    MaxDepth        int
    // 爬取间隔 单位秒
    CrawlInterval   int
    // 爬取任务所使用的go routine数
    ThreadCount     int
    // 任务通用配置
    TaskCommonCfg   *TaskCommonConfig
    // 站点爬取间隔timer表
    TimerTable      sync.Map
}

// Create a new scheduler.
func NewScheduler() *Scheduler {
    return new(Scheduler)
}

// Initialize scheduler by config.
// Initialize scheduler's task queue by seeds.
func (s *Scheduler) Init(config loader.Config, seeds []string) {
    // we have checked TargetUrl in config's check,
    // so we do not need to check again,
    // just ignore the possible error
    targetUrlPattern, _ := regexp.Compile(config.TargetUrl)
    taskCommonCfg := &TaskCommonConfig{
        CrawlTimeout: config.CrawlTimeout,
        OutputDirectory: config.OutputDirectory,
        TargetUrlPattern: targetUrlPattern,
    }

    // initialize task queue
    s.TaskQue.Init()
    for _, seed := range seeds {
        task := &Task{
            Url:       seed,
            Depth:     0,
            CommonCfg: taskCommonCfg,
        }
        s.TaskQue.Append(task)
    }

    // use buffered channel to control max
    // number of concurrent go routines
    s.TaskChan = make(chan struct{}, config.ThreadCount)

    s.MaxDepth = config.MaxDepth

    s.CrawlInterval = config.CrawlInterval

    s.ThreadCount = config.ThreadCount

    s.TaskCommonCfg = taskCommonCfg
}

// Start to run tasks.
func (s *Scheduler) Start() {
    log.Logger.Info("start to run tasks")

    for {
        if s.TaskQue.Len() == 0 && len(s.TaskChan) == 0 {
            // 将新任务加入任务队列这一操作包含在任务中
            // len(s.TaskChan) == 0说明后续一定没有新任务被加入到s.TaskQue中
            log.Logger.Info("ms.TaskQue has been empty")
            break
        }

        // s.TaskQue.Len() == 0 && len(s.TaskChan) != 0说明当前还有任务在运行
        // 只是还没有生成新的任务加入到任务队列 有了uselessTask的存在 后续一定会有任务被加入到任务队列
        // s.TaskQue.Remove()可能会等待一段时间 但不会阻塞

        // s.TaskQue.Len() != 0 && len(s.TaskChan) == 0说明当前没有任务在运行
        // s.TaskQue不为空 直接从s.TaskQue里取任务然后运行即可

        // s.TaskQue.Len() != 0 && len(s.TaskChan) != 0说明当前还有任务在运行
        // s.TaskQue不为空 直接从s.TaskQue里取任务然后运行即可
        task := s.TaskQue.Remove()
        s.RunTask(task.(*Task))
    }

    close(s.TaskChan)
    log.Logger.Info("all tasks done")
}

// Run single task.
func (s *Scheduler) RunTask(task *Task) {
    if task.Depth >= s.MaxDepth {
        return
    }

    // 避免重复抓取
    // LoadOrStore是Go官方提供的sync.Map的一个方法 第一个参数为key 第二个参数为value
    // 如果task.Url已经存在于urlTable中了则返回的ok的值为true 否则ok的值为false并将task.Url加入到urlTable中
    if _, ok := s.UrlTable.LoadOrStore(task.Url, true); ok {
        // 该url的内容正在抓取或者已经抓取过了 直接返回
        return
    }

    s.TaskChan <- struct{}{}
    go func() {
        // uselessTask是为了在任务队列变为空之前排空TaskChan从而优雅退出 uselessTask一进入RunTask方法就会返回不会向TaskChan添加元素
        // 有的任务爬虫任务可能不会取到符合条件的子url（可能某个url下没有子url 也可能有子url但子url不能匹配正则表达式）
        // 不管有没有符合条件的子url都往任务队列里加一个uselessTask可以保证在Start方法的for循环里遇到
        // s.TaskQue.Len() == 0 && len(s.TaskChan) != 0的情况下不会阻塞
        uselessTask := NewUselessTask(s.MaxDepth)

        defer func() {
            log.Logger.Info("task %s done", task.Url)
            // append useless task
            s.TaskQue.Append(uselessTask)
            <- s.TaskChan
        }()

        // 控制抓取间隔 防止被封禁
        hostName, err := parser.ParseHostName(task.Url)
        if err != nil {
            log.Logger.Error("%s: parser.ParseHostName(): %s", task.Url, err.Error())
            return
        }
        timer, ok := s.TimerTable.LoadOrStore(hostName, time.NewTimer(time.Duration(s.CrawlInterval) * time.Second))
        if ok {
            select {
            case <- timer.(*time.Timer).C:
            }
            timer.(*time.Timer).Reset(time.Duration(s.CrawlInterval) * time.Second)
        }

        log.Logger.Info("start to crawl %s", task.Url)
        urlList, err := task.Run()
        if err != nil {
            log.Logger.Error("%s", err.Error())
            return
        }

        // generate new tasks
        for _, url := range urlList {
            nextTask := &Task{
                Url:       url,
                Depth:     task.Depth + 1,
                CommonCfg: s.TaskCommonCfg,
            }
            s.TaskQue.Append(nextTask)
        }
    }()
}