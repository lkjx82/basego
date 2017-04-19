package id

import (
	"sync"
	"time"
)

/*
	Snowflake like id Generator

  	<41it: millsec, 10bit:server Id, 12bit:count>
*/

const (
	IdGenTimeMask   = int64(0x7FFFFFFFFFC00000)
	IdGenServMask   = int64(0x00000000003FF000)
	IdGenCntMask    = int64(0x0000000000000FFF)
	IdGenTimeOffset = 22
	IdGenServOffset = 12
)

// -----------------------------------------------------------------------------
//
type IdGen struct {
	servId        int64
	cnt           int16
	lastTimeMillS int64
	// beginTime     time.Time
	m sync.Mutex
}

// -----------------------------------------------------------------------------
// servId: server's id
func NewIdGen(servId int16) *IdGen {
	if uint(servId) >= 1024 {
		return nil
	}
	ig := IdGen{}
	ig.servId = int64(servId)
	return &ig
}

// -----------------------------------------------------------------------------
// gen id with lock for thread safe
func (this *IdGen) GenId() int64 {
	this.m.Lock()
	defer this.m.Unlock()
	id := int64(0)
	for {
		id = time.Now().UnixNano() / 10000000
		// if overflow, sleep 100 millsec
		if (int64(this.cnt)&IdGenCntMask == 0) && id == this.lastTimeMillS {
			<-time.After(time.Microsecond * 100)
		} else {
			break
		}
	}

	this.lastTimeMillS = id
	id = (id << IdGenTimeOffset) | ((this.servId << IdGenServOffset) & IdGenServMask) | (int64(this.cnt) & IdGenCntMask)
	this.cnt = this.cnt + 1
	return id
}

// -----------------------------------------------------------------------------
// gen id without lock
func (this *IdGen) GenIdUnsafe() int64 {
	id := int64(0)
	for {
		id = time.Now().UnixNano() / 10000000
		// if overflow, sleep 100 millsec
		if (int64(this.cnt)&IdGenCntMask == 0) && id == this.lastTimeMillS {
			<-time.After(time.Microsecond * 100)
		} else {
			break
		}
	}

	this.lastTimeMillS = id
	id = (id << IdGenTimeOffset) | ((this.servId << IdGenServOffset) & IdGenServMask) | (int64(this.cnt) & IdGenCntMask)
	this.cnt = this.cnt + 1
	return id
}

// -----------------------------------------------------------------------------
// get server Id from id
func Id2ServId(id int64) int16 {
	return int16((id & IdGenServMask) >> IdGenServOffset)
}

// -----------------------------------------------------------------------------
