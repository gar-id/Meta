package mqttAgent

import (
	"crypto/md5"
	"fmt"

	"MetaHandler/agent/caches"

	"MetaHandler/tools"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connection() mqtt.Client {
	// Create mqtt connection options
	localIPBytes := []byte(tools.GetLocalIP())
	md5IP := fmt.Sprintf("%x", md5.Sum(localIPBytes))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%v:%v", caches.MetaHandlerAgent.AgentConfig.Server.MQTTHost, caches.MetaHandlerAgent.AgentConfig.Server.MQTTPort))
	opts.SetClientID(fmt.Sprintf("Meta_Handler_MQTT_Agent_%v", md5IP))
	opts.SetUsername(tools.DefaultString(caches.MetaHandlerAgent.AgentConfig.Server.MQTTAgentUser, "agent"))
	opts.SetPassword(tools.DefaultString(caches.MetaHandlerAgent.AgentConfig.Server.MQTTAgentPass, "agent"))
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.AutoReconnect = true
	opts.ConnectRetry = true

	// Check messageHandler
	opts.SetDefaultPublishHandler(messagePubHandler)

	// Start connection
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Return
	return client
}
