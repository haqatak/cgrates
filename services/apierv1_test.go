/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*/

package services

import (
	"sync"
	"testing"

	"github.com/cgrates/birpc"
	v1 "github.com/cgrates/cgrates/apier/v1"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func TestNewAPIerSv1Service(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	filterSChan := make(chan *engine.FilterS, 1)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	server := cores.NewServer(nil)
	db := NewDataDBService(cfg, nil, false, srvDep)
	stordb := NewStorDBService(cfg, false, srvDep)

	internalChan := make(chan birpc.ClientConnector, 1)
	shdChan := utils.NewSyncedChan()
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, internalChan, srvDep)

	schS := NewSchedulerService(cfg, db, nil, filterSChan, server, internalChan, nil, anz, srvDep)

	apiSv1 := NewAPIerSv1Service(cfg, db, stordb, filterSChan, server, schS, new(ResponderService), internalChan, nil, anz, srvDep)

	if apiSv1 == nil {
		t.Errorf("Expected APIerSv1Service to be created, got nil")
	}
}

func TestAPIerSv1Service_Reload(t *testing.T) {
	apiService := &APIerSv1Service{}
	err := apiService.Reload()
	if err != nil {
		t.Errorf("Expected Reload to return no error, got %v", err)
	}
}

func TestAPIerSv1Service_IsRunning(t *testing.T) {
	apiService := &APIerSv1Service{}
	if apiService.IsRunning() {
		t.Errorf("Expected service to not be running")
	}

	apiService.api = &v1.APIerSv1{}
	if !apiService.IsRunning() {
		t.Errorf("Expected service to be running")
	}
}

func TestAPIerSv1Service_ServiceName(t *testing.T) {
	apiService := &APIerSv1Service{}
	if apiService.ServiceName() != utils.APIerSv1 {
		t.Errorf("Expected service name %s, got %s", utils.APIerSv1, apiService.ServiceName())
	}
}

func TestAPIerSv1Service_GetAPIerSv1(t *testing.T) {
	apiService := &APIerSv1Service{}
	api := &v1.APIerSv1{}
	apiService.api = api
	if apiService.GetAPIerSv1() != api {
		t.Errorf("Expected GetAPIerSv1 to return correct APIerSv1")
	}
}

func TestAPIerSv1Service_ShouldRun(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	apiService := &APIerSv1Service{cfg: cfg}

	cfg.ApierCfg().Enabled = false
	if apiService.ShouldRun() {
		t.Errorf("Expected ShouldRun to be false")
	}

	cfg.ApierCfg().Enabled = true
	if !apiService.ShouldRun() {
		t.Errorf("Expected ShouldRun to be true")
	}
}

func TestAPIerSv1Service_GetAPIerSv1Chan(t *testing.T) {
	ch := make(chan *v1.APIerSv1, 1)
	apiService := &APIerSv1Service{APIerSv1Chan: ch}
	if apiService.GetAPIerSv1Chan() != ch {
		t.Errorf("Expected GetAPIerSv1Chan to return correct channel")
	}
}

func TestAPIerSv1Service_Start_AlreadyRunning(t *testing.T) {
	apiService := &APIerSv1Service{
		api: &v1.APIerSv1{},
	}
	err := apiService.Start()
	if err != utils.ErrServiceAlreadyRunning {
		t.Errorf("Expected Start to return ErrServiceAlreadyRunning, got %v", err)
	}
}
