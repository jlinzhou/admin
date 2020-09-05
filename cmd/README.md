指针就是地址。指针变量就是存储地址的变量。

*p:解引用、间接引用。

栈帧：用来给函数运行提供内存空间。取内存于stack上。

​			当函数调用时，产生栈帧。函数调用结束，释放栈帧。

​			栈帧存储：1.局部变量。2形参。（形参与局部变量存储地位等同）3.内存字段描述值

指针使用注意：

​	空指针：未被初始化的指针。var p *int    *p-->err

​	野指针：被一片无效的地址空间初始化。



申请的内存空间在栈上会自动释放，在堆上不会自动释放,new创建的是在堆上。



格式化输出：

​	%q:以Go语言格式显示字符串。默认带有“”符

​	%v:显示对应数据详细信息

变量存储：

​	等号 左边的变量，代表变量所指向的内存空间

​	等号 右边的变量， 代表变量内存空间存储的数据值

指针的函数传参：

​	传地址（引用）：将地址值作为函数参数，返回后传递。

​	传值（数据）：将实参的值拷贝一份给形参。

​	传引用：在A栈帧内部，修改B栈帧中的变量值





## 切片

​	为什么用切片：

​		1.数组的容量固定，不能自动拓展。

​		2.值传递，数组作为函数参数时，将整个数组值拷贝一份给形参。

​		在Go语言中。我们几乎可以在所有场景中，使用切片替换数组使用。

​	切片的本质，不是一个数组的指针，是一种数据结构，用来操作数组内部元素.

​		runtime/slice.go  type  slice struct { *p len cap }

​	切片的使用：

​		数组和切片定义区别：

​			创建数组时[]指定数组长度。

​			创建切片时，[]为空，或者...

​		切片名称[low:high:max]

​		low:起始下标位置

​		high:结束下标位置 len=high-low

​		容量：cap=max-low  是从它的第一个元素开始数，到其底层数组元素末尾的个数 。

​		截取数组，初始化切片时，切片容量跟随原数组。

​		s[:high:max]:从0开始，到high结束。（不包含）

​		s[low:]:从low开始，到末尾

​		s[:high]:从0开始，到high结束。容量跟随原先容量。【常用】

​	切片创建：

​		1.自动推导类型创建切片。slice:=[]int{1,2,3,4}

​		2.slice:=make([]int,长度，容量)

​		3.slice:=make([]int,长度)创建切片时，没有指定容量，容量==长度

​		**注意：make只能创建slice、map和channel，并且返回一个有初始值（非零）的对象**

​	切片做函数参数---传引用。（传地址）

​	append:

​		向切片增加元素时，切片的容量会自动增长。1024以下时，以两倍方式增长。

​	copy:

​		copy(目标位置切片，源切片)，拷贝过程中，直接对应位置拷贝

​		copy(data[idx:],data[idx+1:]) return data[:len(data)-1]


## map

字典、映射   key-value  key:唯一、无序。不能是引用类型数据。

​	map不能使用cap()

​	创建方式：

​		1.var m1 map[int]string  ---不能存储数据

​		2.m2:=map[int]string --能存储数据

​		3.m3:=make(map[int]string)--默认len=0

​		4.m4:=make(map[int]string,10)



​	删除map:

​		delete()函数

​	strings.Fields()将字符串拆分成字符串切片



## 结构体

​	是一种数据类型

普通变量定义和初始化

​		1.顺序初始化：依次将结构体内部成员初始化。

​		2.制定成员初始化。---未初始化的成员变量，取该数据类型对应得默认值

​	结构体变量的比较和赋值：

​		1.比较：只能使用==和!= 不能使用> < >= <= 

​		2.相同结构体类型(成员变量的类型、个数、顺序一致)变量之间可以直接赋值。

​	结构体地址：

​		结构体变量的地址==结构体首个元素的地址。

​	结构体传参：

​		unSafe.Sizeof(变量名)---->此种类型的变量所占用内存空间的大小

​		将结构体变量的值拷贝一份，传参。--几乎不用，内存消耗大，效率低。

指针变量定义和初始化：

​		1.顺序初始化：依次将结构体内部所有成员初始化。

​		2.new(person)

​	结构体地址：

​		结构体变量的地址==结构体首个元素的地址。

​	结构体指针传参：

​		unSafe.Sizeof(指针)：不管何种类型的指针，在64位操作系统下，大小一致，均为8字节。

​		将结构体变量的地址值，传递（传引用）。---使用频率非常高！！！

​	结构体指针做函数返回值：

​		不能返回局部变量的地址值。------局部变量保存在栈帧上，函数调用结束后，栈帧释放，局部变量的地址，不再受系统保护，随时可能分配给其他程序。

​		可以返回局部变量的值。



## 字符串

​	1.字符串分割 strings.Split("1.2.3",".")

