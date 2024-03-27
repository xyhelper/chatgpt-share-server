#!/bin/bash
# This script is used to initialize the environment for the project
# Usage: ./init.sh
# Author: LiDong
# Date: 2022-11-08

# 出错退出
set -e

# 如果 REMOTE_CONTAINERS 为 true,则为容器开发环境,进行相关配置
if [ "$REMOTE_CONTAINERS" = "true" ]; then
    # 容器开发环境,配置容器内部的环境变量
    echo "Configuring environment variables for container development environment"
    # 记录hostname到 .devcontainer/data/hostname.txt 如果 .devcontainer/data 目录不存在,则创建
    if [ ! -d ".devcontainer/data" ]; then
        mkdir .devcontainer/data
        chmod 777 .devcontainer/data
    fi
    echo "$(hostname)" >.devcontainer/data/hostname.txt

    # # 配置goproxy
    # echo "Configuring goproxy"
    go env -w GO111MODULE=on
    go env -w GOPROXY=https://goproxy.cn,direct

    # # 配置npm mirror
    # echo "Configuring npm mirror"
    npm config set registry https://registry.npmmirror.com
    yarn config set registry https://registry.npmmirror.com

    # # 安装pgit
    # echo "Installing pgit"
    # curl -o pgit https://gitee.com/gcslaoli/pgit/raw/main/shell/pgit && chmod +x pgit && sudo mv pgit /usr/local/bin

    # 安装cool-tools
    echo "Installing cool-tools ..."
    go install github.com/cool-team-official/cool-admin-go/cool-tools@latest
    # 安装gf
    echo "Installing gf use mirror ..."
    wget -O gf \
        https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH) &&
        chmod +x gf &&
        ./gf install -y &&
        rm ./gf

fi
