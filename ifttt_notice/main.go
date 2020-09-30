package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	eventName = os.Getenv("IFTTT_EVENT_NAME")
	key       = os.Getenv("IFTTT_KEY")
)

var api = fmt.Sprintf("https://maker.ifttt.com/trigger/%v/with/key/%v", eventName, key)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "message")
		return
	}

	buf := &bytes.Buffer{}
	msg := fmt.Sprintf("{\"value1\":%q}", os.Args[1])
	_, err := buf.WriteString(msg)
	if err != nil && err == io.EOF {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(buf.String())

	if eventName != "" && key != "" {
		resp, err := http.Post(api, "application/json", buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer resp.Body.Close()

		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Fprintln(os.Stdout, string(result))
		return
	}
	fmt.Fprintln(os.Stderr, "Failed to get IFTTT event name and key")

}
