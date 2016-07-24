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
	frontColors []color.Color
	bkgColors   []color.Color
	disturlvl   DisturLevel
	fonts       []*truetype.Font
	size        image.Point
}

type StrType int

const (
	NUM   StrType = 0 // 数字
	LOWER         = 1 // 小写字母
	UPPER         = 2 // 大写字母
	ALL           = 3 // 全部
)

type DisturLevel int

const (
	NORMAL DisturLevel = 4
	MEDIUM             = 8
	HIGH               = 16
)

func New() *Captcha {
	c := &Captcha{
		disturlvl: NORMAL,
		size:      image.Point{82, 32},
	}
	c.frontColors = []color.Color{color.Black}
	c.bkgColors = []color.Color{color.White}
	return c
}

// AddFont 设置字体
func (c *Captcha) AddFont(path string) error {
	fontdata, erro := ioutil.ReadFile(path)
	if erro != nil {
		return erro
	}
	font, erro := freetype.ParseFont(fontdata)
	if erro != nil {
		return erro
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

func (c *Captcha) RandFont() *truetype.Font {
	return c.fonts[rand.Intn(len(c.fonts))]
}

func (c *Captcha) SetDisturbance(d DisturLevel) {
	if d > 0 {
		c.disturlvl = d
	}
}

func (c *Captcha) SetFrontColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.frontColors = c.frontColors[:0]
		for _, v := range colors {
			c.frontColors = append(c.frontColors, v)
		}
	}
}

func (c *Captcha) SetBkgColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.bkgColors = c.bkgColors[:0]
		for _, v := range colors {
			c.bkgColors = append(c.bkgColors, v)
		}
	}
}

func (c *Captcha) SetSize(w, h int) {
	if w < 48 {
		w = 48
	}
	if h < 20 {
		h = 20
	}
	c.size = image.Point{w, h}
}

// 绘制背景
func (c *Captcha) drawBkg(img *Image) {
	//ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 填充主背景色
	//bgcolorindex := ra.Intn(len(c.bkgColors))
	//bkg := image.NewUniform(c.bkgColors[bgcolorindex])
	//img.FillBkg(bkg)
	img.FillNoiseBkg(c.bkgColors)
}

// 绘制噪点
func (c *Captcha) drawNoises(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	dlen := int(c.disturlvl)
	// 绘制干扰斑点
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		r := ra.Intn(size.Y/20) + 1
		colorindex := ra.Intn(len(c.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, c.frontColors[colorindex])
	}

	// 绘制干扰线
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(size.Y) * o
		h := ra.Intn(size.Y/10) * o
		colorindex := ra.Intn(len(c.frontColors))
		img.DrawLine(x, y, x+w, y+h, c.frontColors[colorindex])
		colorindex++
	}

}

// 绘制文字
func (c *Captcha) drawString(img *Image, str string) {
	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	// 文字大小为图片高度的 0.65
	fsize := int(float64(size.Y) * 0.65)
	// 用于生成随机角度
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 文字之间的距离
	gap := size.X/(len(str)+1) - fsize/6
	// 文字在图形上的起点
	offset_y := int(float64(size.Y) * 0.2)
	offset_x := size.X/(len(str)+1)

	// 逐个绘制文字到图片上
	for i, char := range str {
		// 创建单个文字图片
		// 以高为尺寸创建正方形的图形
		str := NewImage(size.Y, size.Y)
		// 随机取一个前景色
		colorindex := r.Intn(len(c.frontColors))
		
		//随机取一个字体
		font := c.RandFont()
		str.DrawString(font, c.frontColors[colorindex], string(char), float64(fsize), offset_x, offset_y)

		// 转换角度后的文字图形

		//println(r.Float64())
		rs := str.Rotate(float64(r.Intn(60) - 30))
		s := rs.Bounds().Size()
		// 计算文字位置
		left := i*gap - (s.X - size.Y)
		top := size.Y - s.Y
		clip := image.Rect(left, top, left+s.X, top+s.Y)
		// 绘制到图片上
		draw.Draw(img, clip, rs, image.ZP, draw.Over)
	}
}

// Create 生成一个验证码图片
func (c *Captcha) Create(num int, t StrType) (*Image, string) {
	if num <= 0 {
		num = 4
	}
	dst := NewImage(c.size.X, c.size.Y)
	tmp := NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(tmp);
	str := string(c.randStr(num, int(t)))
	c.drawString(tmp, str)
	tmp.distortTo(dst,(4.0+rand.Float64()*4.0), (100.0+rand.Float64()*100.0))
	return dst, str
}

func (c *Captcha) CreateCustom(str string) *Image {
	if len(str) == 0 {
		str = "unkown"
	}
	dst := NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(dst);
	c.drawString(dst, str)
	return dst
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
