package gui

import (
	"io/ioutil"
	"os"
)

func loadIcon(path string) ([]byte ,error){
	f ,e := os.OpenFile(path ,os.O_RDONLY ,0755)
	if e != nil{
		return nil ,e
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}