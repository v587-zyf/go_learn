package elasticSearch

import (
	"context"
	"fmt"
)

func Get(indexName, id string) {
	res, err := GetEsClient().Get().Index(indexName).Id(id).Do(context.Background())
	if err != nil {
		fmt.Println("search err:", err)
	}
	fmt.Println("search source:", string(res.Source))
}
