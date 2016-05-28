package captcha

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/afocus/captcha/graphics"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Image 图片
type Image struct {
	*image.RGBA
}

// NewImage 创建一个新的图片
func NewImage(w, h int) *Image {
	img := &Image{image.NewRGBA(image.Rect(0, 0, w, h))}
	return img
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	return -1
}

// DrawLine 画直线
// Bresenham算法(https://zh.wikipedia.org/zh-cn/布雷森漢姆直線演算法)
// x1,y1 起点 x2,y2终点
func (img *Image) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	dx, dy, flag := int(math.Abs(float64(x2-x1))),
		int(math.Abs(float64(y2-y1))),
		false
	if dy > dx {
		flag = true
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
	}
	ix, iy := sign(x2-x1), sign(y2-y1)
	n2dy := dy * 2
	n2dydx := (dy - dx) * 2
	d := n2dy - dx
	for x1 != x2 {
		if d < 0 {
			d += n2dy
		} else {
			y1 += iy
			d += n2dydx
		}
		if flag {
			img.Set(y1, x1, c)
		} else {
			img.Set(x1, y1, c)
		}
		x1 += ix
	}
}

func (img *Image) drawCircle8(xc, yc, x, y int, c color.Color) {
	img.Set(xc+x, yc+y, c)
	img.Set(xc-x, yc+y, c)
	img.Set(xc+x, yc-y, c)
	img.Set(xc-x, yc-y, c)
	img.Set(xc+y, yc+x, c)
	img.Set(xc-y, yc+x, c)
	img.Set(xc+y, yc-x, c)
	img.Set(xc-y, yc-x, c)
}

// DrawCircle 画圆
// xc,yc 圆心坐标 r 半径 fill是否填充颜色
func (img *Image) DrawCircle(xc, yc, r int, fill bool, c color.Color) {
	size := img.Bounds().Size()
	// 如果圆在图片可见区域外，直接退出
	if xc+r < 0 || xc-r >= size.X || yc+r < 0 || yc-r >= size.Y {
		return
	}
	x, y, d := 0, r, 3-2*r
	for x <= y {
		if fill {
			for yi := x; yi <= y; yi++ {
				img.drawCircle8(xc, yc, x, yi, c)
			}
		} else {
			img.drawCircle8(xc, yc, x, y, c)
		}
		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// DrawString 写字
func (img *Image) DrawString(font *truetype.Font, c color.Color, str string, fontsize float64, x, y int) {
	ctx := freetype.NewContext()
	// default 72dpi
	ctx.SetDst(img)
	ctx.SetClip(img.Bounds())
	ctx.SetSrc(image.NewUniform(c))
	ctx.SetFontSize(fontsize)
	ctx.SetFont(font)
	// 写入文字的位置
	pt := freetype.Pt(x, y+ctx.PointToFixed(fontsize).Ceil())
	ctx.DrawString(str, pt)
}

// Rotate 旋转
func (img *Image) Rotate(angle float64) image.Image {
	nb := image.NewRGBA(img.Bounds())
	graphics.Rotate(nb, img, &graphics.RotateOptions{Angle: angle})
	return nb
}

// 填充背景色
func (img *Image) FillBkg(c color.Color) {
	draw.Draw(img, img.Bounds(), image.NewUniform(c), image.ZP, draw.Src)
}
