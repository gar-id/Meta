package types

type MetaHandlerAgent struct {
	AgentConfig struct {
		Server struct {
			HTTPPort               string `yaml:"httpPort" json:"http_port"`                        // fill with http port (example: 3000)
			HTTPHost               string `yaml:"httpHost" json:"http_host"`                        // fill with http host (example: 127.0.0.1)
			AgentToken             string `yaml:"agentToken" json:"agent_token"`                    // will be use for centralissh API agent token
			EncryptedAgentToken    string `yaml:"encryptedAgentToken" json:"encrypted_agent_token"` // will be use for encrypted centralissh API agent token
			MQTTPort               string `yaml:"mqttPort" json:"mqtt_port"`                        // fill with mqtt port (example: 3000)
			MQTTHost               string `yaml:"mqttHost" json:"mqtt_host"`                        // fill with mqtt host (example: 127.0.0.1)
			MQTTAgentUser          string `yaml:"mqttAgentUser" json:"mqtt_agent_user"`
			EncryptedMQTTAgentPass string `yaml:"encryptedMqttAgentPass" json:"encrypted_mqtt_agent_pass"`
			MQTTAgentPass          string `yaml:"mqttAgentPass" json:"mqtt_agent_pass"`
		} `json:"server" yaml:"server"`
		Meta struct {
			Hostname    string   `json:"hostname" yaml:"hostName"`
			IP          string   `json:"ip" yaml:"ip"`
			Services    []string `json:"services" yaml:"services"`
			Environment string   `json:"environment" yaml:"environment"`
			Stunnel     stunnel  `json:"stunnel" yaml:"stunnel"`
		} `json:"meta" yaml:"meta"`
		Log struct {
			Location string `yaml:"logLocation" json:"log_location"` // centralissh log location. example: /var/log/CentraliSSH/CentraliSSH.log
			Level    string `yaml:"logLevel" json:"log_level"`       // centralissh log level
		} `yaml:"log" json:"log"`
	} `yaml:"agentConfig" json:"agent_config"`
}

type stunnel struct {
	Directory   string `json:"directory" yaml:"directory"`
	ServiceName string `json:"service_name" yaml:"serviceName"`
}
