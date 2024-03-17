package backendapi

import (
	"backend/config"
	"backend/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func BetaFeatures(r *ghttp.Request) {
	ctx := r.GetCtx()
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
		"Authorization": "Bearer " + carinfo.AccessToken,
		"User-Agent":    r.Header.Get("User-Agent"),
	}).Post(ctx, config.CHATPROXY+"/backend-api/settings/beta_features?feature=plugins&value=true")
	if err != nil {
		g.Log().Error(ctx, err)
		r.Response.Status = 500
		r.Response.WriteJson(g.Map{
			"detail": err.Error(),
		})
		return
	}
	r.Response.Status = res.StatusCode
	r.Response.Write(res.ReadAllString())
}
