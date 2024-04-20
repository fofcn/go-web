package index

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func InitRouter(public *gin.RouterGroup) {
	indexRouter := NewIndexRouter(NewIndexService())
	public.GET("/", indexRouter.Index)
}

type IndexRouter struct {
	is IndexService
}

var (
	ir     *IndexRouter
	irOnce sync.Once
)

func NewIndexRouter(is IndexService) *IndexRouter {
	irOnce.Do(func() {
		ir = &IndexRouter{
			is: is,
		}
	})

	return ir
}

func (i *IndexRouter) Index(c *gin.Context) {
	c.JSON(200, "Hello go-web")
}
