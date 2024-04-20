package pdf

import "go-web/pkg/scheduler"

type PdfService interface {
	SplitPdf(filepath string, pages_per_file int) (*PdfSumitTaskDto, error)
}

type pdfservice struct {
	scheduler *scheduler.Scheduler
}

type PdfSplitTask struct {
	PagesPerFile int    `json:"pages_per_file"`
	FilePath     string `json:"pdf_path"`
	TaskType     string `json:"task_type"`
	SubTaskType  string `json:"sub_task_type"`
}

func NewPdfService() PdfService {
	return &pdfservice{
		scheduler: scheduler.GetScheduler("rr"),
	}
}

func (p *pdfservice) SplitPdf(filepath string, pages_per_file int) (*PdfSumitTaskDto, error) {
	pdfSplitTask := &PdfSplitTask{
		FilePath:     filepath,
		TaskType:     "pdf",
		SubTaskType:  "splitter",
		PagesPerFile: pages_per_file,
	}
	task := scheduler.NewTask(scheduler.TaskTypeCSVSplitter, pdfSplitTask)
	_, err := p.scheduler.Schedule(task)
	if err != nil {
		return nil, err
	}

	return &PdfSumitTaskDto{
		TaskId: task.GetId(),
	}, nil
}
