package pdf

import (
	"go-web/pkg/config"
	"go-web/pkg/scheduler"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func InitRouter(public *gin.RouterGroup) {
	pdfRouter := NewPdfRouter(NewPdfService())
	public.POST("/pdf/split", pdfRouter.SplitPdf)
	public.GET("/task/:id", pdfRouter.GetTaskResult)
}

type PdfRouter struct {
	pdfservice PdfService
	scheduler  *scheduler.Scheduler
}

func NewPdfRouter(pdfservice PdfService) *PdfRouter {
	return &PdfRouter{
		pdfservice: pdfservice,
		scheduler:  scheduler.GetScheduler(config.GetScheduler()),
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

func (cr *PdfRouter) GetTaskResult(c *gin.Context) {
	taskId := c.Param("id")
	if len(taskId) == 0 {
		c.JSON(400, gin.H{
			"msg": "task id is empty",
		})
		return
	}

	if dto, err := cr.scheduler.GetTaskStatus(taskId); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	} else {
		if dto.Data != nil {
			path, ok := dto.Data.(string)
			if ok {
				newPath := strings.Replace(path, "/home/xiaosi", "", 1)
				newDto := &scheduler.TaskResult{
					TaskId: dto.TaskId,
					Status: dto.Status,
					Data:   newPath,
				}
				c.JSON(200, newDto)
				return
			}
		}
		c.JSON(200, dto)
	}
}
