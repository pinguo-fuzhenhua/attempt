package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// func main() {

// Create our Temp File
// tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
// if err != nil {
//     log.Fatal("Cannot create temporary file", err)
// }

// fmt.Println("Created File: " + tmpFile.Name())

// // Example writing to the file
// _, err = tmpFile.Write([]byte("This is a www.twle.cn example!"))
// if err != nil {
//     log.Fatal("Failed to write to temporary file", err)
// }

// // Remember to clean up the file afterwards
// defer os.Remove(tmpFile.Name())
// a := "https://cloud-image.c360dn.com/"
// fmt.Println(a[3:4])
// wb()
// }

func wb() {
	// var b bytes.Buffer
	// for i := 0; i < 10; i++ {
	// 	b.WriteString(strings.Join([]string{"i", "\n"}, ""))
	// }
	// fmt.Println(b.String())
	msg := ""
	a := []string{"111", "222", "333"}
	for _, v := range a {
		msg += v + "\n"
	}
	fmt.Println(msg)
}

func main() {
	c1 := make(chan os.Signal, 1)
	c2 := make(chan os.Signal, 1)

	signal.Notify(c1, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(c2, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c1
		fmt.Println("c1:", s)
	}()

	s := <-c2
	fmt.Println("c2", s)
	time.Sleep(1 * time.Second)
}
