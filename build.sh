#!/bin/bash


#sudo apt-get install libgtk-3-dev libappindicator3-dev -y
#yum install gtk3-devel libappindicator-gtk3 -y

export GOPROXY=https://goproxy.io

go build -o ./cake github.com/nynicg/cake/serv

go build -o ./cakecli  github.com/nynicg/cake/client