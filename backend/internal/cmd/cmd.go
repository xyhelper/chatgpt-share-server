package cmd

import (
	"context"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gsession"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			if cool.IsRedisMode {
				go cool.ListenFunc(ctx)
				redisConfig := &gredis.Config{}
				redisVar, err := g.Cfg().Get(ctx, "redis.cool")
				if err != nil {
					g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
					panic(err)
				}
				if !redisVar.IsEmpty() {
					redisVar.Struct(redisConfig)
					redis, err := gredis.New(redisConfig)
					if err != nil {
						panic(err)
					}
					s.SetSessionStorage(gsession.NewStorageRedis(redis, "gfsession:"))
					g.Log().Debug(ctx, "session storage is redis")
				}
			}

			s.Run()
			return nil
		},
	}
)
