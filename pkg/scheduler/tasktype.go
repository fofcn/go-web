package scheduler

type TaskType string
type SubTaskType string

const (
	TaskTypePdf   TaskType = "pdf"
	TaskTypeCsv   TaskType = "csv"
	TaskTypeImage TaskType = "image"
)

// pdf sub task type
const (
	SubTaskTypePdf2Csv             SubTaskType = "pdf2csv"
	SubTaskTypePdf2Img             SubTaskType = "pdf2img"
	SubTaskTypePdfSplitter         SubTaskType = "pdfsplitter"
	SubTaskTypePdfMerger           SubTaskType = "pdfmerger"
	SubTaskTypePdfRotator          SubTaskType = "pdfrotator"
	SubTaskTypePdfWatermarkRemover SubTaskType = "pdfwatermarkremover"
	SubTaskTypePdfWatermarkAdder   SubTaskType = "pdfwatermarkadder"
)

var validTaskSubTypeMapping = map[TaskType][]SubTaskType{
	TaskTypePdf: {
		SubTaskTypePdf2Csv,
		SubTaskTypePdf2Img,
		SubTaskTypePdfSplitter,
		SubTaskTypePdfMerger,
		SubTaskTypePdfRotator,
		SubTaskTypePdfWatermarkRemover,
		SubTaskTypePdfWatermarkAdder,
	},
	TaskTypeCsv: {
		SubTaskTypePdf2Csv,
	},
}

func IsValidTask(taskType TaskType, subTaskType SubTaskType) bool {
	validSubTypes, ok := validTaskSubTypeMapping[taskType]
	if !ok {
		return false
	}
	for _, validSubType := range validSubTypes {
		if subTaskType == validSubType {
			return true
		}
	}
	return false
}
