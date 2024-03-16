package openai

import (
	_ "backend/openai/api"
	_ "backend/openai/auth"
	_ "backend/openai/backend-api"
	_ "backend/openai/fileserver"
	_ "backend/openai/next"
	_ "backend/openai/public-api"
	"backend/openai/workspace"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	s := g.Server()
	// 根路由
	group := s.Group("/")
	group.GET("/", Index)
	group.GET("/c/:convId", C)
	group.GET("/g/:gizmoId", G)
	group.GET("/gpts/discovery", Discovery)
	group.GET("/gpts/editor", Editor)
	group.GET("/gpts/editor/:slug", Slug)
	group.GET("/g/:gizmoId/c/:convId", GC)
	group.GET("/gpts/mine", Mine)
	group.GET("/gpts", Gpts)

	// 状态相关
	group.GET("/status", Status)
	group.GET("/endpoint", EndPoint)
	group.POST("/carpage", CarPage)

	// 列表
	// group.GET("/list", List)

	// 封禁
	group.GET("/workspace/deactivated", workspace.Deactivated)

}
