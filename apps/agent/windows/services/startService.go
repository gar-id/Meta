package services

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func startService(name string) error {
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

	err = service.Start()
	if err != nil {
		return fmt.Errorf("could not start the service: %v", err)
	}

	var status svc.Status
	for wait := 0; wait <= 5; wait++ {
		status, err = service.Query()
		if err != nil {
			return fmt.Errorf("failed to get service status: %v", err)
		}

		if status.State == svc.Running {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("start did not complete, final state: %v", status.State)
}
