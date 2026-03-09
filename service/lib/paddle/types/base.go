package types

import "net/http"

type ErrorDetailsItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Type             string             `json:"type"`
	Code             string             `json:"code"`
	Detail           string             `json:"detail"`
	DocumentationUrl string             `json:"documentation_url"`
	ErrorDetails     []ErrorDetailsItem `json:"errors"`
}

type Meta struct {
	RequestId string `json:"request_id"`
}

// An unsuccessful call to the Dashboard API will return a 200 response containing
// a field success set to false. Additionally an error object will be returned,
// containing a code referencing the error, and a message in a human-readable format.

// Api doc:https://developer.paddle.com/api-reference/about/errors
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error

	Meta  Meta   `json:"meta"`
	Error *Error `json:"error"`
}

// Api doc:https://developer.paddle.com/api-reference/about/success-responses
type SuccessResponse[T any] struct {
	Meta Meta `json:"meta"`
	Data T    `json:"data"`
}
