package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type FrameInfo struct {
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Format string `json:"format"`
	Offset int64  `json:"offset"`
}

func (f *FrameInfo) toMap() map[string]string {
	tframe := reflect.TypeOf(f)
	vframe := reflect.ValueOf(f)
	resMap := make(map[string]string)
	for i := 0; i < vframe.NumField(); i++ {
		fieldVal := vframe.Field(i)
		tagContent := tframe.Field(i).Tag.Get("json")
		valKind := fieldVal.Kind()
		switch valKind {
		case reflect.String:
			resMap[tagContent] = fieldVal.String()
		default:
			resMap[tagContent] = fieldVal.String()
		}
	}
	return resMap
}

func main() {
	res, err := http.Get("http://www.baidu.com")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	gByte, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(gByte))
	fmt.Println(res.Header)
}
