package models

// PostData POST JSON 表单的数据
type PostData[T any] struct {
	Data T `json:"data"`
}
