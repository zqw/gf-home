#!/usr/bin/env bash

set -x;

##################
#cd to shell dir
##################
if [[ $0 =~ ^\/.* ]]    #判断当前脚本是否为绝对路径，匹配以/开头下的所有
then
  script=$0
else
  script=$(pwd)/$0
fi
script=`readlink -f $script`   #获取文件的真实路径
script_path=${script%/*}     #获取文件所在的目录
realpath=$(readlink -f $script_path)   #获取文件所在目录的真实路径
#echo $script
#echo $script_path
#echo $realpath

cd $realpath
cd ..
##################
#mod to vendor
##################
rm -rf vendor
go mod vendor



##################
#编译
##################
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go



##################
#运行
##################

./main