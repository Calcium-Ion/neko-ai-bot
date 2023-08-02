package api

type Result struct {
	Code        int         `json:"code"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func ErrorUtil(description string) *Result {
	return &Result{
		Code:        500,
		Description: description,
		Data:        nil,
	}
}
