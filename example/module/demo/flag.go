package example

import (
	"flag"
	"fmt"
)

var (
	data string
)

func Do() {
	flag.StringVar(&data, "data", "", "")

	flag.Parse()

	fmt.Println(data)
	fmt.Println("data is", data)
}
