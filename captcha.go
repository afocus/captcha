package captcha

import (
	"code.google.com/p/graphics-go/graphics"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Captcha struct {
	font *truetype.Font
	pool *sync.Pool
}

func New(path string) (*Captcha, error) {
	fontdata, erro := ioutil.ReadFile(path)
	if erro != nil {
		return nil, erro
	}
	font, erro := freetype.ParseFont(fontdata)
	if erro != nil {
		return nil, erro
	}
	c := &Captcha{
		font: font,
		pool: &sync.Pool{New: func() interface{} {
			return freetype.NewContext()
		}},
	}
	return c, nil
}

func (c *Captcha) buildChar(size int, char rune) draw.Image {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	ctx := c.pool.Get().(*freetype.Context)
	defer c.pool.Put(ctx)
	ctx.SetDst(dst)
	ctx.SetClip(dst.Bounds())
	ctx.SetSrc(image.Black)
	ctx.SetFontSize(float64(size))
	ctx.SetFont(c.font)
	ctx.DrawString(string(char), freetype.Pt(0, size))
	return dst
}

func (c *Captcha) Draw(write http.ResponseWriter, w, h int, str string) error {
	if h < 20 {
		h = 20
	}
	if w < 60 {
		w = 60
	}

	padding := h / 4
	fontsize := h - padding*2
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	//draw.Draw(dst, dst.Bounds(), image.White, image.ZP, draw.Src)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	offset := (w - padding*2) / len(str)
	psize := fontsize + padding
	for i, char := range str {
		b := c.buildChar(fontsize, char)
		nb := image.NewRGBA(image.Rect(0, 0, psize, psize))
		if erro := graphics.Rotate(nb, b, &graphics.RotateOptions{Angle: r.Float64()}); erro != nil {
			println(erro.Error())
			continue
		}

		draw.Draw(dst, image.Rect(i*offset, 0, (i+1)*offset, psize), nb, image.ZP, draw.Src)
	}

	return png.Encode(write, dst)
}
