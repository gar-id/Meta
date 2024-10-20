package mqttMochi

import (
	"fmt"
	"log/slog"
	"os"

	"MetaHandler/server/config/caches"

	"MetaHandler/tools"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func Start() error {
	// Create the new MQTT Server.
	server := mqtt.New(nil)

	// Setup Mochi MQTT logging
	level := new(slog.LevelVar)
	filelocation := tools.DefaultString(caches.MetaHandlerServer.MetaHandlerServer.Log.Location, "/var/log/meta-handler")
	filelocation = fmt.Sprintf("%v/Meta-Handler-Server-%v.log", filelocation, "mqtt")
	logfile, err := os.OpenFile(filelocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		tools.ZapLogger("file", "server").Error(fmt.Sprintf("error opening file: %v", err))
	}
	server.Log = slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	switch caches.MetaHandlerServer.MetaHandlerServer.Log.Level {
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warning":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	case "panic":
		level.Set(slog.LevelError)
	case "fatal":
		level.Set(slog.LevelError)
	default:
		level.Set(slog.LevelDebug)
	}
	defer logfile.Sync()

	// Setup connection ACL
	_ = server.AddHook(new(auth.Hook), &auth.Options{
		Ledger: &auth.Ledger{
			Auth: auth.AuthRules{ // Auth disallows all by default
				{
					Username: auth.RString(tools.DefaultString(
						caches.MetaHandlerServer.MetaHandlerServer.MQTT.AdminUser,
						"admin")),
					Password: auth.RString(tools.DefaultString(
						caches.MetaHandlerServer.MetaHandlerServer.MQTT.AdminPass,
						"admin")),
					Allow: true,
				}, {
					Username: auth.RString(tools.DefaultString(
						caches.MetaHandlerServer.MetaHandlerServer.MQTT.AgentUser,
						"agent")),
					Password: auth.RString(tools.DefaultString(
						caches.MetaHandlerServer.MetaHandlerServer.MQTT.AgentPass,
						"agent")),
					Allow: true,
				},
				{
					Remote: "127.0.0.1:*",
					Allow:  true,
				},
				{
					Remote: "localhost:*",
					Allow:  true,
				},
			},
			ACL: auth.ACLRules{ // ACL allows all by default
				{
					Remote: "127.0.0.1:*",
				}, // local superuser allow all
				{
					// user admin can read and write to their own topic
					Username: auth.RString(tools.DefaultString(
						caches.MetaHandlerServer.MetaHandlerServer.MQTT.AdminUser,
						"admin")),
					Filters: auth.Filters{
						"#": auth.ReadWrite,
					},
				},
				{
					// user agent can read and write to their own topic
					Username: auth.RString(tools.DefaultString(
						caches.MetaHandlerServer.MetaHandlerServer.MQTT.AgentUser,
						"agent")),
					Filters: auth.Filters{
						"agent/#":   auth.ReadWrite,
						"updates/#": auth.WriteOnly, // can write to updates, but can't read updates from others
					},
				},
				{
					// Otherwise, no clients have publishing or subscribe permissions
					Filters: auth.Filters{
						"#":         auth.Deny,
						"updates/#": auth.Deny,
					},
				},
			},
		},
	})

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP(listeners.Config{
		ID: "MetaHandler_MQTT",
		Address: fmt.Sprintf("%v:%v",
			caches.MetaHandlerServer.MetaHandlerServer.MQTT.Host,
			caches.MetaHandlerServer.MetaHandlerServer.MQTT.Port)})
	err = server.AddListener(tcp)
	if err != nil {
		tools.ZapLogger("both", "server").Fatal(fmt.Sprintf("Failed to start MetaHandler MQTT Server. %v", err))
		return err
	}

	// Start MQTT Server
	err = server.Serve()
	if err != nil {
		tools.ZapLogger("both", "server").Fatal(fmt.Sprintf("Failed to start MetaHandler MQTT Server. %v", err))
		return err
	}
	return err
}
