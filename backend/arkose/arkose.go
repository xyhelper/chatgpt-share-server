package arkose

import "github.com/gogf/gf/v2/frame/g"

func init() {
	s := g.Server()
	arkoseGroup := s.Group("/")
	arkoseGroup.ALL("/v2/*", ProxyArkose)
	arkoseGroup.ALL("/fc/*", ProxyArkose)
	arkoseGroup.ALL("/cdn/fc/*any", ProxyArkose)
	arkoseGroup.ALL("/rtig/image", ProxyArkose)
	arkoseGroup.ALL("/params/sri/dapib/*", ProxyArkose)
	arkoseGroup.ALL("/dapib/*any", ProxyArkose)
}
