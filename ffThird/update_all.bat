:: u 更新本地已存在的包以及依赖的包, 否则只通过网络下载本地不存在的包
:: v 信息输出
@echo off

go get -u -v github.com/nsf/gocode
go get -u -v github.com/rogpeppe/godef
go get -u -v github.com/zmb3/gogetdoc
go get -u -v github.com/golang/lint/golint
go get -v -v github.com/golang/protobuf
go get -u -v github.com/ramya-rao-a/go-outline
go get -u -v sourcegraph.com/sqs/goreturns
go get -u -v golang.org/x/tools/cmd/gorename
go get -u -v github.com/tpng/gopkgs
go get -u -v github.com/acroca/go-symbols
go get -u -v golang.org/x/tools/cmd/guru
go get -u -v github.com/cweill/gotests/...
go get -u -v golang.org/x/tools/cmd/godoc
go get -u -v github.com/fatih/gomodifytags
go get -u -v github.com/josharian/impl
go get -u -v github.com/derekparker/delve/cmd/dlv

go get -u -v github.com/alexmullins/zip
go get -u -v github.com/davecgh/go-spew

go get -u -v github.com/lexical005/mysql
go get -u -v github.com/lexical005/toml
go get -u -v github.com/lexical005/xlsx

pause
