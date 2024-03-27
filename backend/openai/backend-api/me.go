package backendapi

import (
	"backend/config"
	"backend/utility"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Me(r *ghttp.Request) {
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
	res, err := g.Client().SetHeaderMap(g.MapStrStr{
		"Accept":          "*.*",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Authorization":   "Bearer " + carinfo.AccessToken,
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"User-Agent":      r.Header.Get("User-Agent"),
	}).Get(ctx, config.CHATPROXY+"/backend-api/me")
	if err != nil {
		g.Log().Error(ctx, err)
		r.Response.Status = 500
		r.Response.WriteJson(g.Map{
			"detail": err.Error(),
		})
		return
	}
	// r.Response.Status = res.StatusCode
	resStr := res.ReadAllString()
	if res.StatusCode != http.StatusOK {
		r.Response.Status = res.StatusCode
		r.Response.Write(resStr)
		return
	}
	resJson := gjson.New(resStr)
	// resJson.Set("id", "user-"+usertoken)
	resJson.Set("email", carid)
	resJson.Set("name", carid)
	resJson.Set("picture", "/avatars.png")
	resJson.Set("phone_number", "+1911")
	resJson.Set("orgs.data.0.description", "OpenAI")
	// resJson.Set("orgs.data.0.name", "OpenAI")
	// resJson.Set("orgs.data.0.id", "openai")
	// resJson.Dump()
	r.Response.WriteJson(resJson)
}
