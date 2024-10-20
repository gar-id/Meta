package types

// Hosts struct

type HostResponse struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date  string        `json:"date"`
		Hosts []HostDetails `json:"hosts"`
	}
}

type HostDelete struct {
	HostIP  string `json:"host_ip" form:"host_ip" yaml:"hostIP"`
	Confirm bool   `json:"confirm" form:"confirm" yaml:"confirm"`
}

type HostDetails struct {
	HostTags HostTags `json:"host_tags" yaml:"hostTags"`
	Hostname string   `json:"hostname" yaml:"hostName"`
	Host     string   `json:"host" yaml:"host"`
	Port     string   `json:"port" yaml:"port"`
}

type HostTags struct {
	Region      string `json:"region" yaml:"region"`
	Project     string `json:"project" yaml:"project"`
	Environment string `json:"environment" yaml:"environment"`
}
