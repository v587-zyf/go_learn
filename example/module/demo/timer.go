package example

import (
	"fmt"
	"time"
)

func Timer() {
	fn := func() {
		fmt.Println("----------")
	}

	var t *time.Timer
	tt := 1699804800
	ttt := time.Unix(int64(tt), 0)
	if time.Now().Before(ttt) {
		t = time.AfterFunc(1*time.Second, fn)
		fmt.Println(1)
	} else {
		t = time.AfterFunc(2*time.Second, fn)
		fmt.Println(2)
	}
	t.Stop()

	//fmt.Println(t)/*
	//t = time.AfterFunc(10*time.Second, fn)
	//fmt.Println("1")
	//t.Stop()
	//
	//fmt.Println("2")
	//t = time.AfterFunc(-1*time.Second, fn)
	//t.Reset(-10)
	//fmt.Println("3")
	//fmt.Println("4")
	//fmt.Println("5")*/

	//t = time.AfterFunc(0*time.Second, fn)
	//
	//t = time.AfterFunc(3*time.Second, fn)
	//
	//t = time.AfterFunc(10*time.Second, fn)
}
