package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

type Captcha struct {
	colors      []color.Color
	strNum      int
	modal       int
	disturbance int
	font        *truetype.Font
}

// 颜色
type Color struct {
	R uint8
	G uint8
	B uint8
}

const (
	NUM   = 0 // 数字
	LOWER = 1 // 小写字母
	UPPER = 2 // 大写字母
	ALL   = 3 // 全部
)

const (
	NORMAL = 6
	MEDIUM = 10
	HIGH   = 16
)

func New() *Captcha {
	c := &Captcha{
		strNum:      4,
		disturbance: NORMAL,
		modal:       NUM,
	}
	c.colors = []color.Color{color.Black}
	return c
}

// SetFont 设置字体
func (c *Captcha) SetFont(path string) error {
	fontdata, erro := ioutil.ReadFile(path)
	if erro != nil {
		return erro
	}
	font, erro := freetype.ParseFont(fontdata)
	if erro != nil {
		return erro
	}
	c.font = font
	return nil
}

// SetOpt 设置配置信息
// params
// modal 0:纯数字 1:数字+小写字母 2:数字+大小写字母
// disturbance 干扰 数值越大越干扰
// strNum 验证码文字个数
func (c *Captcha) SetOpt(strNum int, modal int, disturbance int, colors ...Color) {
	c.strNum = strNum
	c.modal = modal
	if len(colors) > 0 {
		c.colors = c.colors[:0]
		for _, v := range colors {
			c.colors = append(c.colors, color.RGBA{v.R, v.G, v.B, 255})
		}
	}
	if disturbance < 4 {
		disturbance = 4
	}
	c.disturbance = disturbance

}

// 绘制背景
func (c *Captcha) drawBkg(img *Image) {

	// 填充主背景色
	//img.FillBkg(bgcolor)

	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	co := color.RGBA{0, 0, 0, 20}
	// 绘制干扰斑点
	for i := 0; i < c.disturbance; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		r := ra.Intn(size.Y / 4)
		img.DrawCircle(x, y, r, true, co)
	}
	colorindex := 0
	// 绘制干扰线
	for i := 0; i < c.disturbance; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(size.Y) * o
		h := ra.Intn(size.Y/10) * o
		if colorindex == len(c.colors) {
			colorindex = 0
		}
		img.DrawLine(x, y, x+w, y+h, c.colors[colorindex])
		colorindex++
	}

}

// 绘制文字
func (c *Captcha) drawString(img *Image, str string) {
	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	// 文字大小为图片高度的1/2
	fsize := size.Y / 2
	// 用于生成随机角度
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 文字之间的距离
	offset := size.X/len(str) - fsize/6
	// 文字在图形上的起点
	p := fsize / 2
	colorindex := 0
	// 逐个绘制文字到图片上
	for i, char := range str {
		// 创建单个文字图片
		// 以高为尺寸创建正方形的图形
		str := NewImage(size.Y, size.Y)
		if colorindex == len(c.colors) {
			colorindex = 0
		}
		str.DrawString(c.font, c.colors[colorindex], string(char), float64(fsize), p, p)
		colorindex++
		// 转换角度后的文字图形
		r := str.Rotate(r.Float64())
		// 计算文字位置
		clip := image.Rect(i*offset, 0, (i+1)*offset+fsize, offset+fsize)
		// 绘制到图片上
		draw.Draw(img, clip, r, image.ZP, draw.Over)
	}
}

// Create 生成一个验证码图片
func (c *Captcha) Create(w, h int) (*Image, string) {
	if h < 20 {
		h = 20
	}
	if w < 60 {
		w = 60
	}
	dst := NewImage(w, h)
	c.drawBkg(dst)
	str := string(c.randStr(c.strNum, c.modal))
	c.drawString(dst, str)
	return dst, str
}

var fontKinds = [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}

// 生成随机字符串
// size 个数 kind 模式
func (c *Captcha) randStr(size int, kind int) []byte {
	ikind, result := kind, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(3)
		}
		scope, base := fontKinds[ikind][0], fontKinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
