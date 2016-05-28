package main

import (
	"image/png"
	"net/http"

	"github.com/afocus/captcha"
)

var cap *captcha.Captcha

func main() {

	cap = captcha.New()

	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	cap.SetOpt(5, captcha.ALL, captcha.NORMAL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		img, str := cap.Create(128, 64)
		png.Encode(w, img)
		println(str)
	})

	http.ListenAndServe(":8085", nil)
}
