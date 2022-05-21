package model

type Response struct {
	Status ResponseStatus `json:"status"`
	Data   interface{}    `json:"data"`
}

type ResponseStatus string

const (
	ResponseSuccess ResponseStatus = "success"
	ResponseError   ResponseStatus = "error"
)
