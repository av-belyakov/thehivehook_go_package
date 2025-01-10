package elasticsearchapi

import "github.com/elastic/go-elasticsearch/v8"

type Settings struct {
	Port    int
	Host    string
	User    string
	Passwd  string
	IndexDB string
}

type ElasticsearchDB struct {
	settings Settings
	client   *elasticsearch.Client
}
