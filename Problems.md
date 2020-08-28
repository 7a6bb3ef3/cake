### 开发问题合集

> 这个文件用于记录一些开发中出现的Bug/问题，以及采取的解决方案


#### TCP连接关闭的恰当时机

不同于其他TCP应用，作为代理方首先是无从得知src(net.Conn，下同)dst究竟有多少数据要写，多少数据要读，并且HTTP并非是完全的一问一答模式，TLS有握手，
100状态码也可能会在中途出现，因此上行跟下行一定是需要并发执行，以下就是一个错误的示例：
```
_ ,_ = io.Copy(dst ,src)
_ ,_ = io.Copy(src ,dst)
```
至少应当跑一个`Copy`到Goroutine里
```
go func(){
    _ ,_ = io.Copy(dst ,src)
}
_ ,_ = io.Copy(src ,dst)
```
但这样写也可能有问题，往往在Handle一个net.Conn时，会在连接建立之后，写入`defer dst.Close()`，后一个Copy跑完就Close了连接，想来应该是不太合适的，
因此往代码里加了`sync.WaitGroup`，两个`Copy`都跑到了协程里。

然后就出现了新的问题：src和dst关不掉，两条协程并没有很顺利的跑完然后`wg.Done()`。目前认为造成这个情况的原因是，作为代理没有将浏览器的EOF通知到主机，
导致浏览器和远端主机感知不到对方的EOF，一个解决方案是，上行过程中读完src的数据后，执行`dst.(*net.TcpConn).CloseWrite()`：
```
	go func() {
		defer func() {
			wg.Done()
			dst.(*net.TCPConn).CloseWrite()
		}()
		up, e := ahoy.CopyConn(dst, src, outboundEnv)
		if e != nil {
			log.Info(addr.Address(), " src request.", e)
			return
		}
	}()
```

这里其实查了一些资料，又自己抓包、测试之后基本确认，`CloseWrite`或者说`SHUT_WR`在网络层是给了对端一个EOF。
而`Close()`会导致读写均不可用，不能提前调用。在收到src和dst的EOF之后，再调用`Close()`应该是合适的。


#### WINDOWS注册表更改后无法立刻生效

以下为修改windows系统代理的代码段，修改后无法立即生效，需要打开`ms-settings:network-proxy`(即设置里的代理页面)，观察一眼才能生效，
这种需要观测它一下它的状态才坍缩下来的感觉很是神必，于是给代理加一个开关的功能遇到了瓶颈。每次更新后强制call出代理页面也算是一个办法，但太暴力了。

```
	key ,e := registry.OpenKey(registry.CURRENT_USER ,"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings" ,registry.ALL_ACCESS)
	if e != nil{
		log.Println(e)
		return
	}
	defer key.Close()
	key.SetBinaryValue("ProxyEnable" ,[]byte{0})
```

最后找到的办法是，直接修改`Software\Microsoft\Windows\CurrentVersion\Internet Settings\Connections\DefaultConnectionSettings`就可以立即生效，
说来也是奇妙深刻，对`ProxyEnable` `ProxyServer`的修改也会反映到`DefaultConnectionSettings`上面，但是就是没法生效。
利用微软的ProcessMonitor检查了注册表修改记录后，确认了之前打开设置页面才生效的原因，其实是因为打开页面的同时svchost会重新修改注册表，
`DefaultConnectionSettings`也是这么找到的。
