package captcha

import (
	"errors"
	"image/color"
	"io"
	"strings"

	"github.com/qAison/captcha/draw"
)

var (
	ErrNotFound = errors.New("captcha: id not found")

	globalStore    = NewMemoryStore(60)
	globalDraw     = draw.New()
	globalGenValue = draw.RandAlphaDigit
	globalIdLen    = 32
	globalValueLen = 4
)

func init() {
	globalDraw.SetSize(100, 30)
	globalDraw.Disturbance.SetNormal()
	globalDraw.SetFrontColor(
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{72, 72, 72, 255},
		color.RGBA{255, 0, 204, 255},
		color.RGBA{255, 102, 51, 255},
		color.RGBA{255, 153, 51, 255},
		color.RGBA{204, 0, 204, 255},
		color.RGBA{153, 153, 204, 255},
		color.RGBA{51, 102, 255, 255},
	)
	globalDraw.SetBackgroundColor(
		color.RGBA{255, 255, 255, 255}, // 白色
	)
}

//-----------------------------------------------------------------------------------------------------------//

// 设置 缓存
func SetStore(store Store) {
	globalStore = store
}

// 设置 验证码 画图 对象
func SetDraw(capt *draw.Draw) {
	globalDraw = capt
}

// 设置 生成 验证码 的值为 全数字
func SetGenValueDigit() {
	globalGenValue = draw.RandDigit
}

// 设置 生成 验证码 的值为 全字母
func SetGenValueAlpha() {
	globalGenValue = draw.RandAlpha
}

// 设置 生成 验证码 的值为 字母 + 数字
func SetGenValueAlphaDigit() {
	globalGenValue = draw.RandAlphaDigit
}

// 设置 生成 验证码 ID的长度
func SetGenIdLen(length int) {
	if length > 0 {
		globalIdLen = length
	}
}

// 设置 生成 验证码 值的长度
func SetGenValueLen(length int) {
	if length > 0 {
		globalValueLen = length
	}
}

//-----------------------------------------------------------------------------------------------------------//

// 默认 值 长度
func New() (id string) {
	return NewLen(globalValueLen)
}

// 指定 值 长度
func NewLen(length int) (id string) {
	id = draw.RandAlphaDigit(globalIdLen)
	globalStore.Set(id, globalGenValue(length))
	return
}

//-----------------------------------------------------------------------------------------------------------//

// 刷新指定 ID 的值【过期后将不能 重新刷新】
func Reload(id string) bool {
	value := globalStore.Get(id)
	if value == "" {
		return false
	}
	globalStore.Set(id, globalGenValue(len(value)))
	return true
}

// 刷新指定 ID 的值，不存在时，自动生成
func ReloadGen(id string, length int) {
	if value := globalStore.Get(id); value != "" {
		length = len(value)
	}
	globalStore.Set(id, globalGenValue(length))
}

//-----------------------------------------------------------------------------------------------------------//

// 获取 PNG图片 字节
func GetImage(id string) ([]byte, error) {
	value := globalStore.Get(id)
	if value == "" {
		return nil, ErrNotFound
	}

	return globalDraw.Create(value).EncodedPNG()
}

// 图片写入 IO
func WriteImage(w io.Writer, id string) error {
	value := globalStore.Get(id)
	if value == "" {
		return ErrNotFound
	}
	return globalDraw.Create(value).WriteTo(w)
}

//-----------------------------------------------------------------------------------------------------------//

// 校验【不区分大小写】
func Verify(id, value string) bool {
	if val := globalStore.Get(id); val == "" {
		return false
	} else {
		return strings.ToLower(value) == strings.ToLower(val)
	}
}

//-----------------------------------------------------------------------------------------------------------//
