# golang实现的验证码 golang captcha

**不依赖任何第三方图形库 Does not rely on third party graphics library **

之前找了很多go的验证码实现,发现要么太简单,要么需要依赖很多东西。如`imagemagic`。所以自己没事用go实现了一个验证码。

## 优点

1. 使用简单
2. 不依赖第三方图形库 直接go get 就Ok
3. 丰富自定义设置(字体,多颜色,验证码大小,文字模式,文字数量)
4. 中文注释:) (你能看懂，自己修改成自己用的)


## demo

![color](http://afocus.github.io/captcha/demo1.png)

![sample](http://afocus.github.io/captcha/demo2.png)

## 使用 Start using it

Download and install it:
```
go get github.com/afocus/captcha
```

#### 最简单的示例 sample use

```go
cap = captcha.New()
cap.SetFont("comic.ttf")
img,str := cap.Create(100,32)
```

#### 设置 set options

```go
cap = captcha.New()
cap.SetFont("comic.ttf")
cap.SetOpt(
    4, // 验证码文字数量 string length
    captcha.ALL, // 文字模式 0:纯数字 1:小写字母 2:大小写字母 3:数字+大小写字母 string modal
    captcha.Normal, // 干扰强度 disturbance
    captcha.Color{255,0,0}, // 文字以及干扰线颜色，可以设置多个 font color
    captcha.Color{0,0,255},
)

img,str := cap.Create(100,32)
img1,str1 := cap.Create(200,64)
```

#### for web

look `examples/main.go`




