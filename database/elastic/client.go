package elastic

import (
	"context"
	"errors"

	"github.com/JREAMLU/j-kit/constant"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

// Elastic Elastic client
type Elastic struct {
	client *elastic.Client
	Infos  []*elastic.PingResult
	Codes  []int
}

// NewElastic new elastic
func NewElastic(debug bool, urls []string) (*Elastic, error) {
	if len(urls) == constant.ZeroInt {
		return nil, errors.New(constant.ESUrlNotEmpty)
	}

	ctx := context.Background()
	options := []elastic.ClientOptionFunc{elastic.SetURL(urls...)}

	if debug {
		options = append(options, elastic.SetTraceLog(logrus.New()))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	var infos = make([]*elastic.PingResult, len(urls))
	var codes = make([]int, len(urls))

	for i, url := range urls {
		info, code, err := client.Ping(url).Do(ctx)
		if err != nil {
			return nil, err
		}
		infos[i] = info
		codes[i] = code
	}

	return &Elastic{
		client: client,
		Infos:  infos,
		Codes:  codes,
	}, nil
}
