import { defaultTheme, defineUserConfig } from "vuepress";

export default defineUserConfig({
  lang: "zh-CN",
  title: "ChatgptShareServer",
  description: "基于官网UI的共享账号方案",
  theme: defaultTheme({
    repo: "xyhelper/chatgpt-share-server",
    docsBranch: "master",
    docsDir: "docs",
    navbar: [
      { text: "首页", link: "/" },
      {
        text: "安装",
        link: "/install/",
      },
      {
        text: "配置",
        link: "/config/",
      },
    ],
    sidebar: {
      "/config/": [
        "/config/README.md",
        "/config/list.md",
        "/config/apiauth.md",
        "/config/freemode.md",
        "/config/oauth.md",
        "/config/resetpasswd.md",
        "/config/roma.md",
      ],
    },
  }),
});
