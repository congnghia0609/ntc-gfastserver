/**
 *
 * @author nghiatc
 * @since Oct 06, 2020
 */

package server

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Index handler
func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

// Hello handler
func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

// StartWebServer start WebServer
func StartWebServer(name string) {

	// Setup Router Handlers.
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)

	// Serve static files from the ./public directory
	r.NotFound = fasthttp.FSHandler("./public", 0) // http://localhost:8080/css/main.css

	log.Printf("======= WebServer[%s] is running on host: %s", name, ":8080")
	// log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
	err := fasthttp.ListenAndServe(":8080", r.Handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
