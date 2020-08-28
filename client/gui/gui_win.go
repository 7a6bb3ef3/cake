// +build windows

package gui

import (
	"github.com/getlantern/systray"
	"github.com/nynicg/cake/lib/log"
	"github.com/skratchdot/open-golang/open"
)

func RunAsIcon(onexit func()){
	go systray.Run(onReady , onexit)
}

func onReady() {
	icob ,e := loadIcon("cake.ico")
	if e != nil{
		log.Error("load icon " ,e)
		panic(e)
	}
	systray.SetTemplateIcon(icob ,icob)
	systray.SetTitle("Cake")
	systray.SetTooltip("Love and Spanner")

	go func() {
		stt := systray.AddMenuItem("Status: OFF", "")
		stt.Disable()
		systray.AddSeparator()
		update := systray.AddMenuItem("Update", "")
		runStop := systray.AddMenuItem("Run", "")
		systray.AddSeparator()
		mQuitOrig := systray.AddMenuItem("Quit", "")
		for {
			select {
			case <-update.ClickedCh:
				open.Run("https://github.com/nynicg/cake")
			case <-runStop.ClickedCh:
				open.Run("https://github.com/nynicg/cake")
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
