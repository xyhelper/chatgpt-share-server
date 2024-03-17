package backendapi

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	s := g.Server()
	backendApiGroup := s.Group("/backend-api")
	backendApiGroup.ALL("/*", ProxyBackend)
	backendApiGroup.GET("/me", Me)                                       // 获取当前用户信息
	backendApiGroup.GET("/conversations", Conversations)                 // 获取会话列表
	backendApiGroup.POST("/conversation", Conversation)                  // 会话
	backendApiGroup.POST("/conversation/gen_title/:id", GenTitle)        // 生成标题
	backendApiGroup.GET("/conversation/:convid", ProxyBackendWithCar)    // 获取会话
	backendApiGroup.GET("/files/:fileid/download", ProxyBackendWithCar)  // 文件下载
	backendApiGroup.POST("/files/:fileid/uploaded", ProxyBackendWithCar) // 文件上传
	backendApiGroup.POST("/files", ProxyBackendWithCar)                  // 文件上传
	backendApiGroup.PATCH("/conversation/:id", ConversationPATCH)        // 修改会话 删除会话
	backendApiGroup.PATCH("/conversations", Conversations)               // 清空会话列表

	// 禁止访问的路径
	backendApiGroup.ALL("/payments/customer_portal", Error404)   // 支付
	backendApiGroup.ALL("/settings/beta_features", BetaFeatures) // 设置
	backendApiGroup.ALL("/shared_conversations", Error404)       // 共享会话
	backendApiGroup.ALL("/accounts/data_export", Error404)       // 导出数据
	backendApiGroup.ALL("/user_system_messages", Error404)       // 系统消息
	backendApiGroup.ALL("/payments/checkout", Error404)          // 支付
	backendApiGroup.ALL("/accounts/*/invites", Error404)         // 邀请
	backendApiGroup.ALL("/accounts/*/users/*", Error404)         // 成员
}

func Error404(r *ghttp.Request) {
	r.Response.WriteStatus(404)
}
