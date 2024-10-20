package setup

import (
	"github.com/gar-id/Meta/apps/agent/windows/services"
)

func Install() {
	services.ServiceAction("install", "MetaHandlerAgent")
}
