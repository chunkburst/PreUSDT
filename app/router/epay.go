package router

import (
	"github.com/chunkburst/PreUSDT/app/handler/epay"
	"github.com/gin-gonic/gin"
)

func epayInit(engine *gin.Engine) {
	epHdr := new(epay.Epay)
	{
		engine.POST("/submit.php", epHdr.Submit)
		engine.GET("/submit.php", epHdr.Submit)
	}
}
