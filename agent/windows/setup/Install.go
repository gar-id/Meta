package setup

import (
	"MetaHandler/agent/windows/services"
)

func Install() {
	services.ServiceAction("install", "MetaHandlerAgent")
}
