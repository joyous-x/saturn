package common

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync/atomic"
	"time"
)

var seed int64 = 0

func init() {
	s0 := time.Now().UnixNano()
	buf := make([]byte, 8)
	if _, err := crand.Read(buf); err == nil {
		s0 ^= int64(binary.BigEndian.Uint64(buf))
	}

	atomic.StoreInt64(&seed, s0)
}

func NewRand() *rand.Rand {
	//r := rand.New(rand.NewSource(atomic.LoadInt64(&seed)))
	//atomic.StoreInt64(&seed, r.Int63())
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
