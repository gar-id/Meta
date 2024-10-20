package mqttCSSH

import (
	"encoding/json"
	"fmt"

	"github.com/gar-id/centralissh/agent/jobs"
	"github.com/gar-id/centralissh/agent/mqtt/types"
	"github.com/gar-id/centralissh/tools"

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
			tools.ZapLogger("file").Error(err.Error())
			return
		}

		// Do action
		tools.ZapLogger("file").Info(fmt.Sprint(payload))
		switch payload.Opt {
		case "sudoers":
			jobs.Sudoers(payload.Action, payload.Data.Sudoers)
		case "keypair":
			jobs.UserKeypair(payload.Action, payload.Data.Keypair)
		}
	}()
}

// Connection handler
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	tools.ZapLogger("both").Info("Connected to CentraliSSH MQTT Server")
}

// Lost connection handler
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	message := fmt.Sprintf("MQTT connection lost. Error: %v", err)
	tools.ZapLogger("both").Info(message)
}
