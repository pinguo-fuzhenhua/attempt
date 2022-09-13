package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {

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
	a := "https://cloud-image.c360dn.com/"
	fmt.Println(a[3:4])
}

func wb() {
	var b bytes.Buffer
	for i := 0; i < 10; i++ {
		b.WriteString(strings.Join([]string{"i", "\n"}, ""))
	}
	fmt.Println(b.String())
}
