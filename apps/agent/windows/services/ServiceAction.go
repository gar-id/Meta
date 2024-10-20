package services

import (
	"fmt"
	"strings"
)

func ServiceAction(action, serviceName string) {
	action = strings.ToLower(action)
	var err error
	switch action {
	case "start":
		err = startService(serviceName)
	case "stop":
		err = stopService(serviceName)
	case "restart":
		err = restartService(serviceName)
	case "install":
		err = installService(serviceName)
	case "status":
		err = statusService(serviceName)
		fmt.Printf("%v for service %v is %v", action, serviceName, err)
		return
	default:
		fmt.Printf("Action %v is not recognized", action)
		return
	}

	if err != nil {
		fmt.Print(err)
	} else if action != "status" {
		fmt.Printf("%v for service %v is success", action, serviceName)
	}
}
