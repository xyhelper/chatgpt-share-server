package fileserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Proxyfile(r *ghttp.Request) {
	ctx := r.GetCtx()
	g.Log().Info(ctx, r.Method, r.URL.Path)
	u, _ := url.Parse("https://files.oaiusercontent.com")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		g.Log().Error(ctx, e)
		writer.WriteHeader(http.StatusBadGateway)
	}
	newreq := r.Request.Clone(ctx)
	newreq.URL.Host = u.Host
	newreq.URL.Scheme = u.Scheme
	newreq.Host = u.Host
	newreq.Header.Set("Referer", "https://chat.openai.com")
	// g.Dump(newreq.Header)
	proxy.ServeHTTP(r.Response.Writer.RawWriter(), newreq)
}
