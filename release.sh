#!/bin/bash
set -e
# 更新 car-list 仓库    
rm -rf ./backend/resource/public/list
git clone -b dist https://github.com/Hanwencc/car-list.git ./backend/resource/public/list
# 检测是否存在目录 ./backend/resource/public/xyhelper
if [ ! -d "./backend/resource/public/xyhelper" ]; then
    echo "Create directory ./backend/resource/public/xyhelper"
    mkdir -p "./backend/resource/public/xyhelper"

    cd frontend
    yarn build
    cd ..
fi

cd backend
gf build main.go -a amd64 -s linux -p ./temp
gf docker main.go -p -t xyhelper/chatgpt-share-server:latest
now=$(date +"%Y%m%d%H%M%S")
# 以当前时间为版本号
docker tag xyhelper/chatgpt-share-server:latest xyhelper/chatgpt-share-server:$now
docker push xyhelper/chatgpt-share-server:$now
echo "release success" $now
# 写入发布日志 release.log
echo $now >> ../release.log
