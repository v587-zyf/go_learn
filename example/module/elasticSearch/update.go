package elasticSearch

import (
	"context"
	"fmt"
)

func Update(indexName, id string, doc interface{}) error {
	_, err := GetEsClient().Update().Index(indexName).Id(id).Doc(doc).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
