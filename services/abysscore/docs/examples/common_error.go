package examples

type UnprocessableEntityResponse struct {
	Message string `json:"message" example:"unprocessable entity"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"422"`
	Path    string `json:"path"`
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

type ForbiddenByAccessLevelResponse struct {
	Message string `json:"message" example:"insufficient access level"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"403"`
	Path    string `json:"path"`
}

type BadRequestResponse struct {
	Message string   `json:"message" example:"bad request"`
	Detail  string   `json:"detail"`
	Errors  []string `json:"errors"`
	Code    int      `json:"code" example:"400"`
	Path    string   `json:"path"`
}
