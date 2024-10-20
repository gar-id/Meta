package config

import (
	"os"
	"path"

	"MetaHandler/server/api"
	"MetaHandler/server/config/caches"
	"MetaHandler/server/databases"
	mqttMochi "MetaHandler/server/mqtt"

	"MetaHandler/tools"

	"gopkg.in/yaml.v2"
)

func LoadMetaHandlerConfig(configPath string) {
	// Load MetaHandler config file
	yamlFile, err := os.ReadFile(tools.DefaultString(configPath, path.Join("/", "etc", "Meta Handler", "config.yml")))
	if err != nil {
		tools.ZapLogger("both", "server").Fatal(err.Error())
	}

	// Parse MetaHandler config file into struct
	err = yaml.Unmarshal(yamlFile, &caches.MetaHandlerServer)
	if err != nil {
		tools.ZapLogger("both", "server").Fatal(err.Error())
	}
}

func Bootstrap(configPath string) {
	LoadMetaHandlerConfig(configPath)

	// Run migration
	databases.Bootstrap()
	tools.ZapLogger("both", "server").Info("Meta Handler database migration finished")

	// Run MetaHandler MQTT
	tools.ZapLogger("both", "server").Info("Starting Meta Handler MQTT Server")
	mqttMochi.Start()

	// Run MetaHandler API
	tools.ZapLogger("both", "server").Info("Starting Meta Handler RestAPI Server")
	api.Start()

}
