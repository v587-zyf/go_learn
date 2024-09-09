package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

func GetGet(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()

		if r := recover(); r != nil {
			stackBytes := debug.Stack()
			fmt.Println("r", r, "Error", stackBytes)
		}
	}()
	r.ParseForm()

	defer func() {
		// 自定义返回内容
		backData := struct {
			Code int
			Msg  string
		}{
			200,
			"success",
		}
		bytes, err := json.Marshal(backData)
		if err != nil {
			fmt.Fprint(w, nil)
			return
		}
		w.Write(bytes)
	}()

	// param
	param1 := r.Form.Get("param1")

	fmt.Println(param1)
}

func send() {
	param1 := "param1"
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080?param1=%s", param1))
	if err != nil {
		log.Println("http.Get err=", err)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll err=", err)
		return
	}

	fmt.Println(bytes)
}
