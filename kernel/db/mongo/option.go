package mongo

type MongoOption struct {
	Uri string
	DB  string
}

type Option func(o *MongoOption)

func NewMongoOption() *MongoOption {
	return &MongoOption{}
}

func WithUri(uri string) Option {
	return func(o *MongoOption) {
		o.Uri = uri
	}
}

func WithDb(db string) Option {
	return func(o *MongoOption) {
		o.DB = db
	}
}
