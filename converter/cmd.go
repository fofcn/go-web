package converter

type CreationConvertCmd struct {
	Type    string      `json:"type"`
	SubType string      `json:"sub_type"`
	FileId  string      `json:"file_id"`
	Params  interface{} `json:"params"`
}

type ConverterStatusCmd struct {
	TaskId string `json:"task_id"`
}

type FileConvertUserDefCmd struct {
	FileIds []string               `json:"file_id"`
	Params  map[string]interface{} `json:"params"`
}
