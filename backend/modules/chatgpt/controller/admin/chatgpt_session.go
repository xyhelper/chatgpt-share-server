package admin

import (
	"backend/modules/chatgpt/service"
	"context"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type ChatgptSessionController struct {
	*cool.Controller
}

func init() {
	var chatgpt_session_controller = &ChatgptSessionController{
		&cool.Controller{
			Prefix:  "/admin/chatgpt/session",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewChatgptSessionService(),
		},
	}
	// 注册路由
	cool.RegisterController(chatgpt_session_controller)
}

// 增加 InfoBYCarID
type ChatgptSessionInfoBYCarIDReq struct {
	g.Meta `path:"/infobycarid" method:"GET"` // 这里指定路由和请求方法
}
type ChatgptSessionInfoBYCarIDRes struct { // 返回值
	*cool.BaseRes
	Data interface{} `json:"data"`
}

func (c *ChatgptSessionController) Welcome(ctx context.Context, req *ChatgptSessionInfoBYCarIDReq) (res *ChatgptSessionInfoBYCarIDRes, err error) {
	res = &ChatgptSessionInfoBYCarIDRes{
		BaseRes: cool.Ok("Welcome to Cool Admin Go"),
		Data:    gjson.New(`{"name": "Cool Admin Go", "age":0}`),
	}
	return
}
