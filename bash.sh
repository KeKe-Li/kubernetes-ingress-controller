#!/usr/bin/env bash
# 打开合并页面的脚步
other2master() {
open "https://github.com/KeKe-Li/kubernetes-ingress-controller/merge_requests"
}

branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" = "master" ]
then
master
else
other2master $branch
fi