package draw

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//-----------------------------------------------------------------------------------------------------------//

// 干扰度
type disturbance struct {
	value int // 必须 >0 的值
}

func (self *disturbance) Value() int {
	return self.value
}

func (self *disturbance) Set(value int) {
	if value > 0 {
		self.value = value
	}
}

func (self *disturbance) SetNormal() {
	self.value = 4
}

func (self *disturbance) SetMedium() {
	self.value = 8
}

func (self *disturbance) SetHigh() {
	self.value = 16
}

//-----------------------------------------------------------------------------------------------------------//

// 画验证码
type Draw struct {
	Disturbance      disturbance      // 干扰度
	frontColors      []color.Color    // 字体色列表
	backgroundColors []color.Color    // 背景色列表
	fonts            []*truetype.Font // 字体
	size             image.Point      // 图片大小
}

// 验证码
func New() *Draw {
	captcha := &Draw{
		frontColors:      []color.Color{color.Black},      // 黑色
		backgroundColors: []color.Color{color.White},      // 白色
		fonts:            []*truetype.Font{globalDefFont}, // 默认字体
		size:             image.Point{150, 30},
	}
	captcha.Disturbance.SetNormal()
	return captcha
}

// 创建 数字 验证码 图片
func (self *Draw) CreateDigit(length int) (*Image, string) {
	if length <= 0 {
		length = 4
	}
	str := RandDigit(length)
	return self.create(str), str
}

// 创建 字母 验证码 图片
func (self *Draw) CreateAlpha(length int) (*Image, string) {
	if length <= 0 {
		length = 4
	}
	str := RandAlpha(length)
	return self.create(str), str
}

// 创建 字母+数字 验证码 图片
func (self *Draw) CreateAlphaDigit(length int) (*Image, string) {
	if length <= 0 {
		length = 4
	}
	str := RandAlphaDigit(length)
	return self.create(str), str
}

// 创建 指定字符 验证码 图片
func (self *Draw) Create(str string) *Image {
	if len(str) == 0 {
		str = "unkown"
	}
	return self.create(str)
}

//-----------------------------------------------------------------------------------------------------------//

// 创建 指定字符串
func (self *Draw) create(str string) *Image {
	dst := NewImage(self.size.X, self.size.Y)
	self.drawBackground(dst)
	self.drawDisturbance(dst)
	self.drawString(dst, str)
	return dst
}

// 添加 字体 可以设置多个
func (self *Draw) AddFont(paths ...string) error {
	for _, path := range paths {
		fontData, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		font, err := freetype.ParseFont(fontData)
		if err != nil {
			return err
		}
		self.fonts = append(self.fonts, font)
	}
	return nil
}

// 添加 二进制 内容的字体
func (self *Draw) AddFontFromBytes(contents []byte) error {
	font, err := freetype.ParseFont(contents)
	if err != nil {
		return err
	}
	self.fonts = append(self.fonts, font)
	return nil
}

// 设置 字体颜色【覆盖】
func (self *Draw) SetFrontColor(colors ...color.Color) *Draw {
	if len(colors) > 0 {
		self.frontColors = colors
	}
	return self
}

// 设置 背景颜色【覆盖】
func (self *Draw) SetBackgroundColor(colors ...color.Color) *Draw {
	if len(colors) > 0 {
		self.backgroundColors = colors
	}
	return self
}

// 设置 图片大小
func (self *Draw) SetSize(width, height int) *Draw {
	if width < 48 {
		width = 48
	}
	if height < 20 {
		height = 20
	}
	self.size = image.Point{width, height}
	return self
}

// 绘制 背景
func (self *Draw) drawBackground(img *Image) {
	bgcolorindex := rand.Intn(len(self.backgroundColors)) // 随机背景色
	bkg := image.NewUniform(self.backgroundColors[bgcolorindex])
	img.FillBackground(bkg) //填充主背景色
}

// 绘制 噪点
func (self *Draw) drawDisturbance(img *Image) {
	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	dlen := self.Disturbance.Value()

	// 绘制干扰斑点
	for i := 0; i < dlen; i++ {
		x := rand.Intn(size.X)
		y := rand.Intn(size.Y)
		r := rand.Intn(size.Y/20) + 1
		colorindex := rand.Intn(len(self.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, self.frontColors[colorindex])
	}

	// 绘制干扰线
	for i := 0; i < dlen; i++ {
		x := rand.Intn(size.X)
		y := rand.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := rand.Intn(size.Y) * o
		h := rand.Intn(size.Y/10) * o
		colorIndex := rand.Intn(len(self.frontColors))
		img.DrawLine(x, y, x+w, y+h, self.frontColors[colorIndex])
		colorIndex++
	}
}

// 绘制文字
func (self *Draw) drawString(img *Image, str string) {
	tmp := NewImage(self.size.X, self.size.Y)

	fsize := int(float64(self.size.Y) * 0.6) // 文字大小为图片高度的 0.6
	padding := fsize / 4                     // 文字之间的距离，左右各留文字的1/4大小为内部边距
	gap := (self.size.X - padding*2) / (len(str))

	// 逐个绘制文字到图片上
	for i, char := range str {
		str := NewImage(fsize, fsize)                  // 创建单个文字图片，以文字为尺寸创建正方形的图形
		colorIndex := rand.Intn(len(self.frontColors)) // 随机取一个前景色

		// 随机取一个字体
		font := self.fonts[rand.Intn(len(self.fonts))]
		str.DrawString(font, self.frontColors[colorIndex], string(char), float64(fsize))

		// 转换角度后的文字图形
		rs := str.Rotate(float64(rand.Intn(40) - 20))

		// 计算文字位置
		s := rs.Bounds().Size()
		left := i*gap + padding
		top := (self.size.Y - s.Y) / 2

		// 绘制到图片上
		draw.Draw(tmp, image.Rect(left, top, left+s.X, top+s.Y), rs, image.ZP, draw.Over)
	}

	if self.size.Y >= 48 { // 高度大于 48 添加波纹 小于 48 波纹影响用户识别
		tmp.distortTo(float64(fsize)/10, 200.0)
	}

	draw.Draw(img, tmp.Bounds(), tmp, image.ZP, draw.Over)
}

//-----------------------------------------------------------------------------------------------------------//
