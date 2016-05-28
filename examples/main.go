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

	// 不调用SetOpt时 默认为(4,0,4,Color{0,0,0})
	cap.SetOpt(5, 3, 6)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		img, str := cap.Create(128, 64)
		png.Encode(w, img)
		println(str)
	})

	http.ListenAndServe(":8085", nil)
}
