package services

import (
	"fmt"
	"path"

	"golang.org/x/sys/windows/svc/mgr"
)

func installService(name string) error {
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("cannot connect to manager %v", err)
	}
	defer manager.Disconnect()

	serviceCfg := mgr.Config{
		// ServiceType:      16,
		StartType:        mgr.StartAutomatic,
		ErrorControl:     mgr.ErrorCritical,
		BinaryPathName:   path.Join("C", "Program Files", "Meta Handler", "agent.exe"),
		DisplayName:      "Meta Handler Agent",
		Description:      "Created by Tegar, this agent used for failover stunnel, meta service action, and more",
		ServiceStartName: name,
	}

	service, err := manager.CreateService("meta-handler-agent", serviceCfg.BinaryPathName, serviceCfg)
	if err != nil {
		return fmt.Errorf("failed to install service %s: %v", name, err)
	}
	defer service.Close()

	return err
}
