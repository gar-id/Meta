package types

// Group struct
type GroupPayload struct {
	GroupName   string           `json:"group_name" form:"group_name" yaml:"groupName"`
	Role        string           `json:"role" yaml:"role"`
	TargetHosts GroupTargetHosts `json:"target_hosts" yaml:"targetHosts"`
}
type GroupTargetHosts struct {
	HostIDs []string `json:"host_ids" yaml:"hostIDs"`
	HostIPs []string `json:"host_ips" yaml:"hostIPs"`
	Action  string   `json:"action" yaml:"action"`
}

type GroupResponse struct {
	HTTP_Code int    `json:"http_code"`
	Status    string `json:"status"`
	ClientIP  string `json:"client_ip"`
	Data      struct {
		Date  string      `json:"date"`
		Group []GroupData `json:"group"`
	} `json:"data"`
}

type GroupData struct {
	GroupName   string        `json:"group_name" yaml:"groupName"`
	Role        string        `json:"role" yaml:"role"`
	TargetHosts []HostDetails `json:"target_hosts" yaml:"targetHosts"`
}

type GroupDelete struct {
	GroupName string `json:"group_name" form:"group_name" yaml:"groupName"`
	Confirm   bool   `json:"confirm" form:"confirm" yaml:"confirm"`
}
