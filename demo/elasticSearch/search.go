package elasticSearch

import (
	"context"
	"fmt"
	"reflect"
)

func Search(indexName string) error {
	res, err := GetEsClient().
		Search(indexName).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	var typ Student
	for _, item := range res.Each(reflect.TypeOf(typ)) {
		fmt.Println(item) // {小学生 7 [soccer basketball tennis]}
	}

	//for _, v := range res.Hits.Hits {
	//	fmt.Println(v)
	//}

	return err
}
