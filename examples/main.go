package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/victorguidi/umi"
	"github.com/victorguidi/umi/middleware"
)

func main() {
	server := umi.New().WithFlexibleCors()
	server.GET("/", get)
	server.Use(printer())
	server.POST("/oi", post)
	server.Start(":8000")
}

func printer() middleware.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("My custom middleware")
			next(w, r)
		}
	}
}

func get(c *umi.Context) error {
	return c.JSON(map[string]string{"foo": "bar"})
}

func post(c *umi.Context) error {
	var response any
	err := json.NewDecoder(c.Body).Decode(&response)
	if err != nil {
		return c.FAIL(err, 500)
	}
	return c.JSON(response)
}
