/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"ntc-gfastserver/mdb"
	"ntc-gfastserver/post"
	"strconv"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/valyala/fasthttp"
)

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

	dataResp := DataResp{Err: 0, Msg: "Add post successfully", Data: p}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// UpdatePost api update post
func UpdatePost(ctx *fasthttp.RequestCtx) {
	params := make(map[string]interface{})
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	var id int64 = 0
	if params["id"] != nil {
		id = int64(params["id"].(float64))
	}
	title := ""
	if params["title"] != nil {
		title = params["title"].(string)
	}
	body := ""
	if params["body"] != nil {
		body = params["body"].(string)
	}
	fmt.Println("id:", id)
	fmt.Println("title:", title)
	fmt.Println("body:", body)
	// Validate params
	if id == 0 || len(title) == 0 || len(body) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	// id, _ := strconv.ParseInt(sid, 10, 64)
	p := post.GetPost(id)
	if p.ID <= 0 {
		dataResp := DataResp{Err: -1, Msg: "Post is not exist"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}
	// id, _ := mdb.Next(post.TablePost)
	np := post.Post{
		ID:        id,
		Title:     title,
		Body:      body,
		CreatedAt: p.CreatedAt,
		UpdatedAt: time.Now(),
	}
	count, err1 := post.UpdatePost(np)
	if err1 != nil || count < 1 {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Update post fail"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	dataResp := DataResp{Err: 0, Msg: "Update post successfully", Data: np}
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

// GetPosts api get page post
func GetPosts(ctx *fasthttp.RequestCtx) {
	c := nconf.GetConfig()
	paging := c.GetInt64("system.paging")

	mapData := make(map[string]interface{})
	isMore := false
	posts := []post.Post{}
	var page int64 = 1
	pg, _ := ctx.QueryArgs().GetUint("page")
	if pg > 0 {
		page = int64(pg)
	}
	log.Println("page:", page)
	mapData["page"] = page

	total, _ := post.GetTotalPost()
	// finish soon.
	if total == 0 {
		mapData["isMore"] = false
		mapData["posts"] = posts
		dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: mapData}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	maxPage := (total-1)/paging + 1
	// finish soon.
	if page > maxPage {
		mapData["isMore"] = false
		mapData["posts"] = posts
		dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: mapData}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}
	// Get data paging.
	if page < maxPage {
		isMore = true
	}
	skip := (page - 1) * paging

	posts = post.GetSlidePost(skip, paging)
	mapData["isMore"] = isMore
	mapData["posts"] = posts
	dataResp := DataResp{Err: 0, Msg: "Get all posts successfully", Data: mapData}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// DeletePost api get post
func DeletePost(ctx *fasthttp.RequestCtx) {
	sid := ctx.UserValue("id").(string)
	id, _ := strconv.ParseInt(sid, 10, 64)
	count, err := post.DeletePost(id)
	if err != nil || count < 1 {
		dataResp := DataResp{Err: -1, Msg: "Delete post fail"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	dataResp := DataResp{Err: 0, Msg: "Delete post successfully"}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}
