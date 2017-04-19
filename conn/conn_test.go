package conn_test

// -----------------------------------------------------------------------------

import (
	"bytes"
	"fmt"
	conn "github.com/lkjx82/basego/conn"
	"github.com/lkjx82/basego/obj_pool"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
//
func TestPacket(t *testing.T) {

	lenSizes := []int8{1, 2, 4}

	// ------------------
	for i := 0; i < 100; i++ {
		p := conn.NewPacket(200, nil)
		if p == nil {
			t.Fail()
		}
		p.LenSize = lenSizes[i%len(lenSizes)]

		b := []byte(`1234567890abcdefg`)
		p.Append(b)

		if p.DataLen != len(b)+int(p.LenSize) {
			fmt.Println(p.DataLen, len(b), p.LenSize)
			t.Fail()
		}

		p1 := p.Dup()
		if !(bytes.Equal(p1.Data, p.Data) && p.DataLen == p1.DataLen) {
			t.Fail()
		}
		p1.Release()

		p2 := p.Clone()
		if !(bytes.Equal(p2.Data, p.Data) && p.DataLen == p2.DataLen) {
			t.Fail()
		}

		p2.Release()

		//		time.AfterFunc(time.Second, func() { p2.Release() })

		//		fmt.Println(&p, &p1, &p2)
		fmt.Printf("%p, %p, %p\r", p, p1, p2)

		p.Release()
	}

	objPool := obj_pool.NewObjPool(func() interface{} {
		return conn.NewPacket(200, nil)
	})

	// ------------------
	for i := 0; i < 100; i++ {
		p := conn.NewPacket(200, objPool)
		if p == nil {
			t.Fail()
		}
		p.LenSize = lenSizes[i%len(lenSizes)]

		b := []byte(`1234567890abcdefg`)
		p.Append(b)

		if p.DataLen != len(b)+int(p.LenSize) {
			fmt.Println(p.DataLen, len(b), p.LenSize)
			t.Fail()
		}

		p2 := p.Dup()
		if !(bytes.Equal(p2.Data, p.Data) && p.DataLen == p2.DataLen) {
			t.Fail()
		}
		p2.Release()

		p2 = p.Clone()
		if !(bytes.Equal(p2.Data, p.Data) && p.DataLen == p2.DataLen) {
			t.Fail()
		}
		p2.Release()

		fmt.Println(p.Data)
		p.Release()
	}

	// ------------------
	// ------------------
	for i := 0; i < 100; i++ {
		p := conn.NewPacket(210, objPool)
		if p != nil {
			t.Fail()
		}
	}

}

//func TestPacket(t *testing.T) {
//	lenSize := []int8{1, 2, 4}
//	for _, i := range lenSize {
//		p := conn.NewPacket(100, nil)

//	}
//}

// -----------------------------------------------------------------------------

func TestServ(t *testing.T) {
	svr := conn.NewServ(2, 1024, 4096, 20)

	objPool := obj_pool.NewObjPool(func() interface{} {
		return conn.NewPacket(200, nil)
	})

	svr.Listen(":9009",
		func(c *conn.Conn) {
			i := 0
			// recv
			go func() {
				for {
					p := conn.NewPacket(20, objPool)
					i++
					if i == 3 {
						c.Close()
						return
					}
					err := c.Read(p)
					fmt.Println("listen recv:", err, p)
					// proccess p
					<-time.After(time.Second * 1)
					c.Send(p)
				}
			}()
		})

	fmt.Println("---------------------")

	svr.Conn("localhost:9009",
		func(c *conn.Conn) {
			i := 0
			// recv
			go func() {
				for {
					p := conn.NewPacket(20, objPool)

					i++
					if i == 5 {
						//						c.Close()
					}
					err := c.Read(p)
					fmt.Println("conn recv", err, p)
					if err != nil {
						c.Close()
						return
					}
					// proccess p
					<-time.After(time.Second * 1)
					c.Send(p)
				}
			}()

			p := conn.NewPacket(100, nil)
			p.LenSize = 2

			b := []byte(`1234567890abcdefg`)
			p.Append(b)
			fmt.Println(p.Data)

			c.Send(p)

			<-time.After(time.Second * 5)

		})

	<-time.After(time.Second * 30)

}

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
