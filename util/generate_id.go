package util

import (
	"cana.io/clap/pkg/base"
	"math/rand"
	"sync/atomic"
	"time"
)

var nextStart int32 = 1
var timeLast = base.TimeStart

func UniqueId() uint64 {
	nowTime := getTime()
	if nowTime < timeLast {
		println(nowTime)
		nowTime = timeLast
	} else {
		timeLast = nowTime
	}
	pastTime := nowTime - base.TimeStart
	getNext := atomic.AddInt32(&nextStart, rand.Int31n(9) + 1)
	if getNext > base.NextLimit {
		atomic.StoreInt32(&nextStart, 1)
	}
	return (pastTime << base.TimeBits) | (base.Seq() << base.SeqBits) | (uint64(getNext))
}

func getTime() uint64 {
	return uint64(time.Now().Unix())
}
