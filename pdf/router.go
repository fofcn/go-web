package pdf

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitRouter(public *gin.RouterGroup) {
	pdfRouter := NewPdfRouter(NewPdfService())
	public.POST("/pdf/split", pdfRouter.SplitCsv)
}

type PdfRouter struct {
	pdfservice PdfService
}

func NewPdfRouter(pdfservice PdfService) *PdfRouter {
	return &PdfRouter{
		pdfservice: pdfservice,
	}
}

func (cr *PdfRouter) SplitCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	pages_per_file := c.Request.FormValue("pages_per_file")

	path := "./tmp/"
	err = c.SaveUploadedFile(file, path+file.Filename)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ipages_per_file, err := strconv.Atoi(pages_per_file)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if err := cr.pdfservice.SplitPdf(path+file.Filename, ipages_per_file); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "ok",
	})
}
