package file

import "github.com/gin-gonic/gin"

func InitRouter(public *gin.RouterGroup) {
	csvrouter := NewCsvRouter(NewCsvService())
	public.POST("/split", csvrouter.SplitCsv)
}

type CsvRouter struct {
	csvservice CsvService
}

func NewCsvRouter(csvservice CsvService) *CsvRouter {
	return &CsvRouter{
		csvservice: csvservice,
	}
}

func (cr *CsvRouter) SplitCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	path := "./tmp/"
	err = c.SaveUploadedFile(file, path+file.Filename)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if err := cr.csvservice.SplitCsv(path + file.Filename); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "ok",
	})
}
