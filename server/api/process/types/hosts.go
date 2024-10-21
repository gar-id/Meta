package types

// Hosts struct

type ServerResponse struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date    string          `json:"date"`
		Servers []ServerDetails `json:"servers"`
	}
}

// type HostDelete struct {
// 	HostIP  string `json:"host_ip" form:"host_ip" yaml:"hostIP"`
// 	Confirm bool   `json:"confirm" form:"confirm" yaml:"confirm"`
// }

type ServerDetails struct {
	Hostname         string `json:"hostname" yaml:"hostName"`
	PublicIP         string `json:"public_ip" yaml:"publicIp"`
	Environment      string `json:"environment" yaml:"environment"`
	Status           string `json:"status" yaml:"status"`
	ActiveConnection string `json:"active_connection" yaml:"activeConnection"`
}

// type HostTags struct {
// 	Region      string `json:"region" yaml:"region"`
// 	Project     string `json:"project" yaml:"project"`
// 	Environment string `json:"environment" yaml:"environment"`
// }
