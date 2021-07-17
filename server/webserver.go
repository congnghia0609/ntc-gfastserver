/**
 *
 * @author nghiatc
 * @since Oct 06, 2020
 */

package server

import (
	"fmt"
	"github.com/congnghia0609/ntc-gfastserver/handler"
	"log"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Index handler
func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

// Hello handler
func Hello(ctx *fasthttp.RequestCtx) {
	// set some headers and status code first
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusOK)

	// then write the first part of body
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))

	// then set more headers
	ctx.Response.Header.Set("Foo-Bar", "baz")

	// then write more body
	fmt.Fprintf(ctx, "this is the second part of body\n")

	// // then override already written body
	// ctx.SetBody([]byte("this is completely new body contents"))
}

// StartWebServer start WebServer
func StartWebServer(name string) {
	// Config
	c := nconf.GetConfig()
	address := c.GetString(name + ".addr")

	// Setup Router Handlers.
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)
	// Post Handler
	r.GET("/post/{id}", handler.GetPost)
	r.GET("/posts", handler.GetPosts)
	r.POST("/post", handler.AddPost)
	r.PUT("/post", handler.UpdatePost)
	r.DELETE("/post/{id}", handler.DeletePost)
	// Tag Handler
	r.GET("/tag/{id}", handler.GetTag)
	r.GET("/tags", handler.GetTags)
	r.POST("/tag", handler.AddTag)
	r.PUT("/tag", handler.UpdateTag)
	r.DELETE("/tag/{id}", handler.DeleteTag)

	// Serve static files from the ./public directory
	r.NotFound = fasthttp.FSHandler("./public", 0) // http://localhost:8080/css/main.css

	log.Printf("======= WebServer[%s] is running on host: %s", name, address)
	// log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
	err := fasthttp.ListenAndServe(address, r.Handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
