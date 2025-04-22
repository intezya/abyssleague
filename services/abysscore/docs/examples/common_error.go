package examples

type UnprocessableErrorResponse struct {
	Message string   `json:"message" example:"unprocessable entity"`
	Detail  string   `json:"detail"`
	Errors  []string `json:"errors"`
	Code    int      `json:"code" example:"422"`
	Path    string   `json:"path"`
}

type InternalServerErrorResponse struct {
	Message string `json:"message" example:"internal server error"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"500"`
	Path    string `json:"path"`
}

type TooManyRequestsResponse struct {
	Message string `json:"message" example:"too many requests"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"429"`
	Path    string `json:"path"`
}
