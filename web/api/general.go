package api

type SuccessResponse struct {
	Status string `json:"status" enums:"Success"`
}

type FailureResponse struct {
	Status string `json:"status" enums:"Failure"`
	Error  string `json:"error"`
}
