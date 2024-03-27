package auth

import (
    "backend/config" // 导入配置模块
    "backend/utility" // 导入工具包

    "github.com/cool-team-official/cool-admin-go/cool" // 导入第三方库
    "github.com/gogf/gf/v2/encoding/gjson" // 导入GoFrame的JSON处理库
    "github.com/gogf/gf/v2/frame/g" // 导入GoFrame框架
    "github.com/gogf/gf/v2/net/ghttp" // 导入GoFrame的HTTP模块
    "github.com/gogf/gf/v2/util/gconv" // 导入GoFrame的类型转换库
    "github.com/narqo/go-badge" // 导入徽章生成库
)

// Login函数处理登录请求
func Login(r *ghttp.Request) {
    ctx := r.GetCtx() // 获取请求上下文
    if r.Method == "GET" { // 处理GET请求，显示登录页面
        displayLoginPage(r, "")
        return
    } else if r.Method == "POST" { // 处理POST请求，执行登录逻辑
        handlePostLogin(ctx, r)
    }
}

// displayLoginPage函数用于显示登录页面
func displayLoginPage(r *ghttp.Request, errorMsg string) {
    badgeSVG, _ := badge.RenderBytes("登录", "点击登录", "blue") // 生成徽章
    r.Response.WriteTpl("login.html", g.Map{ // 渲染登录模板
        "badge":   string(badgeSVG), // 将徽章SVG字符串传递给模板
        "BuyLink": "https://chat.bjp666.link", // 购买按钮链接
        "error":   errorMsg, // 错误信息
    })
}

// handlePostLogin函数处理登录表单提交
func handlePostLogin(ctx g.Ctx, r *ghttp.Request) {
    req := r.GetMapStrStr() // 获取POST请求的表单数据
    loginVar := g.Client().PostVar(ctx, config.OauthUrl, req) // 发送请求到OAuth服务
    loginJson := gjson.New(loginVar) // 将返回的数据解析为JSON
    code := loginJson.Get("code").Int() // 获取登录结果的状态码

    if code != 1 { // 如果登录失败
        msg := loginJson.Get("msg").String() // 获取失败信息
        displayLoginPage(r, msg) // 显示登录页面，并展示错误信息
        return
    }

    // 登录成功，选择剩余次数最多的账号
    selectedCarid, _ := findAccountWithMostCalls()
    // 设置用户会话信息
    r.Session.Set("usertoken", req["usertoken"])
    r.Session.Set("carid", selectedCarid)
    r.Session.SetMaxAge(432000) // 会话有效期设为5天
    r.Response.RedirectTo("/") // 重定向到首页
}

// findAccountWithMostCalls函数查找剩余次数最多的账号
func findAccountWithMostCalls() (string, int) {
    maxCount := -1
    selectedCarid := ""
    carids := []string{"carid1", "carid2", "carid3"} // 假定的账号ID列表

    for _, carid := range carids {
        count := utility.GetStatsInstance(carid).GetCallCount() // 获取每个账号的剩余次数
        if count > maxCount {
            maxCount = count
            selectedCarid = carid
        }
    }
    return selectedCarid, maxCount // 返回剩余次数最多的账号ID和次数
}

// LoginToken函数处理基于令牌的登录
func LoginToken(r *ghttp.Request) {
    ctx := r.GetCtx()
    req := r.GetMapStrStr() // 获取请求参数
    resptype := req["resptype"] // 获取响应类型

    loginVar := g.Client().PostVar(ctx, config.OauthUrl, req) // 向OAuth服务发送请求
    loginJson := gjson.New(loginVar) // 解析返回的JSON数据
    code := loginJson.Get("code").Int() // 获取状态码

    if code != 1 { // 如果验证失败
        msg := loginJson.Get("msg").String() // 获取错误信息
        if resptype == "json" { // 如果客户端期望JSON响应
            r.Response.WriteJson(g.Map{
                "code": 0,
                "msg":  msg,
            })
            return
        } else { // 如果客户端期望页面响应
            displayLoginPage(r, msg) // 显示登录页面，并展示错误信息
            return
        }
    }

    // 验证成功，设置会话信息
    r.Session.Set("usertoken", req["usertoken"])
    r.Session.Set("carid", req["carid"])
    if resptype == "json" { // 客户端期望JSON响应
        r.Response.WriteJson(g.Map{
            "code": 1,
            "msg":  "登录成功",
        })
        return
    } else { // 重定向到首页
        r.Response.RedirectTo("/")
    }
}