​	2.字符串以空格分割strings.Fileds("")

​	3.判断字符串结束标记HasSuffix   strings.HasSuffix("test.abc","abc")

​	4.判断字符串起始标记HasPrefix  strings.HasPrefix("test.abc","test")

## 文件

​	打开创建文件:

​		1.创建文件 Create:文件不存在，文件存在，将文件清空。

​				参数:name ，打开文件的路径：绝对路径、相对路径

​		2.打开文件Open: 以只读文件打开。

​				参数:name ，打开文件的路径：绝对路径、相对路径

​		3.打开文件OpenFile:以只读，只写、读写方式打开。

​				参1：name ，打开文件的路径：绝对路径、相对路径

​				参2：打开文件权限：O_RDONLY、O_WRONLY 、O_RDWR

​				参3:一般传6/7

​	写文件：

​		按字符串写：WriteString()-->n个写入的字符个数

​			**回车换行：windows:\r\n   linux:\n**

​		按位置写

​			Seek():修改文件的读写指针位置。

​			参1：偏移量。正：向文件尾偏。负：向文件头偏。

​			参2：偏移量起始位置。

​				io.SeekStart  文件起始位置

​				io.SeekCurrent  文件当前位置

​				io.SeekEnd  文件结尾位置

​			返回值:表示文件起始位置，到文件读写指针位置的偏移量。

​			off,_:=f.Seek(-5,io.SeekEnd)

​		按字节写

​			WriteAt():在文件指定偏移位置写入[]byte，通常搭配Seek()

​			参1：待写入的数据

​			参2：偏移量

​			返回值：写入文件的字节数

​			n,_=f.WriteAt([]byte("111"),off)

​	读文件：

​		**按行读(常用)**

​			1.创建一个带有缓冲区的Reader(读写器)

​				reader:=bufio.NewReader(打开的文件指针)

​			2.从reader的缓冲区中，读取指定长的数据，数据的长度取决于参数dlime

​				buf,err:=reader.ReadBytes('\n')按行读。

​				判断到达文件末尾:  if err!=nil && err== io.EOF  到文件结尾。

​							文件结束标记，是要单独读一次获取到的。

​			**缓存区：内存中的一块区域，用来减少物理硬盘访问操作**。<<计算硬件及组成原理>>---机械工业出版社

​	按字节读写文件。

​		read():按字节读文件

​		write()：按字节写文件

​	目录操作：

​		打开目录：OpenFile

​			参1：name ，打开目录的路径：绝对路径、相对路径

​			参2：打开目录权限：O_RDONLY、O_WRONLY 、O_RDWR

​			参3:os.ModeDir

​		返回值：返回一个可以读写目录的指针

​		读目录：Readdir

## 并行并发

​	并行：在同一时刻，有多条指令在多个处理器上同时执行。	借助多核cpu实现    (真  并行)

​	并发：在同一时刻只能有一条指令执行，但多个进程指令被快速的轮换执行。

​		宏观：用户体验上，程序在并行执行。

​		微观：多个计划任务，顺序执行，在飞快的切换。轮换使用cpu时间轮换。  (假 并行)

​	进程并发：

​		程序：编译成功得到的二进制文件。 占用磁盘空间。 死的  1

​		进程：运行起来的程序。占用系统资源。（内存）      活的  N

​	进程状态：

​		初始态、就绪态、运行态、挂起（阻塞）态、终止（停止）态。

​	孤儿进程:	父进程先于子进程结束，则子进程成为孤儿进程，子进程的父进程成为init进程，称为init进程领养孤儿进程。

​	僵尸进程：进程终止，父进程尚未回收，子进程残留资源（PCB）存放于内核中，变成僵尸（Zombie）进程.



### 线程并发

​	LWP  **轻量级的进程 **，在linux依然是进程。  **最小的执行单位 **。----cpu分配时间轮片的对象。

​	进程： **最小的系统资源分配单位。 **

​	目的是为了并行的争夺cpu，提高cpu执行速度。

​	同步：

​		协同步调，规划先后顺序。

​		线程同步机制：

​			互斥锁(互斥量)：建议锁，拿到锁以后，才能访问数据，没有拿到锁的线程，阻塞等待，等到锁的线程释放锁。

​			读写锁：一把锁（读属性、写属性）。写独占，读共享。写锁优先级高。

### 协程并发

​	协程，coroutine， **轻量级线程。 **

​	python、lua、Ruset...

​	21世纪。

​	目的是为了程序的执行效率。在阻塞等待期间可以做其他的事情。

总结：进程、线程、协程都可以完成并发。

​			稳定性强、节省资源、效率高。

​	


###  Goroutine:go程
 并不是由操作系统调度的 ， 避免了上下文切换的额外耗费 

**奉行通过通信来共享内存，而不是通过共享内存来通信**

