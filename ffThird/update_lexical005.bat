:: u 更新本地已存在的包以及依赖的包, 否则只通过网络下载本地不存在的包
:: v 信息输出
@echo off

go get -u -v github.com/lexical005/mysql
go get -u -v github.com/lexical005/toml
go get -u -v github.com/lexical005/xlsx

pause
