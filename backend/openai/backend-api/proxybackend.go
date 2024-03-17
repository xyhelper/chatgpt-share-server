package backendapi

import (
	"backend/config"
	"backend/modules/chatgpt/model"
	"backend/utility"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

func ProxyBackend(r *ghttp.Request) {
	ctx := r.GetCtx()
	// usertoken := r.Session.MustGet("usertoken").String()
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

	Authorization := r.Header.Get("Authorization")
	if Authorization != "" {
		r.Header.Set("Authorization", "Bearer "+carinfo.AccessToken)
	}

	u, _ := url.Parse(config.CHATPROXY)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		g.Log().Error(ctx, e)
		writer.WriteHeader(http.StatusBadGateway)
	}
	newreq := r.Request.Clone(ctx)
	newreq.URL.Host = u.Host
	newreq.URL.Scheme = u.Scheme
	newreq.Host = u.Host
	// g.Dump(newreq.Header)
	newreq.Header.Set("authkey", config.AUTHKEY)
	proxy.ServeHTTP(r.Response.Writer.RawWriter(), newreq)
}
func ProxyBackendWithCar(r *ghttp.Request) {
	ctx := r.GetCtx()
	// usertoken := r.Session.MustGet("usertoken").String()
	carid := r.Session.MustGet("carid").String()
	conv := r.GetRouter("convid").String()

	refer := r.Referer()
	g.Log().Info(ctx, "refer:", refer)

	referArr := gstr.Split(refer, "/c/")
	if len(referArr) > 1 {
		conv = referArr[1]
		g.Log().Info(ctx, "conv:", conv, "from referer", refer)
	}

	// chatgptaccountid := r.Header.Get("ChatGPT-Account-ID")
	if conv != "" {
		g.Log().Info(ctx, "conv:", conv)
		// 查询会话
		result, err := cool.DBM(model.NewChatgptConversations()).Where(g.Map{
			"convid": conv,
		}).Cache(gdb.CacheOption{
			Duration: time.Hour,
			Name:     conv,
			Force:    true,
		}).One()
		if err != nil {
			g.Log().Error(ctx, err)
			r.Response.Status = 500
			r.Response.WriteJson(g.Map{
				"detail": "Internal Server Error",
			})
			return
		}
		if result == nil {
			r.Response.Status = 404
			r.Response.WriteJson(g.Map{
				"detail": "Can't load conversation " + conv,
			})
			return
		}
		email := result["email"].String()
		caridVar, err := g.Redis("cool").Get(ctx, "email:"+email)
		if err != nil {
			g.Log().Error(ctx, err)
			r.Response.Status = 500
			r.Response.WriteJson(g.Map{
				"detail": "Internal Server Error",
			})
			return
		}
		carid = caridVar.String()
		g.Log().Info(ctx, conv, "email:", email, "carid:", caridVar.String())

		if carid == "" {
			r.Response.Status = 404
			r.Response.WriteJson(g.Map{
				"detail": "The car " + conv + " belongs to is unavailable",
			})
			return
		}
		// r.Session.Set("carid", carid)
		chatgptaccountid := result["chatgptaccountid"].String()
		if chatgptaccountid != "" {
			r.Header.Set("ChatGPT-Account-ID", chatgptaccountid)
		} else {
			r.Header.Del("ChatGPT-Account-ID")
		}
		// r.Session.Set("carid", carid)
		// r.Session.Set("chatgptaccountid", chatgptaccountid)
		r.Session.Set("convcarid", carid)
		r.Session.Set("convchatgptaccountid", chatgptaccountid)
	}

	carinfo, err := utility.CheckCar(ctx, carid)
	if err != nil {
		g.Log().Error(ctx, err)
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"detail": "Authentication credentials were not provided.",
		})
		return
	}

	Authorization := r.Header.Get("Authorization")
	if Authorization != "" {
		r.Header.Set("Authorization", "Bearer "+carinfo.AccessToken)
	}
	r.Header.Del("Accept-Encoding")
	u, _ := url.Parse(config.CHATPROXY)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		g.Log().Error(ctx, e)
		writer.WriteHeader(http.StatusBadGateway)
	}
	g.Log().Debug(ctx, "ModifyResponse", r.Request.URL.Path)

	proxy.ModifyResponse = func(resp *http.Response) error {
		g.Log().Debug(ctx, "ModifyResponse", r.Request.URL.Path)
		// 如果 路径为 /backend-api/files/file-7FQ4EGrvYXGuQrNpmh8usd1p/download 格式
		if gstr.Contains(r.Request.URL.Path, "/backend-api/files/") && gstr.Contains(r.Request.URL.Path, "/download") {
			// 获取返回内容
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			// 修改返回内容
			body = []byte(strings.ReplaceAll(string(body), "https://files.oaiusercontent.com", config.FILESERVER))
			// 写入返回内容
			resp.Body = io.NopCloser(strings.NewReader(string(body)))

		}
		// 如果 路径 为/backend-api/files/file-TraKt2p9Q6y4JL5W1GBjqCaU/uploaded 格式
		if gstr.Contains(r.Request.URL.Path, "/backend-api/files/") && gstr.Contains(r.Request.URL.Path, "/uploaded") {
			// 获取返回内容
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			// 修改返回内容
			body = []byte(strings.ReplaceAll(string(body), "https://files.oaiusercontent.com", config.FILESERVER))
			// 写入返回内容
			resp.Body = io.NopCloser(strings.NewReader(string(body)))
		}
		// 如果路径为 /backend-api/files
		if gstr.Contains(r.Request.URL.Path, "/backend-api/files") {
			// 获取返回内容
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			// 修改返回内容
			body = []byte(strings.ReplaceAll(string(body), "https://files.oaiusercontent.com", config.FILESERVER))
			// 写入返回内容
			resp.Body = io.NopCloser(strings.NewReader(string(body)))
		}

		return nil
	}
	newreq := r.Request.Clone(ctx)
	newreq.URL.Host = u.Host
	newreq.URL.Scheme = u.Scheme
	newreq.Host = u.Host
	// g.Dump(newreq.Header)
	newreq.Header.Set("authkey", config.AUTHKEY)
	proxy.ServeHTTP(r.Response.Writer.RawWriter(), newreq)
}
