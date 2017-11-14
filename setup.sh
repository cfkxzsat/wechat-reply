#! /bin/bash

go install wechat.go
project_name=$GOBIN/wechat-reply
cp -r data/ $project_name
