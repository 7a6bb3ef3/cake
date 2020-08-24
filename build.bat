@echo off

set GOPROXY=https://goproxy.io

go build -o ./cakeserv.exe github.com/nynicg/cake/serv

go build -o ./cakecli.exe  github.com/nynicg/cake/client