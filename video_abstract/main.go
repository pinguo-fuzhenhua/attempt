package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

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
	// a := new(AA)
	// b := a.C
	// fmt.Println(b)
	// cWluaXUtZGV2ZWxvcGVyOm1remlwNC56aXA=
	s := "https://cloud-image.c360dn.com/FiA2iId8UUT9vcSsXl1dmVLYPbL2"
	fmt.Println(unsafe.Sizeof(s))
	str1 := base64.StdEncoding.EncodeToString([]byte("qiniu-developer:mkzip4.zip"))
	fmt.Println(str1)
	a := strconv.Itoa(int(time.Now().Unix())) + "-index.txt"
	fmt.Println(a)
}
