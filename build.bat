@echo off

set GOPROXY=https://goproxy.io

go build -o ./cake.exe github.com/nynicg/cake/serv

echo if you want a non-GUI client, REMOVE the [-ldflags "-H=windowsgui"] in build params below and use [cakecli -nonGui] to run client
echo NOTE:if a client was compiled from -ldflags "-H=windowsgui" ,better not exec it with -nonGui flag
go build -o ./cakecli.exe -ldflags "-H=windowsgui" github.com/nynicg/cake/client

:: nonGui ver
:: go build -o ./cakecli.exe github.com/nynicg/cake/client
:: after conpiled ,exec "cakecli -nonGui [OPTIONS...]"