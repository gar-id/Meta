package types

import "time"

// UserUpdate struct
type UserDetails struct {
	Username         string          `json:"username" form:"username" yaml:"userName"`
	MFAEnabled       bool            `json:"mfa_enabled" form:"mfa_enabled" yaml:"mfaEnables"`
	Role             string          `json:"role" form:"role" yaml:"role"`
	Disabled         bool            `json:"disabled" form:"disabled" yaml:"disabled"`
	PublicKey        string          `json:"public_key" form:"public_key" yaml:"publicKey"`
	PublicKeyExpired time.Time       `json:"public_key_expired" form:"public_key_expired" yaml:"publicKeyExpired"`
	Token            string          `json:"token"`
	Group            []UserGroupJSON `json:"group" yaml:"group"`
}
type UserGroupJSON struct {
	GroupName string `json:"group_name" yaml:"groupName"`
	Sudo      string `json:"sudo" yaml:"sudo"`
}

type UserResponse struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date  string        `json:"date"`
		Users []UserDetails `json:"users"`
	}
}

type UserDelete struct {
	Username string `json:"username" form:"username" yaml:"userName"`
	Confirm  bool   `json:"confirm" form:"confirm" yaml:"confirm"`
}
