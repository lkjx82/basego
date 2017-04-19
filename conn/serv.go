package conn

import (
	"github.com/lkjx82/basego/obj_pool"
	"net"
	"sync"

	"fmt"
)

// -----------------------------------------------------------------------------
//
func newConn() interface{} {
	return &Conn{}
}

// -----------------------------------------------------------------------------
//
func NewServ(lenSize, maxConn, maxPacketSize, sndQSize int) *Serv {
	s := Serv{}
	s.pool = obj_pool.NewObjPool(newConn)
	s.maxConn = maxConn
	s.lenSize = lenSize
	s.maxPacketSize = maxPacketSize
	s.sndQSize = sndQSize
	s.connMap = make(map[*Conn]struct{})
	return &s
}

// -----------------------------------------------------------------------------
// 连接
type Serv struct {
	pool          obj_pool.ObjPool   // conn pool
	maxConn       int                // max client connect allowed
	ln            net.Listener       // listener
	connMap       map[*Conn]struct{} // conn map
	mutex         sync.Mutex         // lock for map
	lenSize       int                // 长度字段的字节数
	maxPacketSize int                // 最大包的大小
	sndQSize      int
	// connCnt int             // connCnt

}

// -----------------------------------------------------------------------------
// 启动 server
func (this *Serv) Listen(addr string, newConnFunc NewConnFunc) {
	if ln, err := net.Listen("tcp", addr); err != nil {
		fmt.Println(err)
		return
	} else {
		this.ln = ln
	}

	//
	go func() {
		for {
			if nc, err := this.ln.Accept(); err != nil {
				// fmt.Println(err.(net.Error))
			} else {
				if this.maxConn <= len(this.connMap) {
					nc.Close()
					continue
				}

				if c := this.pool.Alloc(); c != nil {
					conn := c.(*Conn)
					conn.Init(this.lenSize, this.sndQSize, nc, this)
					this.mutex.Lock()
					this.connMap[conn] = struct{}{}
					this.mutex.Unlock()
					go newConnFunc(conn)
				}
			}
		}
	}()
}

// -----------------------------------------------------------------------------

func (this *Serv) Close() {
	this.ln.Close()
}

// -----------------------------------------------------------------------------
//
func (this *Serv) FreeConn(c *Conn) {
	this.mutex.Lock()
	delete(this.connMap, c)
	this.mutex.Unlock()
}

// -----------------------------------------------------------------------------
//
func (this *Serv) Conn(addr string, newConnFunc NewConnFunc) {
	if nc, err := net.Dial("tcp", addr); err != nil {
		fmt.Println(err)
	} else {
		if c := this.pool.Alloc(); c != nil {
			conn := c.(*Conn)
			conn.Init(this.lenSize, this.sndQSize, nc, this)
			this.mutex.Lock()
			this.connMap[conn] = struct{}{}
			this.mutex.Unlock()
			newConnFunc(conn)
		}
	}
}

// -----------------------------------------------------------------------------
