:: u 更新本地已存在的包以及依赖的包, 否则只通过网络下载本地不存在的包
:: v 信息输出
@echo off

go get -v github.com/nsf/gocode
go get -v github.com/rogpeppe/godef
go get -v github.com/zmb3/gogetdoc
go get -v github.com/golang/lint/golint
go get -v github.com/ramya-rao-a/go-outline
go get -v sourcegraph.com/sqs/goreturns
go get -v golang.org/x/tools/cmd/gorename
go get -v github.com/tpng/gopkgs
go get -v github.com/acroca/go-symbols
go get -v golang.org/x/tools/cmd/guru
go get -v github.com/cweill/gotests/...
go get -v golang.org/x/tools/cmd/godoc
go get -v github.com/fatih/gomodifytags
go get -v github.com/josharian/impl
go get -v github.com/derekparker/delve/cmd/dlv

go get -v github.com/alexmullins/zip
go get -v github.com/davecgh/go-spew

go get -v google.golang.org/grpc

go get -v github.com/lexical005/mysql
go get -v github.com/lexical005/protobuf
go get -v github.com/lexical005/toml
go get -v github.com/lexical005/xlsx

pause
