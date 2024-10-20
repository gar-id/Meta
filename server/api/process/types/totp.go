package types

// TOTP struct
type TOTPResponse struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date    string `json:"date"`
		Secret  string `json:"secret"`
		Message string `json:"message"`
		TOTPURI string `json:"totp_uri"`
	} `json:"data"`
}

type TOTP struct {
	Username string `json:"username" form:"username"`
	TOTPCode string `json:"totp_code" form:"totp_code"`
}
