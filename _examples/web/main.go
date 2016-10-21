package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/qAison/captcha"
)

var (
	formTemplate = template.Must(template.New("example").Parse(formTemplateSrc))
)

const (
	CST_ID_LEN  = 32
	CST_VAL_LEN = 4
)

func init() {
	captcha.SetGenIdLen(CST_ID_LEN)
}

//-----------------------------------------------------------------------------------------------------------//

type captchaHandler struct {
	length int
}

func Server(length int) http.Handler {
	return &captchaHandler{
		length: length,
	}
}

func (self *captchaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, id := path.Split(r.URL.Path)
	if id = strings.TrimSpace(id); len(id) != CST_ID_LEN {
		http.NotFound(w, r)
		return
	}

	if r.FormValue("reload") != "" {
		captcha.ReloadGen(id, self.length)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")

	buf := new(bytes.Buffer)
	if err := captcha.WriteImage(buf, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.ServeContent(w, r, id+".png", time.Now(), bytes.NewReader(buf.Bytes()))
	}
}

//-----------------------------------------------------------------------------------------------------------//

func showFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	d := struct {
		CaptchaId string
	}{
		captcha.NewLen(CST_VAL_LEN),
	}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.Verify(r.FormValue("captchaId"), r.FormValue("captchaValue")) {
		io.WriteString(w, "xxx 验证码错误 xxx")
	} else {
		io.WriteString(w, "~~~ 验证码正确  ~~~")
	}
	io.WriteString(w, "<br><br><a href='/'>继续再尝试</a>")
}

// ulimit -HSn 65535
func main() {
	http.HandleFunc("/", showFormHandler)
	http.HandleFunc("/process", processFormHandler)
	http.Handle("/captcha/", Server(CST_VAL_LEN))
	addr := "0.0.0.0:8666"
	fmt.Println("Server is at " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

const formTemplateSrc = `<!doctype html>
<head><title>验证码示例</title></head>
<body>
<script>
function setSrcQuery(e, q) {
	var src  = e.src;
	var p = src.indexOf('?');
	if (p >= 0) {
		src = src.substr(0, p);
	}
	e.src = src + "?" + q
}

function reload() {
	setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
	return false;
}
</script>

<form action="/process" method=post>
	<p>请输入下面的验证码:</p>
	<p><img id=image src="/captcha/{{.CaptchaId}}" alt="验证码图片"></p>

	<a href="#" onclick="reload()">刷新</a>

	<input type=hidden name=captchaId value="{{.CaptchaId}}"><br>
	<input name=captchaValue>
	<input type=submit value=提交>
</form>
`
