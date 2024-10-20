package types

// General struct
type General struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date    string `json:"date"`
		Message string `json:"message"`
	} `json:"data"`
}
type Welcome struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date    string `json:"date"`
		Message string `json:"message"`
		Version string `json:"version"`
	} `json:"data"`
}

type ErrorHandler struct {
	HTTP_Code int    `json:"http_code"`
	ClientIP  string `json:"client_ip"`
	Status    string `json:"status"`
}
