package openai

import (
	"backend/modules/chatgpt/model"

	"github.com/cool-team-official/cool-admin-go/cool"
	baseservice "github.com/cool-team-official/cool-admin-go/modules/base/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

// 分页返回车辆列表
func CarPage(r *ghttp.Request) {
	// ctx := r.GetCtx()
	reqJson, err := r.GetJson()
	if err != nil {
		r.Response.WriteJson(g.Map{
			"detail": "Unable to parse request body.",
		})
		return
	}
	page := reqJson.Get("page").Int()
	size := reqJson.Get("size").Int()
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	// .Limit(size).Offset((page - 1) * size).AllAndCount(false)
	record, count, err := cool.DBM(model.NewChatgptSession()).Fields("carID", "status", "isPlus").OrderDesc("sort").OrderAsc("id").Page(page, size).AllAndCount(false)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"detail": err.Error(),
		})
		return
	}
	notice := baseservice.NewBaseSysParamService().HtmlByKey("notice")
	// 去除 <html><body> </body></html>
	notice = gstr.Replace(notice, "<html><body>", "", -1)
	notice = gstr.Replace(notice, "</body></html>", "", -1)
	// 如果 notice 为 keyName notfound
	if notice == "keyName notfound" {
		notice = "暂无公告"
	}
	r.Response.WriteJson(g.Map{
		"code":     1000,
		"messages": "success",
		"data": g.Map{
			"list": record,
			"pagination": g.Map{
				"page":  page,
				"size":  size,
				"total": count,
			},
		},
		"notice": notice,
	})
}
