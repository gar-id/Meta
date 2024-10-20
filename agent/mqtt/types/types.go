package types

type MQTTPayload struct {
	Type    string      `json:"type"` // type = stunnel, service
	Stunnel MQTTStunnel `json:"stunnel"`
	Service MQTTService `json:"service"`
}

type MQTTStunnel struct {
	Action     string `json:"action"` // action = [failover]
	FileConfig string `json:"file_config"`
	// File config example :
	// [my-stunnel-conf]
	// Client = yes
	// Accept = 127.0.0.1:23
	// Connect =84.15.X.X:2030
}

type MQTTService struct {
	ServiceName string `json:"service_name"` // service name example = [stunnel,meta-server]
	Action      string `json:"action"`       // action = [stop,start,status,restart]
}
