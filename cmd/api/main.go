package main

import (
	"net/http"
	"time"

	"github.com/dnsx2k/mongo-sharding-lookup-service/cmd/api/httphandlers"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
)

func main() {
	//TODO: get cache size from config
	cache, err := lru.New(10_000)
	if err != nil {
		panic(err)
	}

	lookupHttpHandler := httphandlers.New(cache)
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
