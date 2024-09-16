package redis

import (
	"strconv"
	"strings"
)

func ParseGateData(src string) (string, int32) {
	data := strings.Split(src, "_")
	if len(data) != 2 {
		return "", 0
	}

	id, err := strconv.Atoi(data[1])
	if err != nil {
		return "", 0
	}

	//addr := strings.Replace(data[0], "Thunder", "", 1)
	return data[0], int32(id)
}
