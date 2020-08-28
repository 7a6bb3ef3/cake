// +build linux darwin

package main

// although https://github.com/getlantern/systray is a cross platfrom library ,
// in order to avoid importing more dependency(gtk3...) , leave it empty
// TODO
func RunAsIcon(onexit func()){

}