package fileserver

import "github.com/gogf/gf/v2/frame/g"

func init() {
	s := g.Server()
	s.BindHandler("/file-{fileid}", Proxyfile)
}
