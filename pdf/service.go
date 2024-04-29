package pdf

import (
	"context"
	"go-web/pkg/config"
	"go-web/pkg/scheduler"
)

type PdfService interface {
	SplitPdf(filename string, filepath string, pages_per_file int) (*PdfSumitTaskDto, error)
}

type pdfservice struct {
	scheduler *scheduler.Scheduler
}

type PdfSplitTask struct {
	PagesPerFile  int    `json:"pages_per_file"`
	FilePath      string `json:"pdf_path"`
	TaskType      string `json:"task_type"`
	SubTaskType   string `json:"sub_task_type"`
	FileOrginName string
}

func NewPdfService() PdfService {
	context.Background()
	return &pdfservice{
		scheduler: scheduler.GetScheduler(config.GetScheduler()),
	}
}

func (p *pdfservice) SplitPdf(filename, filepath string, pages_per_file int) (*PdfSumitTaskDto, error) {
	return nil, nil
}
