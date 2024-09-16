package robot

import (
	"fmt"
)

func InputString(des string) string {
	for {
		out := ""
		//log.Info(des)
		fmt.Println(des)
		if _, err := fmt.Scanln(&out); err != nil {
			fmt.Println("input err,please try again")
			//log.Info("input err,please try again")
		}
		return out
	}
}

func InputInt32(des string) int32 {
	for {
		out := int32(0)
		//log.Info(des)
		fmt.Println(des)
		if _, err := fmt.Scanln(&out); err != nil {
			//log.Info("input err,please try again")
			fmt.Println("input err,please try again")
		}
		//time.Sleep(200 * time.Millisecond)
		return out
	}
}
