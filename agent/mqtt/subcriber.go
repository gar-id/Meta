package mqttAgent

import (
	"fmt"
	"time"

	"MetaHandler/tools"
)

func Subscriber(topic string) {
	// Create connection
	client := connection()

	message := fmt.Sprintf("Subscribed to topic: %s", topic)
	tools.ZapLogger("file", "agent").Info(message)
	for client.IsConnected() {
		// Start to subscribe
		token := client.Subscribe(topic, 1, nil)
		token.Wait()
		time.Sleep(time.Second)
	}
}
