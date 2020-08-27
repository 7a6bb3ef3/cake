package log

import "testing"

func TestLog(t *testing.T) {
	Debug("1")
	Info("2")
	Warn("3")
	Error("4")
	Panic("5")
}
