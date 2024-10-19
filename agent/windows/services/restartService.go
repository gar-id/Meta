package services

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func restartService(name string) error {
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

	if status.State != svc.Stopped {
		_, err = service.Control(svc.Cmd(svc.Stop))
		if err != nil {
			return fmt.Errorf("could not stop the service: %v", err)
		}
	}

	for wait := 0; wait <= 15; wait++ {
		status, err = service.Query()
		if err != nil {
			return fmt.Errorf("failed to get service status: %v", err)
		}

		if status.State == svc.Stopped {
			break
		}
		time.Sleep(100 * time.Microsecond)
	}

	err = startService(name)
	return err
}
