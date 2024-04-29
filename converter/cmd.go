package converter

type CreationConvertCmd struct {
	Type    string `json:"type"`
	SubType string `json:"sub_type"`
	FileId  string `json:"file_id"`
}

type ConverterStatusCmd struct {
	TaskId string `json:"task_id"`
}
