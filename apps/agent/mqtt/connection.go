package mqttCSSH

import (
	"crypto/md5"
	"fmt"

	"github.com/gar-id/centralissh/agent/datas"
	"github.com/gar-id/centralissh/tools"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connection() mqtt.Client {
	// Create mqtt connection options
	localIPBytes := []byte(tools.GetLocalIP())
	md5IP := fmt.Sprintf("%x", md5.Sum(localIPBytes))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%v:%v", datas.CentraliSSHAgentConfig.CentraliSSHAgent.Server.MQTTHost, datas.CentraliSSHAgentConfig.CentraliSSHAgent.Server.MQTTPort))
	opts.SetClientID(fmt.Sprintf("CentraliSSH_MQTT_Agent_%v", md5IP))
	opts.SetUsername(tools.DefaultString(datas.CentraliSSHAgentConfig.CentraliSSHAgent.Server.MQTTAgentUser, "agent"))
	opts.SetPassword(tools.DefaultString(datas.CentraliSSHAgentConfig.CentraliSSHAgent.Server.MQTTAgentPass, "agent"))
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
