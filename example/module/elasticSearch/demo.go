package elasticSearch

func Do() {
	host := "http://127.0.0.1:9200"
	err := Init(host)
	if err != nil {
		panic(err)
	}

	//id := "1"
	indexName := "student"
	//Create(indexName, id)
	err = Search(indexName)
	if err != nil {
		panic(err)
	}
	err = Query(indexName)
	if err != nil {
		panic(err)
	}
}
