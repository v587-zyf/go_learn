package elasticSearch

import (
	"context"
	"fmt"
)

func Del(indexName, id string) error {
	_, err := GetEsClient().Delete().Index(indexName).Id(id).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
