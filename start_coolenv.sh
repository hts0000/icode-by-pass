#!/bin/bash

# start local go server
go run ./go-test/server/main.go

# --add-host: 改为自己的docker的网关ip
docker run --name coolenv -p 18000:18000 --add-host apis.imooc.com:172.17.0.1 -e ICODE="xxxxxx" -e GODEBUG=x509ignoreCN=0 coolenv-bypass-icode
