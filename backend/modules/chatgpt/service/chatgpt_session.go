package service

import (
	"backend/config"
	"backend/modules/chatgpt/model"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

type ChatgptSessionService struct {
	*cool.Service
}

func NewChatgptSessionService() *ChatgptSessionService {
	return &ChatgptSessionService{
		&cool.Service{
			Model: model.NewChatgptSession(),
			UniqueKey: g.MapStrStr{
				"email": "邮箱不能重复",
				"carID": "车号不能重复",
			},
			NotNullKey: g.MapStrStr{
				"email":    "邮箱不能为空",
				"password": "密码不能为空",
				"carID":    "车号不能为空",
			},
			PageQueryOp: &cool.QueryOp{
				FieldEQ:      []string{"email", "password", "remark", "officialSession", "carID"},
				KeyWordField: []string{"email", "password", "remark", "officialSession", "carID"},
			},
		},
	}
}

// ModifyBefore 新增/删除/修改之前的操作
func (s *ChatgptSessionService) ModifyBefore(ctx g.Ctx, method string, param map[string]interface{}) (err error) {
	if method == "Delete" {
		ids := gjson.New(param["ids"]).Array()
		for _, id := range ids {
			record, err := cool.DBM(s.Model).Where("id=?", id).One()
			if err != nil {
				g.Log().Error(ctx, "ChatgptSessionService.ModifyBefore", "get record error", err)
				continue
			}
			carid := record["carID"].String()
			cool.CacheManager.Remove(ctx, "session:"+carid)
			email := record["email"].String()
			cool.CacheManager.Remove(ctx, "email:"+email)
		}
	}
	return
}

// ModifyAfter 新增/删除/修改之后的操作
func (s *ChatgptSessionService) ModifyAfter(ctx g.Ctx, method string, param map[string]interface{}) (err error) {

	g.Log().Debug(ctx, "ChatgptSessionService.ModifyAfter", method, param)
	// g.Dump(param)
	// 新增/修改 之后，更新session
	if method != "Add" && method != "Update" {
		return
	}
	officialSession := gjson.New(param["officialSession"])
	refreshCookie := officialSession.Get("refreshCookie").String()
	// 如果没有officialSession，就去获取
	getSessionUrl := config.CHATPROXY + "/getsession"
	sessionVar := g.Client().SetHeader("authkey", config.AUTHKEY).SetCookie("arkoseToken", gconv.String(param["arkoseToken"])).PostVar(ctx, getSessionUrl, g.Map{
		"username":      param["email"],
		"password":      param["password"],
		"authkey":       config.AUTHKEY,
		"refreshCookie": refreshCookie,
	})
	sessionJson := gjson.New(sessionVar)
	if sessionJson.Get("accessToken").String() == "" {
		g.Log().Error(ctx, "ChatgptSessionService.ModifyAfter", "get session error", sessionJson)
		detail := sessionJson.Get("detail").String()
		if detail != "" {
			err = gerror.New(detail)
			cool.DBM(s.Model).Where("carid=?", param["carID"]).Update(g.Map{
				"officialSession": sessionJson.String(),
				"status":          0,
			})
		} else {
			err = gerror.New("get session error")
		}
		return
	}
	email := sessionJson.Get("user.email").String()
	models := sessionJson.Get("models").Array()
	_, err = cool.DBM(s.Model).Where("carid=?", param["carID"]).Update(g.Map{
		"email":           email,
		"officialSession": sessionJson.String(),
		"isPlus":          len(models) > 1,
		"status":          1,
	})
	if err != nil {
		g.Log().Error(ctx, "ChatgptSessionService.ModifyAfter", "update session error", err)
		return
	}
	cool.CacheManager.Set(ctx, "session:"+gconv.String(param["carID"]), sessionJson.String(), 0)
	cool.CacheManager.Set(ctx, "email:"+email, gconv.String(param["carID"]), 0)
	return
}

func init() {
	ctx := gctx.GetInitCtx()
	sessionRecords, err := cool.DBM(NewChatgptSessionService().Model).All()
	if err != nil {
		panic(err)
	}
	for _, record := range sessionRecords {
		sessionJson := gjson.New(record["officialSession"])
		// sessionJson.Dump()
		email := sessionJson.Get("user.email").String()
		g.Log().Debug(ctx, "init session", record["carID"], email)
		if email != "" {
			cool.CacheManager.Set(ctx, "session:"+record["carID"].String(), record["officialSession"].String(), 0)
			cool.CacheManager.Set(ctx, "email:"+email, record["carID"].String(), 0)
		}
	}
}
