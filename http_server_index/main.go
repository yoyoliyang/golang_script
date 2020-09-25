// very simple http server for directory
// 一个非常简单的http方式目录访问脚本
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var path = flag.String("w", "./", "path")
var port = flag.String("p", "8008", "port")
var host = flag.String("h", "0.0.0.0", "host")

func main() {

	flag.Parse()
	h := *host + ":" + *port

	http.HandleFunc("/", handle)
	fmt.Println("starting http server:", h, *path)
	log.Fatal(http.ListenAndServe(h, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	hf := http.FileServer(http.Dir(*path))
	hf.ServeHTTP(w, r)
	log.Println(r.RemoteAddr, r.URL.Path)
}
