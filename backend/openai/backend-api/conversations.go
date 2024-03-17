package backendapi

import (
	"backend/config"
	"backend/modules/chatgpt/model"
	"backend/utility"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ConvItem struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	CreateTIme string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func Conversations(r *ghttp.Request) {
	ctx := r.Context()
	// 获取当前登录用户信息
	usertoken := r.Session.MustGet("usertoken").String()
	if usertoken == "" {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"detail": "Authentication credentials were not provided.",
		})
		return
	}
	method := r.Method
	carid := r.Session.MustGet("carid").String()
	carinfo, err := utility.CheckCar(ctx, carid)

	if err != nil {
		g.Log().Error(ctx, err)
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"detail": "Authentication credentials were not provided.",
		})
		return
	}
	// 会话列表
	if method == "GET" {
		offset := r.Get("offset").Int()
		limit := r.Get("limit").Int()
		// total := 100
		modelquery := cool.DBM(model.NewChatgptConversations()).As("a").LeftJoin("chatgpt_session", "b", "a.email=b.email").Fields("a.createTime", "a.updateTime", "a.convid", "a.title")
		if config.DISALLOW_ROAM {
			modelquery.Where(g.Map{
				"a.usertoken":        usertoken,
				"a.email":            carinfo.Email,
				"a.chatgptaccountid": r.Header.Get("ChatGPT-Account-ID"),
			})

		} else {
			modelquery.Where(g.Map{
				"a.usertoken": usertoken,
				"b.status":    "1",
			})
		}
		items, total, err := modelquery.OrderDesc("a.updateTime").Limit(limit).Offset(offset).AllAndCount(false)
		if err != nil {
			g.Log().Error(ctx, err)
			r.Response.Status = 500
			r.Response.WriteJson(g.Map{
				"detail": err.Error(),
			})
			return
		}
		if items == nil {
			r.Response.WriteJson(g.Map{
				"total":                     total,
				"items":                     g.Slice{},
				"limit":                     limit,
				"offset":                    offset,
				"has_missing_conversations": false,
			})
			return
		}
		var convItems []ConvItem
		// 遍历items
		for _, item := range items {
			convItem := ConvItem{
				Id:         item["convid"].String(),
				Title:      item["title"].String(),
				CreateTIme: item["createTime"].String(),
				UpdateTime: item["updateTime"].String(),
			}
			convItems = append(convItems, convItem)
		}
		r.Response.WriteJson(g.Map{
			"total":                     total,
			"items":                     convItems,
			"limit":                     limit,
			"offset":                    offset,
			"has_missing_conversations": false,
		})
		return
	}

	// 清除所有会话
	if method == "PATCH" {
		cool.DBM(model.NewChatgptConversations()).Where(g.Map{
			"usertoken": usertoken,
			// "email":            carinfo.Email,
			// "chatgptaccountid": r.Header.Get("ChatGPT-Account-ID"),
		}).Delete()
		r.Response.WriteJson(g.Map{
			"success": true,
			"message": nil,
		})
		return
	}
}
