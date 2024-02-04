#!/usr/bin/env sh

# 确保脚本抛出遇到的错误
set -e

# 生成静态文件
yarn docs:build


# 进入生成的文件夹
cd docs/.vuepress/dist

git init -b gh-pages
git add -A
git commit -m "deploy | $(date +'%Y-%m-%d %H:%M:%S')"


git push -f git@github.com:xyhelper/chatgpt-share-server gh-pages:gh-pages
