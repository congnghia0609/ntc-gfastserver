/**
 *
 * @author nghiatc
 * @since Oct 06, 2020
 */

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"ntc-gfastserver/mdb"
	"ntc-gfastserver/post"
	"strconv"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// DataResp is struct data response
type DataResp struct {
	Err  int         `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

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

func printJSON(ctx *fasthttp.RequestCtx, json string) {
	ctx.Response.Header.Set("content-type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(json)
}

// AddPost api add post
func AddPost(ctx *fasthttp.RequestCtx) {
	params := make(map[string]interface{})
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	title := ""
	if params["title"] != nil {
		title = params["title"].(string)
	}
	body := ""
	if params["body"] != nil {
		body = params["body"].(string)
	}
	fmt.Println("title:", title)
	fmt.Println("body:", body)
	// Validate params
	if len(title) == 0 || len(body) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	id, _ := mdb.Next(post.TablePost)
	p := post.Post{
		ID:        id,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err1 := post.InsertPost(p)
	if err1 != nil {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Add post fail"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	bp, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bp))

	dataResp := DataResp{Err: 0, Msg: "Add post successfully", Data: p}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// GetPost api get post
func GetPost(ctx *fasthttp.RequestCtx) {
	sid := ctx.UserValue("id").(string)
	id, _ := strconv.ParseInt(sid, 10, 64)
	p := post.GetPost(id)
	if p.ID <= 0 {
		dataResp := DataResp{Err: -1, Msg: "Post is not exist"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	bp, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bp))

	dataResp := DataResp{Err: 0, Msg: "Get post successfully", Data: p}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// GetAllPosts api get all post
func GetAllPosts(ctx *fasthttp.RequestCtx) {
	posts := post.GetAllPost()

	dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: posts}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
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
	r.GET("/post/{id}", GetPost)
	r.GET("/posts", GetAllPosts)
	r.POST("/post", AddPost)

	// Serve static files from the ./public directory
	r.NotFound = fasthttp.FSHandler("./public", 0) // http://localhost:8080/css/main.css

	log.Printf("======= WebServer[%s] is running on host: %s", name, address)
	// log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
	err := fasthttp.ListenAndServe(address, r.Handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
