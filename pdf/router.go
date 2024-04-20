package pdf

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func InitRouter(public *gin.RouterGroup) {
	pdfRouter := NewPdfRouter(NewPdfService())
	public.POST("/pdf/split", pdfRouter.SplitPdf)
}

type PdfRouter struct {
	pdfservice PdfService
}

func NewPdfRouter(pdfservice PdfService) *PdfRouter {
	return &PdfRouter{
		pdfservice: pdfservice,
	}
}

func (cr *PdfRouter) SplitPdf(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	pages_per_file := c.Request.FormValue("pages_per_file")

	path := "/home/xiaosi/tmp/"
	newfilename, err := uuid.NewUUID()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		return
	}

	absServerPath := path + newfilename.String()
	err = c.SaveUploadedFile(file, absServerPath)
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

	if dto, err := cr.pdfservice.SplitPdf(file.Filename, absServerPath, ipages_per_file); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	} else {
		c.JSON(200, dto)
	}

}
