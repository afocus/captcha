# golang实现的验证码 golang captcha
丰富自定义设置(字体,多颜色,验证码大小,文字模式,文字数量,干扰强度)

## demo

![color](https://raw.githubusercontent.com/qAison/captcha/master/_examples/1.png)![sample](https://raw.githubusercontent.com/qAison/captcha/master/_examples/2.png)![color2](https://raw.githubusercontent.com/qAison/captcha/master/_examples/3.png)

## 使用 Start using it

Download and install it:
```
go get github.com/qAison/captcha
```


#### 最简单的示例 sample use

```go
// 验证码 ID
id := captcha.New()

oFile, _ := os.Create(id+".png")
defer oFile.Close()

// 写入文件
captcha.WriteImage(oFile, id)

```

#### 设置 set options

```go
cap := draw.New()

// 设置 图片大小
cap.SetSize(100, 30)

// 设置 干扰度
cap.Disturbance.SetNormal()

// 设置 字体
cap.SetFrontColor("comic.ttf", "xxx.ttf")

// 设置 字体颜色
cap.SetFrontColor(color.RGBA{255, 255, 255, 255})

// 设置 多个 背景色，将随机使用
cap.SetBackgroundColor(
    color.RGBA{255, 0, 0, 255},
    color.RGBA{0, 0, 255, 255},
    color.RGBA{0, 153, 0, 255},
    color.RGBA{185, 123, 131, 255},
    color.RGBA{185, 123, 43, 255},
    color.RGBA{82, 146, 114, 255},
    color.RGBA{82, 69, 114, 255},
    color.RGBA{22, 69, 114, 255},
)

// 创建 全数字 图片
img, val := cap.CreateDigit(4)

// 创建 全字母 图片
img, val := cap.CreateAlpha(4) 

// 创建 字母 + 数字 图片
img, val := cap.CreateAlphaDigit(4) 

// 创建 自定义字符 图片
img := cap.Create("abc123") 
		

#### 网站中如果使用? how to use for web

look `_examples/web/main.go`




