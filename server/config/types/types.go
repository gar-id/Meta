package types

type ServerConfig struct {
	MetaHandlerServer struct {
		API struct {
			HTTPPort            string `yaml:"httpPort" json:"http_port"`                        // fill with http port (example: 3000)
			HTTPHost            string `yaml:"httpHost" json:"http_host"`                        // fill with http port (example: 127.0.0.1)
			MainToken           string `yaml:"mainToken" json:"main_token"`                      // will be use for meta handler API main token
			AgentToken          string `yaml:"agentToken" json:"agent_token"`                    // will be use for meta handler API agent token
			EncryptedMainToken  string `yaml:"encryptedMainToken" json:"encrypted_main_token"`   // will be use for meta handler API main token
			EncryptedAgentToken string `yaml:"encryptedAgentToken" json:"encrypted_agent_token"` // will be use for meta handler API agent token
		} `yaml:"api" json:"api"`
		MQTT struct {
			Port               string `yaml:"mqttPort" json:"mqtt_port"` // fill with mqtt port (example: 3000)
			Host               string `yaml:"mqttHost" json:"mqtt_host"` // fill with mqtt port (example: 127.0.0.1)
			AdminUser          string `yaml:"mqttAdminUser" json:"mqtt_admin_user"`
			AdminPass          string `yaml:"mqttAdminPass" json:"mqtt_admin_pass"`
			EncryptedAdminPass string `yaml:"encryptedMqttAdminPass" json:"encrypted_mqtt_admin_pass"`
			AgentUser          string `yaml:"mqttAgentUser" json:"mqtt_agent_user"`
			AgentPass          string `yaml:"mqttAgentPass" json:"mqtt_agent_pass"`
			EncryptedAgentPass string `yaml:"encryptedMqttAgentPass" json:"encrypted_mqtt_agent_pass"`
		} `yaml:"mqtt" json:"mqtt"`
		Log struct {
			Level    string `yaml:"logLevel" json:"log_level"`       // meta handler log level. You can choose between debug, info, warning, error, panic, fatal
			Location string `yaml:"logLocation" json:"log_location"` // meta handler log location. example: /var/log/meta handler/meta handler.log
		} `yaml:"log" json:"log"`
		Database struct {
			Driver   string `yaml:"databaseDriver" json:"database_driver"`      // meta handler database driver. Set one of those: postgres, mysql, sqlite
			TimeZone string `yaml:"databaseTimeZone" jsone:"database_timezone"` // meta handler database timezone
			Postgres struct {
				Host              string `yaml:"postgresqlHost" json:"postgresql_host"`                            // Default is 127.0.0.1
				Port              string `yaml:"postgresqlPort" json:"postgresql_port"`                            // Default is 5432
				User              string `yaml:"postgresqlUser" json:"postgresql_user"`                            // Default is postgres
				Password          string `yaml:"postgresqlPassword" json:"postgresql_password"`                    // Default is postgres
				EncryptedPassword string `yaml:"encryptedPostgresqlPassword" json:"encrypted_postgresql_password"` // Default is postgres
				DB                string `yaml:"postgresqlDB" json:"postgresql_db"`                                // default is postgres
				SSLMode           string `yaml:"postgresqlSSLMode" json:"postgresql_ssl_mode"`                     // default is disable
				TimeZone          string `yaml:"postgresqlTimeZone" json:"postgresql_timezone"`                    // default to Etc/UTC
			} `yaml:"postgres" json:"postgres"`
			Mysql struct {
				Host              string `yaml:"mysqlHost" json:"mysql_host"`                            // Default is 127.0.0.1
				Port              string `yaml:"mysqlPort" json:"mysql_port"`                            // Default is 5432
				User              string `yaml:"mysqlUser" json:"mysql_user"`                            // Default is root
				Password          string `yaml:"mysqlPassword" json:"mysql_password"`                    // Default is empty
				EncryptedPassword string `yaml:"encryptedMysqlPassword" json:"encrypted_mysql_password"` // Default is empty
				DB                string `yaml:"mysqlDB" json:"mysql_db"`                                // default is meta handler
				Charset           string `yaml:"mysqlCharset" json:"mysql_charset"`                      // default is utf8mb4
				ParseTime         string `yaml:"mysqlParseTime" json:"mysql_parse_time"`                 // default is True
				Loc               string `yaml:"mysqlLoc" json:"mysql_loc"`                              // default to Local
			} `yaml:"mysql" json:"mysql"`
			Sqlite struct {
				Location string `yaml:"sqliteLocation" json:"sqlite_location"` // default to /opt/meta handler/database/meta handler.db
			} `yaml:"sqlite" json:"sqlite"`
		} `yaml:"database" json:"database"`
	} `yaml:"metaHandler" json:"meta_handler"`
}
