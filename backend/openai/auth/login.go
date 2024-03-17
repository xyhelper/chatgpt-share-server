package auth

import (
	"backend/config"
	"backend/utility"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/narqo/go-badge"
)

func Login(r *ghttp.Request) {
	ctx := r.GetCtx()
	method := r.Method
	if method == "GET" {
		req := r.GetMapStrStr()

		carid := req["carid"]
		// 	r.Response.WriteTpl("login.html")
		// 	return
		// }
		carInfo, err := utility.CheckCar(ctx, carid)
		if err != nil {
			g.Log().Error(ctx, err)
			badge, err := badge.RenderBytes("😭", "      翻车|不可用", "grey")
			if err != nil {
				g.Log().Error(ctx, err)
				r.Response.WriteTpl("login.html")
			}
			r.Response.WriteTpl("login.html", g.Map{"badge": string(badge)})

			return
		}
		var badgeSVG []byte

		count := utility.GetStatsInstance(carid).GetCallCount()
		expTime := cool.CacheManager.MustGetExpire(ctx, "clears_in:"+carid)
		teamExpTime := cool.CacheManager.MustGetExpire(ctx, "team_clears_in:"+carid)
		expInt := gconv.Int(expTime.Seconds())
		teamExpInt := gconv.Int(teamExpTime.Seconds())
		if expInt > 0 || teamExpInt > 0 {
			if expInt > 0 && teamExpInt > 0 {
				// 两者都有
				badgeSVG, err = badge.RenderBytes(carInfo.IsPlusStr, "            😡停运｜将于"+gconv.String(min(expInt, teamExpInt))+"秒后恢复", "red")
			}
			if expInt > 0 && teamExpInt == 0 {
				// 只有个人
				badgeSVG, err = badge.RenderBytes(carInfo.IsPlusStr, "            😡PLUS停运｜将于"+gconv.String(expInt)+"秒后恢复", "red")
			}
			if expInt == 0 && teamExpInt > 0 {
				// 只有团队
				badgeSVG, err = badge.RenderBytes(carInfo.IsPlusStr, "            😡TEAM停运｜将于"+gconv.String(teamExpInt)+"秒后恢复", "red")
			}

			// badgeSVG, err = badge.RenderBytes(carInfo.IsPlusStr, "            😡停运｜将于"+gconv.String(expInt)+"秒后恢复", "red")
		} else {
			if count > 20 {
				badgeSVG, err = badge.RenderBytes(carInfo.IsPlusStr, "    😅繁忙|可用", "yellow")
			} else {
				badgeSVG, err = badge.RenderBytes(carInfo.IsPlusStr, "    😊空闲|推荐", "green")
			}
		}

		if err != nil {
			g.Log().Error(ctx, err)
			r.Response.WriteTpl("login.html")
		}
		// fmt.Printf("%s", badge)

		r.Response.WriteTpl("login.html", g.Map{"badge": string(badgeSVG)})
		return
	} else {
		req := r.GetMapStrStr()
		loginVar := g.Client().PostVar(ctx, config.OauthUrl, req)
		loginJson := gjson.New(loginVar)
		// loginJson.Dump()
		code := loginJson.Get("code").Int()
		if code != 1 {
			msg := loginJson.Get("msg").String()
			r.Response.WriteTpl("login.html", g.Map{
				"error": msg,
				"carid": req["carid"],
			})
			return
		} else {
			r.Session.Set("usertoken", req["usertoken"])
			r.Session.Set("carid", req["carid"])
			r.Response.RedirectTo("/")
		}
	}
}

func LoginToken(r *ghttp.Request) {
	ctx := r.GetCtx()
	req := r.GetMapStrStr()
	resptype := req["resptype"]

	loginVar := g.Client().PostVar(ctx, config.OauthUrl, req)
	loginJson := gjson.New(loginVar)
	// loginJson.Dump()
	code := loginJson.Get("code").Int()
	if code != 1 {
		msg := loginJson.Get("msg").String()
		if resptype == "json" {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  msg,
			})
			return
		} else {
			r.Response.WriteTpl("login.html", g.Map{
				"error": msg,
				"carid": req["carid"],
			})
			return
		}
	} else {
		r.Session.Set("usertoken", req["usertoken"])
		r.Session.Set("carid", req["carid"])
		if resptype == "json" {
			r.Session.Set("usertoken", req["usertoken"])
			r.Session.Set("carid", req["carid"])
			r.Response.WriteJson(g.Map{
				"code": 1,
				"msg":  "登录成功",
			})
			return
		} else {
			r.Session.Set("usertoken", req["usertoken"])
			r.Session.Set("carid", req["carid"])
			r.Response.RedirectTo("/")
		}
	}
}

// 从两个整数中获取最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
