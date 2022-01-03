package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// const url = `https://apis.imooc.com/?cid=108&icode=%sinternal`
const url = `https://apis.imooc.com`

func main() {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", content)
}
