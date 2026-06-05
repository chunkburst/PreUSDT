package router

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/chunkburst/PreUSDT/app/log"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/chunkburst/PreUSDT/static"
	"github.com/gin-gonic/gin"
)

func staticInit(e *gin.Engine) {
	customPath := model.GetK(model.PaymentStaticPath)
	if customPath != "" && utils.IsExist(customPath) {
		initCustomPayment(e, customPath)

		return
	}

	initDefaultPayment(e)
}

func initCustomPayment(e *gin.Engine, path string) {
	tmpl := template.New("customer").Funcs(template.FuncMap{"toJson": toJson})

	template.Must(tmpl.ParseGlob(filepath.Join(path, "views", "*.html")))
	template.Must(tmpl.ParseFS(static.Secure, "secure/secure.html"))
	e.SetHTMLTemplate(tmpl)

	e.StaticFS("/payment/assets", http.Dir(filepath.Join(path, "assets")))
	e.StaticFS("/secure/assets", http.FS(subFS(static.Secure, "secure/assets")))

	log.Info("成功注册自定义静态资源路径：", path)
}

func initDefaultPayment(e *gin.Engine) {
	tmpl := template.New("default").Funcs(template.FuncMap{"toJson": toJson})

	template.Must(tmpl.ParseFS(static.Secure, "secure/secure.html"))
	template.Must(tmpl.ParseFS(static.Payment, "payment/views/*.html"))
	e.SetHTMLTemplate(tmpl)

	e.StaticFS("/payment/assets", http.FS(subFS(static.Payment, "payment/assets")))
	e.StaticFS("/secure/assets", http.FS(subFS(static.Secure, "secure/assets")))
}

func subFS(src fs.FS, dir string) fs.FS {
	sub, _ := fs.Sub(src, dir)

	return sub
}

func toJson(v any) template.JS {
	data, err := json.Marshal(v)
	if err != nil {
		return "null"
	}

	return template.JS(data)
}
