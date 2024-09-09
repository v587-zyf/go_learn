package http

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

func PostGet(w http.ResponseWriter, r *http.Request) {
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

	// must post
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// param
	paramMap := make(map[string]string)
	_ = json.NewDecoder(r.Body).Decode(&paramMap)

	param1Str := paramMap["param1"]

	fmt.Println(param1Str)
}

func PostSend() {
	sendStruct := struct {
		Param1 string `json:"param1"`
	}{
		"param1",
	}
	sendData, _ := json.Marshal(sendStruct)
	fmt.Println("sendData", string(sendData))
	resp, err := http.Post("http://127.0.0.1:8080", "application/json", bytes2.NewBuffer(sendData))
	if err != nil {
		fmt.Println("post err", err)
		return
	}
	defer resp.Body.Close()

	backBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("io.ReadAll err=", err)
		return
	}

	fmt.Println(backBytes)
}
