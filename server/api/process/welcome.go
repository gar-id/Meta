package process

import (
	"fmt"
	"time"

	"MetaHandler/core/caches"
	APItypes "MetaHandler/server/api/process/types"
)

func WelcomeGeneral(httpCode int, statusText, message string) APItypes.Welcome {
	var result = APItypes.Welcome{
		HTTP_Code: httpCode,
		Status:    statusText,
		Data: struct {
			Date    string "json:\"date\""
			Message string "json:\"message\""
			Version string "json:\"version\""
		}{
			Date:    fmt.Sprint(time.Now().String()),
			Message: message,
			Version: caches.MetaHandlerVersion,
		}}

	return result
}
