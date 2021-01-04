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
	"ntc-gfastserver/tag"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddTag api add tag
func AddTag(ctx *fasthttp.RequestCtx) {
	params := make(map[string]interface{})
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	name := ""
	if params["name"] != nil {
		name = params["name"].(string)
	}
	fmt.Println("name:", name)
	// Validate params
	if len(name) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	id := primitive.NewObjectID()
	t := tag.Tag{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err1 := tag.InsertTag(t)
	if err1 != nil {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Add tag fail"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	dataResp := DataResp{Err: 0, Msg: "Add tag successfully", Data: t}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// UpdateTag api update tag
func UpdateTag(ctx *fasthttp.RequestCtx) {
	params := make(map[string]interface{})
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("param:", params)
	var id string = ""
	if params["id"] != nil {
		id = params["id"].(string)
	}
	name := ""
	if params["name"] != nil {
		name = params["name"].(string)
	}
	fmt.Println("id:", id)
	fmt.Println("name:", name)
	// Validate params
	if len(id) == 0 || len(name) == 0 {
		dataResp := DataResp{Err: -1, Msg: "Parameters invalid"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	t := tag.GetTag(oid)
	if t == nil {
		dataResp := DataResp{Err: -1, Msg: "Tag is not exist"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}
	nt := tag.Tag{
		ID:        oid,
		Name:      name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: time.Now(),
	}
	count, err1 := tag.UpdateTag(nt)
	if err1 != nil || count < 1 {
		fmt.Println("err1:", err1)
		dataResp := DataResp{Err: -1, Msg: "Update tag fail"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	dataResp := DataResp{Err: 0, Msg: "Update tag successfully", Data: nt}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// GetTag api get tag
func GetTag(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dataResp := DataResp{Err: -1, Msg: "TagId invalid"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}
	t := tag.GetTag(oid)
	if t == nil {
		dataResp := DataResp{Err: -1, Msg: "Tag is not exist"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	dataResp := DataResp{Err: 0, Msg: "Get tag successfully", Data: t}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// GetAllTags api get all tag
func GetAllTags(ctx *fasthttp.RequestCtx) {
	tags := tag.GetAllTag()
	dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: tags}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// GetTags api get page tag
func GetTags(ctx *fasthttp.RequestCtx) {
	c := nconf.GetConfig()
	paging := c.GetInt64("system.paging")

	mapData := make(map[string]interface{})
	isMore := false
	tags := []tag.Tag{}
	var page int64 = 1
	pg, _ := ctx.QueryArgs().GetUint("page")
	if pg > 0 {
		page = int64(pg)
	}
	log.Println("page:", page)
	mapData["page"] = page

	total, _ := tag.GetTotalTag()
	// finish soon.
	if total == 0 {
		mapData["isMore"] = false
		mapData["tags"] = tags
		dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: mapData}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	maxPage := (total-1)/paging + 1
	// finish soon.
	if page > maxPage {
		mapData["isMore"] = false
		mapData["tags"] = tags
		dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: mapData}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}
	// Get data paging.
	if page < maxPage {
		isMore = true
	}
	skip := (page - 1) * paging

	tags = tag.GetSlideTag(skip, paging)
	mapData["isMore"] = isMore
	mapData["tags"] = tags
	dataResp := DataResp{Err: 0, Msg: "Get all tags successfully", Data: mapData}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}

// DeleteTag api get tag
func DeleteTag(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dataResp := DataResp{Err: -1, Msg: "TagId invalid"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}
	count, err := tag.DeleteTag(oid)
	if err != nil || count < 1 {
		dataResp := DataResp{Err: -1, Msg: "Delete tag fail"}
		resp, _ := json.Marshal(dataResp)
		printJSON(ctx, string(resp))
		return
	}

	dataResp := DataResp{Err: 0, Msg: "Delete tag successfully"}
	resp, _ := json.Marshal(dataResp)
	printJSON(ctx, string(resp))
}
