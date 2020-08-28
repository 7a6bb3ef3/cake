@echo off

set GOPROXY=https://goproxy.io

go build -o ./cake.exe github.com/nynicg/cake/serv

go build -o ./cakecli.exe -ldflags "-H=windowsgui" github.com/nynicg/cake/client