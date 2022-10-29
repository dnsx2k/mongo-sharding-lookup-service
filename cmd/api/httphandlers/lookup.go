package httphandlers

import (
	"net/http"

	"github.com/dnsx2k/mongo-sharding-lookup-service/pkg/dto"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
)

type HTTPHandlerContext struct {
	cache *lru.Cache
}

// Setup - setup for HTTP gin handler
func (sc *HTTPHandlerContext) Setup(route gin.IRouter) {
	route.GET("lookup/:key", sc.handleGet)
	route.POST("lookup", sc.handlePost)
}

// New - factory function for HTTP lookup handler
func New(cache *lru.Cache) *HTTPHandlerContext {
	return &HTTPHandlerContext{cache: cache}
}

func (sc *HTTPHandlerContext) handleGet(gCtx *gin.Context) {
	key := gCtx.Param("key")
	val, ok := sc.cache.Get(key)
	if !ok {
		gCtx.Status(http.StatusNotFound)
		return
	}

	gCtx.JSON(http.StatusOK, dto.Lookup{
		Key:      key,
		Location: val.(string),
	})
}

func (sc *HTTPHandlerContext) handlePost(gCtx *gin.Context) {
	var newLookupEntry dto.LookupBatch
	if err := gCtx.BindJSON(&newLookupEntry); err != nil {
		gCtx.JSON(http.StatusBadRequest, err)
		return
	}

	//TODO: Save to DB

	gCtx.Status(http.StatusOK)
}
