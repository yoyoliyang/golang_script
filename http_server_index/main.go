// very simple http server for directory
// 一个非常简单的http方式目录访问脚本，使用了内置的FileServer方式来索引目录
// 配合模板，使用javascript来修改主题和标题名称
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
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
	tmpl, err := template.New("index").Parse(js)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	err = tmpl.Execute(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	// fmt.Fprintln(w, js)

	log.Println(r.RemoteAddr, r.URL.Path)
}

var js = `
<style type="text/css">
body {
	color: #26b72b;
	margin-top: 10px;
}
 
pre {
	font-size: 20px;
	margin-left: 20px;
}
</style>

<script>
let pre = document.getElementsByTagName("pre")
pre[0].innerHTML = "<h3>Index of {{.URL.Path}}</h3>" + "<hr>" + pre[0].innerHTML
</script>
`
