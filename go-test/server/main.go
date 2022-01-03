package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("####")
	fmt.Fprint(w, `{"msg":"返回什么就返回什么,比如：fuck you imooc","data":null,"code":1000}`)
}

func main() {
	http.HandleFunc("/", handle)
	if err := http.ListenAndServeTLS(":443", "/home/hts0000/ssl/server.crt", "/home/hts0000/ssl/server.key", nil); err != nil {
		// if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
