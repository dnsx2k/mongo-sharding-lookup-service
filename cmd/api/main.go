package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dnsx2k/mongo-sharding-lookup-service/cmd/api/httphandlers"
	"github.com/dnsx2k/mongo-sharding-lookup-service/pkg/dto"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := mongoClient.Database("lookup").Collection("lookup")

	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		panic(err)
	}
	var lookups []dto.Lookup
	if err := cursor.All(ctx, &lookups); err != nil {
		panic(err)
	}

	// Load lookups to cache
	cache, err := lru.New(len(lookups) + 5000)
	if err != nil {
		panic(err)
	}
	for i := range lookups {
		cache.Add(lookups[i].Key, lookups[i].Location)
	}

	lookupHttpHandler := httphandlers.New(cache, collection)
	router := gin.Default()
	apiV1 := router.Group("v1")
	lookupHttpHandler.Setup(apiV1)

	server := http.Server{
		Addr:         "localhost:8085",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}
