
https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/

先理解基本的MPG模型，使用golang代码进行模拟编码执行，理解内部基本运行原理，
先不用理解汇编底层代码，通读上面的这一章，

理解goroutine的历史进程。



- 单线程调度器 · 0.x
  只包含 40 多行代码；
  程序中只能存在一个活跃线程，由 G-M 模型组成；
  
- 多线程调度器 · 1.0
  允许运行多线程的程序；
  全局锁导致竞争严重；
  
- 任务窃取调度器 · 1.1
  引入了处理器 P，构成了目前的 G-M-P 模型；
  在处理器 P 的基础上实现了基于工作窃取的调度器；
  在某些情况下，Goroutine 不会让出线程，进而造成饥饿问题；
  时间过长的垃圾回收（Stop-the-world，STW）会导致程序长时间无法工作；

- 抢占式调度器 · 1.2 ~ 至今
  - 基于协作的抢占式调度器 - 1.2 ~ 1.13
    通过编译器在函数调用时插入抢占检查指令，在函数调用时检查当前 Goroutine 是否发起了抢占请求，实现基于协作的抢占式调度；
    Goroutine 可能会因为垃圾回收和循环长时间占用资源导致程序暂停；
    
  - 基于信号的抢占式调度器 - 1.14 ~ 至今
    实现基于信号的真抢占式调度；
    垃圾回收在扫描栈时会触发抢占调度；
    抢占的时间点不够多，还不能覆盖全部的边缘情况；

- 非均匀存储访问调度器 · 提案
  对运行时的各种资源进行分区；
  实现非常复杂，到今天还没有提上日程；














- single-thread
  - 获取调度器的全局锁；
  - 调用 runtime.gosave 保存栈寄存器和程序计数器；
  - 调用 runtime.nextgandunlock 获取下一个需要运行的 Goroutine 并解锁调度器；
  - 修改全局线程 m 上要执行的 Goroutine；
  - 调用 runtime.gogo 函数运行最新的 Goroutine；

- multi-thread
  - 调度器和锁是全局资源，所有的调度状态都是中心化存储的，锁竞争问严重；
  - 线程需要经常互相传递可运行的 Goroutine，引入了大量的延迟；
  - 每个线程都需要处理内存缓存，导致大量的内存占用并影响数据局部性（Data locality）；
  - 系统调用频繁阻塞和解除阻塞正在运行的线程，增加了额外开销；

- go1.1
  - 如果当前运行时在等待垃圾回收，调用 runtime.gcstopm 函数；
  - 调用 runtime.runqget 和 runtime.findrunnable 从本地或者全局的运行队列中获取待执行的 Goroutine；
  - 调用 runtime.execute 函数在当前线程 M 上运行 Goroutine；



- 但是 1.1 版本中的调度器仍然不支持抢占式调度，程序只能依靠 Goroutine 主动让出 CPU 资源才能触发调度。Go 语言的调度器在 1.2 版本4中引入基于协作的抢占式调度解决下面的问题5：
  - 某些 Goroutine 可以长时间占用线程，造成其它 Goroutine 的饥饿；
  - 垃圾回收需要暂停整个程序（Stop-the-world，STW），最长可能需要几分钟的时间6，导致整个程序无法工作；



所有 Goroutine 在函数调用时都有机会进入运行时检查是否需要执行抢占

因为这里的抢占是通过编译器插入函数实现的，还是需要函数调用作为入口才能触发抢占，所以这是一种协作式的抢占式调度。
[参考下面这个example](https://github.com/zput/Go-Questions/blob/master/goroutine%20%E8%B0%83%E5%BA%A6%E5%99%A8/%E4%B8%80%E4%B8%AA%E8%B0%83%E5%BA%A6%E7%9B%B8%E5%85%B3%E7%9A%84%E9%99%B7%E9%98%B1.md)
```golang
func main() {
    var x int
    threads := runtime.GOMAXPROCS(0)
    for i := 0; i < threads; i++ {
        go func() {
            for { x++ }
        }()
    }
    time.Sleep(time.Second)
    fmt.Println("x =", x)
}

```
> for 循环或者垃圾回收长时间占用线程，这些问题中的一部分直到 1.14 才被基于信号的抢占式调度解决。

- 编译器会在调用函数前插入 runtime.morestack；

- Go 语言运行时会在垃圾回收暂停程序、系统监控发现 Goroutine 运行超过 10ms 时发出抢占请求 StackPreempt；

- 当发生函数调用时，可能会执行编译器插入的 runtime.morestack 函数，它调用的 runtime.newstack 会检查 Goroutine 的 stackguard0 字段是否为 StackPreempt；

- 如果 stackguard0 是 StackPreempt，就会触发抢占让出当前线程；



基于信号的抢占式调度

https://github.com/zmeteor/linux/tree/master/signal/testSignal/testSIGURG
















