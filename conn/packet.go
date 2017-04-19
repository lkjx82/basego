package conn

// -----------------------------------------------------------------------------

import (
	"encoding/binary"
	"basego/obj_pool"
)

// -----------------------------------------------------------------------------
// 从1 byte 里取长度
func parseLen8(b []byte) int {
	return int(uint8(b[0]))
}

// 从2 byte 里取长度
func parseLen16(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}

// 从3 byte 里取长度
func parseLen32(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}

// -----------------------------------------------------------------------------
// 从1 byte 里取长度
func writeLen8(b []byte, leng int) {
	b[0] = byte(leng)
}

// 从2 byte 里取长度
func writeLen16(b []byte, leng int) {
	binary.BigEndian.PutUint16(b, uint16(leng))
}

// 从3 byte 里取长度
func writeLen32(b []byte, leng int) {
	binary.BigEndian.PutUint32(b, uint32(leng))
}

// -----------------------------------------------------------------------------
//
type Packet struct {
	Data    []byte // data
	DataLen int    // data length
	LenSize int8   // len seg size 	// 1, 2, 4
	ref     int    //
	pool    obj_pool.ObjPool
}

// -----------------------------------------------------------------------------
//
func NewPacket(size int, pool obj_pool.ObjPool) *Packet {

	if pool != nil {
		pi := pool.Alloc()
		if pi != nil {
			p := pi.(*Packet)
			if len(p.Data) < size {
				return nil
			} else {
				p.DataLen = 0
				p.ref = 1
				p.pool = pool
				return p
			}
		} else {
			return nil
		}
	}

	p := Packet{}
	p.Data = make([]byte, size, size)
	p.DataLen = 0
	p.ref = 1

	return &p
}

// -----------------------------------------------------------------------------
//
func (this *Packet) Append(b []byte) bool {
	if len(this.Data)-this.DataLen >= len(b) {

		if this.DataLen == 0 {
			this.DataLen = int(this.LenSize)
		}

		switch this.LenSize {
		case 1:
			copy(this.Data[this.DataLen:], b)
			this.DataLen = this.DataLen + len(b)
			writeLen8(this.Data, this.DataLen)
		case 2:
			copy(this.Data[this.DataLen:], b)
			this.DataLen = this.DataLen + len(b)
			writeLen16(this.Data, this.DataLen)
		case 4:
			copy(this.Data[this.DataLen:], b)
			this.DataLen = this.DataLen + len(b)
			writeLen32(this.Data, this.DataLen)

		default:
			return false
		}
	}
	return false
}

// -----------------------------------------------------------------------------
//
func (this *Packet) Release() {
	this.ref--
	if this.ref <= 0 {
		if this.pool != nil {
			this.pool.Free(this)
		}
	}
}

// -----------------------------------------------------------------------------
//
func (this *Packet) Dup() *Packet {
	this.ref++
	return this
}

// -----------------------------------------------------------------------------
//
func (this *Packet) Clone() *Packet {
	p := NewPacket(len(this.Data), this.pool)
	copy(p.Data, this.Data)
	p.DataLen = this.DataLen
	return p
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

//type PacketPool struct {
//	min   int
//	max   int
//	step  int
//	pools []obj_pool.ObjPool
//}

//// -----------------------------------------------------------------------------
////
//func NewPacketPool(min, max int) *PacketPool {
//	pp := PacketPool{}
//	pp.min = min
//	pp.max = max
//	n := (max-min)/1024 + 1
//	pp.pools = make([]obj_pool.ObjPool, n, n)

//	for i := 1; i <= n; i++ {
//		pp.pools[i] = obj_pool.NewObjPool(func() interface{} { return NewPacket(i * 1024) })
//	}
//	return &pp
//}

//// -----------------------------------------------------------------------------
////
//func (this *PacketPool) Alloc(size int) *Packet {
//	return nil
//}

// -----------------------------------------------------------------------------
