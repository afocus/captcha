package draw

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func sign(x int) int {
	if x > 0 {
		return 1
	}
	return -1
}

//-----------------------------------------------------------------------------------------------------------//

// Image 图片
type Image struct {
	*image.RGBA
}

// NewImage 创建一个新的图片
func NewImage(w, h int) *Image {
	img := &Image{image.NewRGBA(image.Rect(0, 0, w, h))}
	return img
}

// 生成 PNG 图片
func (self *Image) EncodedPNG() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, self); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 写入到 io
func (self *Image) WriteTo(w io.Writer) error {
	return png.Encode(w, self)
}

// DrawLine 画直线
// Bresenham 算法 (https://zh.wikipedia.org/zh-cn/布雷森漢姆直線演算法)
// x1,y1 起点 x2,y2终点
func (self *Image) DrawLine(x1, y1, x2, y2 int, c color.Color) {
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
			self.Set(y1, x1, c)
		} else {
			self.Set(x1, y1, c)
		}
		x1 += ix
	}
}

func (self *Image) drawCircle8(xc, yc, x, y int, c color.Color) {
	self.Set(xc+x, yc+y, c)
	self.Set(xc-x, yc+y, c)
	self.Set(xc+x, yc-y, c)
	self.Set(xc-x, yc-y, c)
	self.Set(xc+y, yc+x, c)
	self.Set(xc-y, yc+x, c)
	self.Set(xc+y, yc-x, c)
	self.Set(xc-y, yc-x, c)
}

// DrawCircle 画圆
// xc,yc 圆心坐标 r 半径 fill是否填充颜色
func (self *Image) DrawCircle(xc, yc, r int, fill bool, c color.Color) {
	size := self.Bounds().Size()
	// 如果圆在图片可见区域外，直接退出
	if xc+r < 0 || xc-r >= size.X || yc+r < 0 || yc-r >= size.Y {
		return
	}
	x, y, d := 0, r, 3-2*r
	for x <= y {
		if fill {
			for yi := x; yi <= y; yi++ {
				self.drawCircle8(xc, yc, x, yi, c)
			}
		} else {
			self.drawCircle8(xc, yc, x, y, c)
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
func (self *Image) DrawString(font *truetype.Font, c color.Color, str string, fontsize float64) {
	ctx := freetype.NewContext()
	// default 72dpi
	ctx.SetDst(self)
	ctx.SetClip(self.Bounds())
	ctx.SetSrc(image.NewUniform(c))
	ctx.SetFontSize(fontsize)
	ctx.SetFont(font)
	// 写入文字的位置
	pt := freetype.Pt(0, int(-fontsize/6)+ctx.PointToFixed(fontsize).Ceil())
	ctx.DrawString(str, pt)
}

// Rotate 旋转
func (self *Image) Rotate(angle float64) image.Image {
	return new(rotate).Rotate(angle, self.RGBA).transformRGBA()
}

// 填充背景
func (self *Image) FillBackground(c image.Image) {
	draw.Draw(self, self.Bounds(), c, image.ZP, draw.Over)
}

// 水波纹, amplude=振幅, period=周期
// copy from https://github.com/dchest/captcha/blob/master/image.go
func (self *Image) distortTo(amplude float64, period float64) {
	w := self.Bounds().Max.X
	h := self.Bounds().Max.Y

	oldm := self.RGBA

	dx := 1.4 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			rgba := oldm.RGBAAt(x+int(xo), y+int(yo))
			if rgba.A > 0 {
				oldm.SetRGBA(x, y, rgba)
			}
		}
	}
}

func inBounds(b image.Rectangle, x, y float64) bool {
	if x < float64(b.Min.X) || x >= float64(b.Max.X) {
		return false
	}
	if y < float64(b.Min.Y) || y >= float64(b.Max.Y) {
		return false
	}
	return true
}

//-----------------------------------------------------------------------------------------------------------//

type rotate struct {
	dx   float64
	dy   float64
	sin  float64
	cos  float64
	neww float64
	newh float64
	src  *image.RGBA
}

func radian(angle float64) float64 {
	return angle * math.Pi / 180.0
}

func (self *rotate) Rotate(angle float64, src *image.RGBA) *rotate {
	self.src = src
	srsize := src.Bounds().Size()
	width, height := srsize.X, srsize.Y

	// 源图四个角的坐标（以图像中心为坐标系原点）
	// 左下角,右下角,左上角,右上角
	srcwp, srchp := float64(width)*0.5, float64(height)*0.5
	srcx1, srcy1 := -srcwp, srchp
	srcx2, srcy2 := srcwp, srchp
	srcx3, srcy3 := -srcwp, -srchp
	srcx4, srcy4 := srcwp, -srchp

	self.sin, self.cos = math.Sincos(radian(angle))
	// 旋转后的四角坐标
	desx1, desy1 := self.cos*srcx1+self.sin*srcy1, -self.sin*srcx1+self.cos*srcy1
	desx2, desy2 := self.cos*srcx2+self.sin*srcy2, -self.sin*srcx2+self.cos*srcy2
	desx3, desy3 := self.cos*srcx3+self.sin*srcy3, -self.sin*srcx3+self.cos*srcy3
	desx4, desy4 := self.cos*srcx4+self.sin*srcy4, -self.sin*srcx4+self.cos*srcy4

	// 新的高度很宽度
	self.neww = math.Max(math.Abs(desx4-desx1), math.Abs(desx3-desx2)) + 0.5
	self.newh = math.Max(math.Abs(desy4-desy1), math.Abs(desy3-desy2)) + 0.5
	self.dx = -0.5*self.neww*self.cos - 0.5*self.newh*self.sin + srcwp
	self.dy = 0.5*self.neww*self.sin - 0.5*self.newh*self.cos + srchp
	return self
}

func (self *rotate) pt(x, y int) (float64, float64) {
	return float64(-y)*self.sin + float64(x)*self.cos + self.dy,
		float64(y)*self.cos + float64(x)*self.sin + self.dx
}

func (self *rotate) transformRGBA() image.Image {
	srcb := self.src.Bounds()
	b := image.Rect(0, 0, int(self.neww), int(self.newh))
	dst := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sx, sy := self.pt(x, y)
			if inBounds(srcb, sx, sy) {
				// 消除锯齿填色
				c := bili.RGBA(self.src, sx, sy)
				off := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
				dst.Pix[off+0] = c.R
				dst.Pix[off+1] = c.G
				dst.Pix[off+2] = c.B
				dst.Pix[off+3] = c.A
			}
		}
	}
	return dst
}

//-----------------------------------------------------------------------------------------------------------//
