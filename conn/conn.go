package conn

// -----------------------------------------------------------------------------

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// -----------------------------------------------------------------------------

type NewConnFunc func(*Conn)

//type GenServ struct {
//}

// -----------------------------------------------------------------------------
//
type Conn struct {
	nc        net.Conn         // connect in net
	reader    *bufio.Reader    // 读bufio
	headerLen int              // 1,2,4
	parseFunc func([]byte) int // 解析长度头
	sndQ      chan *Packet     // 网络报发送缓冲 // sndQ      chan []byte      // send quque //
	sndQsize  int              // Send Queue 大小
	mutex     sync.Mutex       // 锁
	isActive  bool             // 是否关闭了
	serv      *Serv            // 所属serv
}

// -----------------------------------------------------------------------------
//
func (this *Conn) Init(headerLen int, sendQueueSize int, nc net.Conn, serv *Serv) {
	switch headerLen {
	case 1:
		this.parseFunc = parseLen8
	case 2:
		this.parseFunc = parseLen16
	case 4:
		this.parseFunc = parseLen32
	default:
		return
	}

	this.nc = nc
	this.serv = serv

	// lenSeg size
	this.headerLen = headerLen
	if this.reader == nil {
		this.reader = bufio.NewReaderSize(this.nc, 4096*2) // 2 倍最大包大小
	} else {
		this.reader.Reset(this.nc)
	}

	this.sndQ = make(chan *Packet, sendQueueSize)
	//	this.sndQ = make(chan []byte, sendQueueSize)

	// send queue
	this.sndQsize = sendQueueSize
	this.isActive = true
	go this.doSend()
}

// -----------------------------------------------------------------------------
// read
func (this *Conn) Read(pack *Packet) error {
	pack.DataLen = 0
	fmt.Println("recv")

	hd, err := this.reader.Peek(this.headerLen)
	if err != nil {
		fmt.Println(this, this.headerLen, "1", err, hd)
		return err
	}

	// 包大小
	packLen := this.parseFunc(hd)

	// 从换存取到 pack 里
	if n, err := this.reader.Read(pack.Data[:packLen]); err != nil {
		fmt.Println(this, "3", err)
		return err
	} else if n < packLen {
		fmt.Println("fuck")
	}

	pack.DataLen = packLen
	return nil
}

// -----------------------------------------------------------------------------
//
func (this *Conn) Close() {
	this.mutex.Lock()
	if this.isActive {
		this.nc.(*net.TCPConn).SetLinger(0)
		this.nc.Close()
		this.isActive = false
		close(this.sndQ)
		this.serv.FreeConn(this)
	}
	this.mutex.Unlock()
}

// -----------------------------------------------------------------------------

//func (this *Conn) Send(b []byte) {
//	this.mutex.Lock()
//	if this.isActive {
//		this.sndQ <- b
//	}
//	this.mutex.Unlock()
//}

// 发送
func (this *Conn) Send(p *Packet) {
	this.mutex.Lock()
	if this.isActive {
		this.sndQ <- p
	}
	this.mutex.Unlock()
}

// -----------------------------------------------------------------------------
//
func (this *Conn) doSend() {
	//	defer fmt.Println("doSend end!!!!!!!!!!!!!!")
	defer this.Close()

	for p := range this.sndQ {
		if p == nil {
			fmt.Println("Break。。。。。。。。。。。。。。。。。。")
			return
		}

		fmt.Println("doSend")

		leftSend := p.DataLen

		for {
			if n, err := this.nc.Write(p.Data[p.DataLen-leftSend : p.DataLen]); err != nil {
				fmt.Println(p.DataLen-leftSend, ":", n, err)
				fmt.Println("Break!!!!!!!!!!!!!!!!!")
				return
			} else {
				//				fmt.Println(leftSend, n)
				leftSend = leftSend - n
				if leftSend == 0 {
					break
				}
			}
		}

		p.Release()
	}
}

// -----------------------------------------------------------------------------
