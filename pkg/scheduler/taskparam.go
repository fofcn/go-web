package scheduler

// pdf
type PdfPerPageSplitterParam struct {
	FileId       string `json:"file_id"`
	PagesPerFile int    `json:"pages_per_file"`
}
