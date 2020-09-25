// very simple http server for directory
// 一个非常简单的http方式目录访问脚本
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var port = flag.String("p", "8008", "port")
var host = flag.String("h", "0.0.0.0", "host")

var path string

func main() {
	flag.Parse()

	if flag.Arg(0) == "" {
		path = "./"
	} else {
		path = flag.Arg(0)
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("error opening file: err:", err)
			return
		}
		defer f.Close()
	}

	h := *host + ":" + *port

	http.HandleFunc("/", handle)
	fmt.Println("Listening on port:", h, path)
	log.Fatal(http.ListenAndServe(h, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	hf := http.FileServer(http.Dir(path))
	hf.ServeHTTP(w, r)
	log.Println(r.RemoteAddr, r.URL.Path)
}
