package utils

import (
	"math/rand"
	"sync"
)

type ISafeRand interface {
	Int63() int64
	Uint32() uint32
	Uint64() uint64
	Int31() int32
	Int() int
	Int63n(n int64) int64
	Int31n(n int32) int32 
	Intn(n int) int 
	Float64() float64
	NormFloat64() float64
	Float32() float32
	Perm(n int) []int
}

type lockedSource struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *lockedSource) Int63() int64 {
	defer r.lk.Unlock()
	r.lk.Lock()
	return r.src.Int63()
}

func (r *lockedSource) Uint64() uint64 {
	defer r.lk.Unlock()
	r.lk.Lock()
	return r.src.Uint64()
}

func (r *lockedSource) Seed(seed int64) {
	defer r.lk.Unlock()
	r.lk.Lock()
	r.src.Seed(seed)
}

type safeRand struct {
	rand *rand.Rand
}

func (r *safeRand) Int63() int64 { return r.rand.Int63() }
func (r *safeRand) Uint32() uint32 { return r.rand.Uint32() }
func (r *safeRand) Uint64() uint64 { return r.rand.Uint64() }
func (r *safeRand) Int31() int32 { return r.rand.Int31() }
func (r *safeRand) Int() int { return r.rand.Int() }
func (r *safeRand) Int63n(n int64) int64 { return r.rand.Int63n(n) }
func (r *safeRand) Int31n(n int32) int32 { return r.rand.Int31n(n) }
func (r *safeRand) Intn(n int) int { return r.rand.Intn(n) }
func (r *safeRand) Float64() float64 { return r.rand.Float64() }
func (r *safeRand) NormFloat64() float64 { return r.rand.NormFloat64() }
func (r *safeRand) Float32() float32 { return r.rand.Float32() }
func (r *safeRand) Perm(n int) []int { return r.rand.Perm(n) }

func NewSafeRand(seed int64) ISafeRand {
	src64, ok := rand.NewSource(seed).(rand.Source64)
	if !ok {
		panic("invalid source")
	}
	src := &lockedSource {
		src: src64,
	}
	return &safeRand{
		rand: rand.New(src),
	}
}