package mongo

import (
	"fmt"
	"time"
)

func LoadTime() {
	var id uint64 = 1
	d := &Test{
		ID:   id,
		Time: time.Unix(0, 0),
	}
	GetTestMongo().InsertOne(d)

	data, _ := GetTestMongo().LoadOne(id)
	fmt.Println(data.Time.IsZero())
}
