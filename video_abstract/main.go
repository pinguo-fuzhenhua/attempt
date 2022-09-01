package main

import "fmt"

type (
	FF func()
	CC interface {
		Do()
	}
)

func (f FF) Do() {
	f()
}

type Hb struct {
	JJ *Ha `json:"jj"`
}

type Ha struct {
	Con bool `json:"con"`
}

type AA struct {
	A []string
	B map[string]int
	C string
}

func main() {
	a := new(AA)
	b := a.C
	fmt.Println(b)
}
