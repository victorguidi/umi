# うみ　ー　Umi

Umi is my simple revision on Golang's `net/http` package. I like the idea of using as much as I can from the `std` library in Golang, but I do miss features from other frameworks like [Echo](https://echo.labstack.com/) that allow me to return errors from my routes.

## How Umi works
Umi uses only the `std` library. All it does is instantiate it's own `Context` holding `http.ResponseWriter` and `*http.Request` and returning in a custom `HandlerFunc(c *Context) error`

## Features
1. Return errors from your route:
``` go
// here post takes Context and the return type is error
func post(c *umi.Context) error {
	var response any
	err := json.NewDecoder(c.Body).Decode(&response)
	if err != nil {
		return err
	}
	return c.JSON(response)
}

func main() {
	server := umi.New()
	server.POST("/awesome", post)
	server.Start(":8000")

}
```
2. Easy middleware implementation. Write your middleware, call with `func (u *Umi) Use(middleware Middleware)` before the routes that you wish to add the middleware to it.
``` go

// Simple endpoint
func get(c *umi.Context) error {
	return c.JSON(map[string]string{"foo":"bar"})
}

// Simple endpoint
func getOne(c *umi.Context) error {
	return c.JSON(map[string]string{"foo":"bar"})
}

// Custom middleware
func printer() middleware.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("My custom middleware")
			next(w, r)
		}
	}
}

func main() {
	server := umi.New()
	server.GET("/no_middleware", getOne)
	server.Use(printer())
	server.GET("/with_middleware", get)
}
```
3. Easy wrapper around `json` and `error`.
``` go
func post(c *umi.Context) error {
	var response any
	err := json.NewDecoder(c.Body).Decode(&response)
	if err != nil {
    // Easy way to customize the status code and the error your are returning
		return c.FAIL(err, 500)
	}
  // JSON already set the header and deals with the encoding
	return c.JSON(response)
}
```
4. Cors Implementation out of the box. Just use `WithFlexibleCors()`
``` go
func main() {
	server := umi.New().WithFlexibleCors()
	server.GET("/no_middleware", getOne)
	server.Use(printer())
	server.GET("/with_middleware", get)
}
```

# DISCLAIMER
I honestly don't recommend using this any serious production environment as its definitely not battle tested. I wrote as an exercise to help me understand better the `net/http` package and I use in personal projects only. Please use with care.
