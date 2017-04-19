package conn

//// data block
//// -----------------------------------------------------------------------------
//import (
//	"gsvr-kit/base/obj_pool"
//)

//// -----------------------------------------------------------------------------

//type DataBlock struct {
//	Data    []byte
//	DataLen int
//	ref     int
//	pool    obj_pool.ObjPool
//}

//// -----------------------------------------------------------------------------

//func NewDataBlock(size int, pool obj_pool.ObjPool) *DataBlock {
//	if pool != nil {
//		db := pool.Alloc ()
//		if db != nil {
//			if len(db.Data) < size {
//				return nil
//			}
//			return db
//		} else {
//	}

//	db := DataBlock{}
//	db.DataLen = size
//	db.pool = pool
//	if pool != nil {
//		b := pool.Alloc()
//		if b != nil {
//			db.Data = b.([]byte)
//			if len(db.Data) < size {
//				pool.Free(b)
//				return nil
//			}
//		} else {
//			return nil
//		}
//	} else {
//		db.Data = make([]byte, size)
//	}

//	db.ref = 1
//	return &db
//}

//// -----------------------------------------------------------------------------

//func (this *DataBlock) Dup() *DataBlock {
//	this.ref++
//	return this
//}

//// -----------------------------------------------------------------------------

//func (this *DataBlock) Clone() *DataBlock {
//	db := NewDataBlock(this.DataLen, this.pool)
//	copy(db.Data, this.Data)
//	return db
//}

//// -----------------------------------------------------------------------------

//func (this *DataBlock) Release() {
//	this.ref--
//	if this.ref <= 0 {
//		if this.pool != nil {
//			this.pool.Free(this)
//		}
//		return
//	}
//}

//// -----------------------------------------------------------------------------

//func (this *DataBlock) Size() int {
//	return this.DataLen
//}

//// -----------------------------------------------------------------------------
//// -----------------------------------------------------------------------------
////
//type MsgBlock struct {
//	db     *DataBlock
//	w      int
//	r      int
//	ref    int
//	cont   *MsgBlock
//	pool   obj_pool.ObjPool
//	dbPool obj_pool.ObjPool
//}

//// -----------------------------------------------------------------------------
////
//func NewMsgBlock(size int, pool obj_pool.ObjPool, dbPool obj_pool.ObjPool) *MsgBlock {
//	db := NewDataBlock(size, dbPool)
//	if db == nil {
//		return nil
//	}
//	mb := MsgBlock{}
//	mb.db = db
//	mb.ref = 1
//	mb.pool = pool
//	return &mb
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Wr(n int) {
//	this.w = this.w + n
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Rd(n int) {
//	this.r = this.r + n
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Rb() []byte {
//	return this.db.Data[this.r:]
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Copy(b []byte) int {
//	n := copy(this.db.Data, b)
//	this.w = n
//	return n
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Append(b []byte) int {
//	n := copy(this.db.Data[this.w:], b)
//	this.w = this.w + n
//	return n
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Cap() int {
//	return this.db.Size() - this.w
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Len() int {
//	return this.r - this.w
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Size() int {
//	return this.db.Size()
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Dup() *MsgBlock {
//	mb := NewMsgBlock(this.Size(), this.pool, this.dbPool)
//	mb.db = this.db.Dup()
//	mb.w = this.w
//	mb.r = this.r
//	return mb
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Clone() *MsgBlock {
//	mb := NewMsgBlock(this.Size(), this.pool, this.dbPool)
//	mb.db = this.db.Clone()
//	mb.w = this.w
//	mb.r = this.r
//	return mb
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Cont(mb *MsgBlock) *MsgBlock {
//	this.cont = mb
//	return this
//}

//// -----------------------------------------------------------------------------
////
//func (this *MsgBlock) Release() *MsgBlock {
//	this.db.Release()
//	this.db = nil
//	this.r = 0
//	this.w = 0

//	if this.pool != nil {
//		this.pool.Free(this)
//	}
//	return nil
//}

//// -----------------------------------------------------------------------------
