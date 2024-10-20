package services

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func stopService(name string) error {
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

	_, err = service.Control(svc.Cmd(svc.Stop))
	if err != nil {
		return fmt.Errorf("could not stop the service: %v", err)
	}
	return nil
}
