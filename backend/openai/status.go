package openai

import (
	"backend/utility"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

func Status(r *ghttp.Request) {
	ctx := r.GetCtx()
	carid := r.Get("carid").String()
	if carid == "" {
		r.Response.WriteJson(g.Map{
			"detail": "Carid is empty.",
		})
		return
	}
	carInfo, err := utility.CheckCar(ctx, carid)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"accountReady": false,
		})
		return
	}
	expTime := cool.CacheManager.MustGetExpire(ctx, "clears_in:"+carid)
	teamExpTime := cool.CacheManager.MustGetExpire(ctx, "team_clears_in:"+carid)
	expInt := gconv.Int(expTime.Seconds())
	teamExpInt := gconv.Int(teamExpTime.Seconds())
	count := utility.GetStatsInstance(carid).GetCallCount()
	r.Response.WriteJson(g.Map{
		"isPlus":         carInfo.IsPlus,
		"accountReady":   true,
		"clears_in":      expInt,
		"team_clears_in": teamExpInt,
		"count":          count,
	})

}

func EndPoint(r *ghttp.Request) {
	ctx := r.GetCtx()
	color := "grey"
	labelColor := "grey"
	lable := ""
	message := ""
	clears_in := 0
	team_clears_in := 0

	carid := r.Get("carid").String()
	if carid == "" {
		r.Response.WriteJson(g.Map{
			"detail": "Carid is empty.",
		})
		return
	}
	carInfo, err := utility.CheckCar(ctx, carid)
	if err != nil {
		color = "grey"
		labelColor = "grey"
		lable = ""
		message = "翻车|不可用"
		r.Response.WriteJson(g.Map{
			"schemaVersion": 1,
			"label":         lable,
			"message":       message,
			"color":         color,
			"labelColor":    labelColor,
			"namedLogo":     "Cockroach Labs",
		})
		return
	}
	expTime := cool.CacheManager.MustGetExpire(ctx, "clears_in:"+carid)
	teamExpTime := cool.CacheManager.MustGetExpire(ctx, "team_clears_in:"+carid)
	expInt := gconv.Int(expTime.Seconds())
	teamExpInt := gconv.Int(teamExpTime.Seconds())
	if expInt > 0 {
		clears_in = expInt
	}
	if teamExpInt > 0 {
		team_clears_in = teamExpInt
	}

	count := utility.GetStatsInstance(carid).GetCallCount()

	if carInfo.IsPlus {
		labelColor = "purple"
		lable = "PLUS"
	} else {
		lable = "3.5"
		labelColor = "blue"
	}
	if clears_in > 0 || team_clears_in > 0 {
		if clears_in > 0 && team_clears_in > 0 {
			// PLUS及team同时停运
			color = "red"
			// 获取最小的时间
			message = "停运｜将于" + gconv.String(min(clears_in, team_clears_in)) + "秒后恢复"
		}
		if clears_in > 0 && team_clears_in == 0 {
			// PLUS停运
			color = "red"
			message = "PLUS停运｜将于" + gconv.String(clears_in) + "秒后恢复"
		}
		if clears_in == 0 && team_clears_in > 0 {
			// team停运
			color = "red"
			message = "TEAM停运｜将于" + gconv.String(team_clears_in) + "秒后恢复"
		}
	} else {
		if count < 20 {
			color = "green"
			message = "空闲|推荐"
		} else {
			color = "yellow"
			message = "繁忙|可用"
		}
	}
	r.Response.WriteJson(g.Map{
		"schemaVersion": 1,
		"label":         lable,
		"message":       message,
		"color":         color,
		"labelColor":    labelColor,
		"namedLogo":     "Cockroach Labs",
	})
}

// 从两个整数中获取最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