​	创建于进程中。直接使用go关键字，放置于函数调用前面，产生一个go程，并发。

​	goroutine的特性：

​		主go程结束，子go程随之退出。

​	

​	runtime包

​		runtime.Gosched():出让当前go程所占用的cpu时间片。当再次获得cpu时，从出让位置继续恢复执行。

​		----时间片轮转调度算法。

​		runtime.Goexit():

​			return:返回m8当前函数调用到调用者那里去，return之前的defer注册生效。

​			Goexit:退出当前go程。结束调用该函数的当前go程，goexit()之前注册的defer都生效。

​		runtime.GOMAXPROCS():设置可以并行计算的cpu核数的最大值，并返回之前的值。

=========

[补充知识点]

​	每当有一个进程启动时，系统会自动打开三个文件：标准输入、标准输出、标准错误。-----对应三个文件stdin、stdout、stderr。当进程运行结束，操作系统自动关闭三个文件。

### channel

​		是一种数据类型，对应一个“管道”。

​		channel的定义：

​			make（chan 在channel中传递的数据类型，容量)  容量=0：无缓冲；容量>0，有缓冲channel。

​			e.g. make(chan int)或make(chan string ,0)

​		len(ch)读取channel剩余读取数据的个数，cap(ch)通道的容量。

​	无缓冲channel:----同步通信

​		通道容量为0，len=0.不能存储数据。

​		channel应用于两个go程中。一个读，另一个写。

​		具备同步的能力。读写同步。

​	有缓冲channel----异步通信

​		通道容量为非0.

​		channel应用于两个go程中。一个读，另一个写。

​	关闭channel

​		确定不再向对端发送数据，使用close(ch)来关闭channel

​		对端可以判断channel是否关闭

​			if num,ok:=<-ch;ok{}

​			可以使用range替换ok

​		总结：1.数据不发送完，不应该关闭

​					2.已经关闭的channel，不能再向其写入数据。

​					3.写端已经关闭channel，可以从中读取数据，读到0.----说明：写端关闭.

​		单向channel:

​			默认的channel是双向的. var ch chan int   ch=make(chan int)

​			单向写channel：var sendCh chan<-int   sendCh=make(chan<-int)不能读操作

​			单向读channel:   var recvCh <-chan int   recvCh=make(<-chan int)

​			转换：

​				1.双向channel可以隐式转换为任意一种单向channel

​						sendCh=ch

​				2.单向channel不能转换为双向channel

​					ch=sendCh    error!

​				**传参，传引用**

​	生产者消费者模型

​		生产者：

​		消费者：

​		缓冲区：1.解耦（降低生产者和消费者之间的耦合度）

​						2.并发 (生产者消费者数量不对等时，能保持正常通信)

​						3.缓存 (生产者和消费者数据处理速度不一致时，暂存数据)

​	time  sleep()  After() NewTimer() NewTicker()

​	定时器

​		Timer:创建定时器，指定定时时长，定时到达后。系统会自动向定时器的成员C写系统当前时间。

​		time.After()

​		time.NewTimer()

​		定时器的停止、重置：

​			1)创建定时器myTimer:=time.NewTimer(2*time.Second)

​			2)停止：myTimer.Stop  ---将定时器归零.  <-myTimer.C会阻塞

​			3)重置：myTimer.Reset(time.Second)

​		time.Stop()设置定时器停止

​		time.Reset()定时器重置

select    其中case语句里必须是一个IO操作。一般不写default，避免忙轮询。

​	作用：用来监听channel上的数据流动方向。读？写？

​	用法：参考switch case语句。但case后面必须是IO操作，不可以任意写判别表达式。

​	注意事项：

​		1.监听的case中，没有满足监听条件，阻塞。

​		2.监听的case中，有多个满足监听条件，任选一个执行。

​		3.可以使用default来处理所有case都不满足监听条件的状况。通常不用（会产生忙轮询）

​		4.select自身不带有循环机制，需借助外层for来循环监听。

​		5.break只能跳出select，类似于switch中的用法。



死锁：

​	1.单go程自己死锁

​		channel应该在至少2个以上的go程中进行通信，否则死锁。

​	2.go程间channel访问顺序导致死锁

​		使用channel一端读（写），要保证另一端写（读）操作，同时有机会执行。否则死锁

​	3.多go程，多channel交叉死锁

​		Ago程，掌握M的同时，尝试N；Bgo程，掌握N的同时尝试拿M.

​	4.在go语言中，尽量不要将互斥锁、读写锁与channel混用。----隐性死锁

互斥锁：（互斥量）

​	A、 B  go程共同访问共享数据。由于cpu调度随机，需要对共享数据访问顺序加以限定（同步）。

