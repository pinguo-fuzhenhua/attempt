package main

import (
	"context"
	"fmt"
	"time"
)

type (
	cctx struct{}
	gctx struct{}
)

func Decay() {

}

func Stop(cancel context.CancelFunc) {
	time.Sleep(time.Second)
	cancel()
	fmt.Println("cancel")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	ch := make(chan struct{}, 1)
	go Stop(cancel)
	select {
	case <-ctx.Done():
		fmt.Println("timeout")
		return
	case <-ch:
		fmt.Println("end")
	}
}
