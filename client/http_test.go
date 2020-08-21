package main

import (
	"net/http"
	"testing"
)

func TestHttp(t *testing.T){
	mux := http.NewServeMux()
	mux.HandleFunc("/" , func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("114514"))
	})
	if e := http.ListenAndServe("127.0.0.1:1918" ,mux);e != nil{
		panic(e)
	}
}
