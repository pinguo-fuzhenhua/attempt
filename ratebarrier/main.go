package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

type RateBarrier struct {
	source []int
	op     uint64
	base   int
}

func NewRateBarrier(base int) *RateBarrier {
	source := make([]int, base, base)
	for i := 0; i < base; i++ {
		source[i] = i
	}

	// 随机排序
	rand.Shuffle(base, func(i, j int) {
		source[i], source[j] = source[j], source[i]
	})

	return &RateBarrier{
		source: source,
		base:   base,
	}
}

func (b *RateBarrier) Rate() int {
	return b.source[int(atomic.AddUint64(&b.op, 1))%b.base]
}

func main() {
	var wg sync.WaitGroup
	times := 200
	wg.Add(times)
	a, b, c := 0, 0, 0
	// 2:3:5
	base := NewRateBarrier(10)
	for i := 0; i < times; i++ {
		go func() {
			rate := base.Rate()
			switch {
			case rate < 6:
				a++
				fmt.Println("this is on 20%")
			case rate >= 6:
				b++
				fmt.Println("this is on 30%")
			}

			wg.Done()
		}()
	}
	wg.Wait()
	println(a, b, c)
}
