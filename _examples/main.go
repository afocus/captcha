package main

import (
	"fmt"
	"os"

	"github.com/qAison/captcha"
)

func main() {
	id := captcha.New()
	fileName := id + ".png"
	if oFile, err := os.Create(fileName); err != nil {
		fmt.Printf("create file error:%s\n", err.Error())
	} else {
		defer oFile.Close()

		if err = captcha.WriteImage(oFile, id); err != nil {
			fmt.Printf("image write error:%s\n", err.Error())
		} else {
			fmt.Printf("create image success:%s\n", fileName)
		}
	}
}
