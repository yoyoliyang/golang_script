// 文件索引服务小脚本
// 根据请求的url path来访问对应的文件。
// 使用ioutil方式来遍历文件夹，生成文件夹索引目录。
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type file struct {
	Name string
	Url  string
	Size int64
}

type fileList []file

var indexTmpl = `
<a href="{{.Url}}">{{.Name}}</a>  {{.Size}}<br/>
`

var port = flag.String("p", "8000", "port")
var host = flag.String("b", "0.0.0.0", "host")
var dirpath = flag.String("d", ".", "dir path")

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	fmt.Printf("Listen and Serve on %v\n", *host+":"+*port)
	err := http.ListenAndServe(*host+":"+*port, mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	var p string
	if strings.HasSuffix(*dirpath, "/") && *dirpath != "/" {
		p = strings.TrimSuffix(*dirpath, "/")
	} else if *dirpath == "/" {
		p = r.URL.Path
	} else {
		p = *dirpath + r.URL.Path
	}
	log.Printf("%v\n", p)
	f, err := os.Open(p)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fInfo, err := f.Stat()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	if fInfo.IsDir() {
		dirs, err := ioutil.ReadDir(p)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		tmpl, err := template.New("index").Parse(indexTmpl)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		for _, dir := range dirs {
			var url string
			name := dir.Name()
			if r.URL.Path == "/" {
				url = r.URL.Path + dir.Name()
			} else {
				url = r.URL.Path + "/" + dir.Name()
			}
			if dir.IsDir() {
				name = dir.Name() + "/"
			}
			err := tmpl.Execute(w, file{name, url, dir.Size()})
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
		}
	} else {
		err := catFile(w, f)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}
}

func catFile(w http.ResponseWriter, f *os.File) error {
	_, err := io.Copy(w, f)
	return err
}
