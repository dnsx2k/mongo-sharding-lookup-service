package httphandlers

import (
	"context"
	"net/http"

	"github.com/dnsx2k/mongo-sharding-lookup-service/pkg/dto"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPHandlerContext struct {
	cache      *lru.Cache
	collection *mongo.Collection
}

// Setup - setup for HTTP gin handler
func (sc *HTTPHandlerContext) Setup(route gin.IRouter) {
	route.GET("lookup/:key", sc.handleGet)
	route.POST("lookup", sc.handlePost)
}

// New - factory function for HTTP lookup handler
func New(cache *lru.Cache, collection *mongo.Collection) *HTTPHandlerContext {
	return &HTTPHandlerContext{cache: cache, collection: collection}
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

	persistIn := make([]interface{}, 0)
	for i := range newLookupEntry.Keys {
		persistIn = append(persistIn, map[string]string{
			"key":      newLookupEntry.Keys[i],
			"location": newLookupEntry.Location,
		})

		sc.cache.Add(newLookupEntry.Keys[i], newLookupEntry.Location)
	}

	_, err := sc.collection.InsertMany(context.Background(), persistIn, nil)
	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, err)
		return
	}

	gCtx.Status(http.StatusOK)
}
