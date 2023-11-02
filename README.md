# 项目名称

mini-spider

## 背景

在调研过程中，经常需要对一些网站进行定向抓取。这是一个使用Go语言开发的迷你定向抓取器，实现对种子链接的抓取，并把URL长相符合特定正则表达式的网页保存到磁盘上。

## 主要功能

实现对种子链接的爬取，并把URL符合特定正则表达式的网页保存到磁盘上。

## 快速开始
### 构建

在项目根目录下执行如下命令：

```bash
sh build.sh
```

默认生成三个目录，其中bin用于存放可执行程序，log用于存放日志，output用于存放下载的网页。

### 参数

* -h 显示帮助
* -c 指定配置文件目录 默认../conf
* -l 指定日志文件目录 默认../log
* -v 显示版本

## 设计思路
### 程序初始化

* 读取命令行参数、初始化日志、读取配置文件、种子文件
* 创建并初始化Scheduler，调用其Start方法开始调度

### 调度主要逻辑

* 退出for循环的条件是任务队列为空且channel为空
* 每次循环，Scheduler从任务队列中取出任务，然后执行任务

### 任务执行主要逻辑

* 判断当前任务的深度，大于等于MaxDepth则返回
* 判断当前任务的URL是否已经爬取过，若是则直接返回，否则开启go routine异步执行任务
* 获取当前任务的域名（站点），检查是否满足爬取间隔的要求
* 根据URL爬取网页，失败则记录日志并返回
* 判断其Content-Type是不是文本，不是则记录日志并返回
* 将爬取到的网页转换成UTF-8格式
* 判断该URL是否满足目标正则表达式，若满足则将其保存至磁盘
* 解析爬取到的网页，将其子URL加入任务队列

### 并发控制

* 最大并发数通过buffered channel控制

### 控制抓取间隔

* 通过sync.Map和time.Timer实现
* sync.Map的Key为hostname，Value为timer
* 每次执行抓取任务前通过任务的URL解析出hostname，通过hostname拿到该站点的timer，等待timer的剩余时间后，重置timer执行抓取任务

### 优雅退出

* 引入useless task保证taskChan在taskQue排空之前排空，详见代码和注释

## 测试
除main包外，每个包下都有单元测试代码，可进行测试。

需要注意的是，在scheduler包内，由于多个测试函数都有创建目录、删除目录操作，进行测试时需要限制并发数为1，可执行如下命令进行测试：

```bash
go test -parallel 1
```
