package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// code: 1000  ->  验证成功
	// code: 1001  ->  icode不正确
	fmt.Println("handle...")
	fmt.Fprint(w, `{"msg":"Yes! I'm imooc.","data":null,"code":1000}`)
}

func main() {
	http.HandleFunc("/", handle)
	fmt.Println("Running...")
	if err := http.ListenAndServeTLS(":443", "./ssl/server.crt", "./ssl/server.key", nil); err != nil {
		// if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
