package main

import "testing"

func TestGetIP(t *testing.T){
	t.Log(getIP("apnic|CN|ipv4|39.0.8.0|2048|20110412|allocated"))
	t.Log(getIP("apnic|CN|ipv4|40.125.128.0|32768|20150223|allocated"))
}
