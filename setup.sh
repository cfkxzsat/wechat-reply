#! /bin/bash

go install wechat.go
project_name=$GOBIN/wechat-reply
mkdir $project_name
cp -r data/ $project_name
cp $GOBIN/wechat $project_name
