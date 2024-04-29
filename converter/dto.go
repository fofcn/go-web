package converter

type CreationConvertDto struct {
	TaskId string `json:"task_id"`
}

type ConverterStatusDto struct {
	TaskId  string `json:"task_id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	SubType string `json:"sub_type"`
}
