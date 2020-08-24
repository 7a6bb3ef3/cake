#!/bin/bash

export GOPROXY=https://goproxy.io

go build -o ./cakeserv github.com/nynicg/cake/serv

go build -o ./cakecli  github.com/nynicg/cake/client