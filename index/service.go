package index

import "sync"

type IndexService interface {
	GetIndex(indexName string) (string, error)
}

type IndexServiceImpl struct {
}

var (
	indexService IndexService
	once         sync.Once
)

func NewIndexService() IndexService {
	once.Do(func() {
		indexService = &IndexServiceImpl{}
	})

	return indexService
}

func (i *IndexServiceImpl) GetIndex(indexName string) (string, error) {
	return "index", nil
}