​	创建mutex(互斥锁)，访问共享数据之前，加锁，访问结束，解锁。在Ago程加锁期间。Bgo程加锁会失败---阻塞。直至Ago程解锁mutex，B从阻塞处，恢复执行。

读写锁：

​	读时共享，写时独占。写锁优先级比读锁高。



条件变量：sync.Cond

​	本身不是锁!经常要与锁结合使用!!

​	使用流程：

​	1.创建条件变量： var cond   sync.Cond

​	2.指定条件变量用的锁  conL=new(sync.Mutex)

​	3.cond.L.Lock()  给公共区加锁（互斥量）

​    4.判断是否到达阻塞条件（缓冲区满/空）---for 循环判断

​			for len(ch)==cap(ch) {  cond.Wait()   ----1) } 阻塞   2）解锁  3)加锁

​	5.访问公共区----读、写数据、打印

​	6.解锁条件变量用的锁 cond.L.Unlock()

​	7.唤醒阻塞在条件变量上的对端。signal()  Broadcast()





## 网络通信

​	协议：一组规则。要求使用协议的双方，必须严格遵守协议内容。

​	网络分层架构：

​		OSI/RM（理论上的标准）应用层 表示层 会话层 传输层  网络层  数据链路层  物理层（物数网传会表应）

​		TCP/IP(事实上的标准)应用层  传输层  网络层  链路层(链网传应)

​		越往右的层，越靠近硬件；越往左的层，越靠近用户。

​	各层功能：

​		链路层：ARP   

​						源MAC----目标MAC

​						ARP协议作用：借助IP获取mac地址

​		网络层：IP

​						源IP-----目标IP

​						IP协议的作用：在网络环境中唯一标识一台主机

​						IP地址本质：2进制数。-----点分十进制IP地址（string)

​		传输层：TCP/UDP

​						port----在一台主机上唯一标识一个进程。

​		应用层：ftp、http、自定义

​						对数据进行封装、解封装

数据通信过程：

​	封装：应用层----传输层----网络层----链路层

​	解封装：链路层---网络层----传输层-----应用层

总结通信过程：

​	1.mac地址（不需要用户指定）  （ARP协议）ip------->mac

​	2.IP地址（需要用户指定）-----确定主机

​	3.port端口号（需要用户指定）----确定程序

​			1.不能使用系统占用的默认端口。5000+端口我们使用（8080）

​			2.65535为端口上限



socket编程：(双向全双工)

​	网络通信过程中：socket一定是成对出现的。

​	

网络设计模式：

​	C/S：

​		优点：数据传输效率高，协议选择灵活。

​		缺点：工作量大，安全性构成威胁。

​	B/S:

​		优点：开发工作较小，不受平台限制，安全威胁小。

​		缺点：缓存数据差，协议选择不灵活。

TCP-CS客户端：

​	1.conn,err:=net.Dial("TCP",服务器的IP+port)

​	2.写数据的服务器conn.Write()

​	3.读取服务器回发的数据conn.Read()

​	4.conn.Close()

TCP-CS并发服务端：

​	1.创建监听套接字listener:=net.Listen("tcp",服务器的IP+port)

​	2.defer listener.Close()

​	3.for 循环阻塞监听客户端连续事件 conn:=listen.Accept()

​	4.创建go程对应每一个客户端进行数据通信 go HandlerConnect()

​	5.实现HandlerConnect(conn net.Conn)

​			1)defer conn.Close()

​			2)获取成功连接的客户端Addr    conn.RemoteAddr()

​			3)for  循环  读取客户端发送数据  conn.Read(buf)

​			4)处理数据 小----大

​			5)回写转换后的数据  conn.Write(buf[:n])



服务器判断关闭：

​	**read读客户端，返回0----对端关闭!**



三次握手

​	客户端主动发起连接请求端，标志位SYN  2000(0),被动接收连接请求端，请求端回复ACK  2001和SYN 7000(0)，客户端回复ACK 7001  三次握手完成---连接建立完成.

​	1.主动发起请求端，发送SYN

​	2.被动建立连接请求端，应答ACK同时发送SYN

​	3.主动发起请求端，发送应答ACK

​	标志TCP三次握手建立完成.---server.Accept()返回。  ----client.Dial()返回

四次挥手(因为半关闭)

​	客户端主动发出关闭请求，标志位FIN 2000(0) ，被动关闭连接请求端，回复ACK 2001，这时主动关闭端处于半关闭，被动关闭连接请求端再次回复FIN  5000(0)，主动关闭端回复ACK 5001 ,四次挥手完成.

​	1.主动关闭连接请求端，发送FIN

​	2.被动关闭连接请求端，应答ACK      标志。半关闭完成。-----close()(服务端)

​	3.被动关闭连接请求端，发送FIN  

​	4.主动关闭请求端，应答ACK   标志。四次挥手建立完成----close()（客户端）

​	



​	





​	