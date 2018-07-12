package elastic

import (
	"context"
	"errors"

	"github.com/JREAMLU/j-kit/constant"
	"github.com/olivere/elastic"
)

// Elastic Elastic client
type Elastic struct {
	client *elastic.Client
	Index  string
	Sort   string
	Info   *elastic.PingResult
	Code   int
}

// NewElastic new elastic
func NewElastic(url string) (*Elastic, error) {
	if url == constant.EmptyStr {
		return nil, errors.New(constant.ESUrlNotEmpty)
	}

	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &Elastic{
		client: client,
		Info:   info,
		Code:   code,
	}, nil
}
