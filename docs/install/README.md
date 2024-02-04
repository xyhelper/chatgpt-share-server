---
sidebar: auto
---
# 安装指南

## 前置条件

- 1 台 2C2G 以上服务器,推荐使用 Ubuntu 20.04 LTS 以上版本,服务器不强制要求可以访问官网,但需要能够访问`demo.xyhelper.cn`
- 安装 Docker 及 docker compose,推荐使用[官方安装脚本](https://docs.docker.com/engine/install/ubuntu/#install-using-the-convenience-script)
- 安装 git,推荐使用`apt install git`
- 安装`caddy`或`nginx`等反向代理服务器,推荐使用`caddy`,推荐使用[官方安装脚本](https://caddyserver.com/docs/install#debian-ubuntu-raspbian)
- 使用 root 用户安装,不建议使用非 root 用户通过`sudo`安装，可能会导致权限问题。

## 常用系统优化

### 系统时间同步

```bash
apt install ntpdate -y
ntpdate time.windows.com
```

### 系统时区设置

```bash
timedatectl set-timezone Asia/Shanghai
```

### 设置最大文件打开数

```bash
echo "fs.inotify.max_user_instances=5120" >> /etc/sysctl.conf
echo "fs.inotify.max_user_watches=2621440" >> /etc/sysctl.conf
echo "fs.file-max=65535" >> /etc/sysctl.conf

sysctl -p
```

### 限制 docker 日志大小

```bash
mkdir -p /etc/docker

cat > /etc/docker/daemon.json <<EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

systemctl restart docker
```

### 设置防火墙

```bash
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 22/tcp

ufw enable
```

## 安装

### 安装 ChatGPT-Share-Server

```bash
curl -sSfL https://raw.githubusercontent.com/xyhelper/chatgpt-share-server/deploy/quick-install.sh | bash
```

### 安装 ChatGPT-Share-Server(手动)

```bash
git clone -b deploy --depth=1 https://github.com/xyhelper/chatgpt-share-server.git chatgpt-share
cd chatgpt-share
./deploy.sh
```

## 反代配置

### Caddy

假定你的域名为`xxx.yourdomain.com`,请将`xxx.yourdomain.com`解析到服务器 IP,然后执行以下命令,注意将`xxx.yourdomain.com`替换为你的域名。

```bash
cat > /etc/caddy/Caddyfile <<EOF
xxx.yourdomain.com {
    reverse_proxy 127.0.0.1:8300
}
EOF

systemctl reload caddy
```

## 管理

后台管理地址: `http://xxx.yourdomain.com/bpjadmin`

默认账号: `admin`

默认密码: `123456`

## 常用管理命令

### 重启

```bash
cd chatgpt-share
docker compose restart
```

### 停止

```bash
cd chatgpt-share
docker compose stop
```

### 启动

```bash
cd chatgpt-share
./deploy.sh
```

### 升级

```bash
cd chatgpt-share
./deploy.sh
```

### 查看日志

```bash 
cd chatgpt-share
docker compose logs -f
```
