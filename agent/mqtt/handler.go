package mqttAgent

import (
	"encoding/json"
	"fmt"

	"MetaHandler/agent/mqtt/types"
	"MetaHandler/agent/windows/services"

	"MetaHandler/tools"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message Subscriber handler
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	go func() {
		// Unmarshal payload
		jsonByte := []byte(msg.Payload())
		var payload types.MQTTPayload
		err := json.Unmarshal(jsonByte, &payload)
		if err != nil {
			tools.ZapLogger("file", "agent").Error(err.Error())
			return
		}

		// Do action
		tools.ZapLogger("file", "agent").Info(fmt.Sprint(payload))
		switch payload.Type {
		case "service":
			services.ServiceAction(payload.Service.ServiceName, payload.Service.ServiceName)
		}
	}()
}

// Connection handler
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	tools.ZapLogger("both", "agent").Info("Connected to Meta Handler MQTT Server")
}

// Lost connection handler
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	message := fmt.Sprintf("MQTT connection lost. Error: %v", err)
	tools.ZapLogger("both", "agent").Info(message)
}
