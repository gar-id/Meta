package services

import (
	"fmt"

	"golang.org/x/sys/windows/svc/mgr"
)

func statusService(name string) error {
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("cannot connect to manager %v", err)
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s does not exist: %v", name, err)
	}
	defer service.Close()
	status, err := service.Query()
	if err != nil {
		return fmt.Errorf("failed to get service status: %v", err)
	}

	switch status.State {
	case 1:
		err = fmt.Errorf("stopped")
	case 2:
		err = fmt.Errorf("start pending")
	case 3:
		err = fmt.Errorf("stop pending")
	case 4:
		err = fmt.Errorf("running")
	case 5:
		err = fmt.Errorf("continue pending")
	case 6:
		err = fmt.Errorf("pause pending")
	case 7:
		err = fmt.Errorf("paused")
	}
	return err
}
