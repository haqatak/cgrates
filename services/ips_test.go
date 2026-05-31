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
	"reflect"
	"sync"
	"testing"

	"github.com/cgrates/birpc"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

// TestIPSCoverage for cover testing
func TestIPSCoverage(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	cfg.IPsCfg().Enabled = true
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	shdChan := utils.NewSyncedChan()
	chS := engine.NewCacheS(cfg, nil, nil)
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	db := NewDataDBService(cfg, nil, false, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan birpc.ClientConnector, 1), srvDep)
	ipS := NewIPService(cfg, db, chS, filterSChan, server, make(chan birpc.ClientConnector, 1), nil, anz, srvDep)

	if ipS.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	ipS2 := &IPService{
		cfg:         cfg,
		dbs:         db,
		cache:       chS,
		fsChan:      filterSChan,
		server:      server,
		cm:          nil,
		ips:         engine.NewIPService(nil, cfg, nil, nil),
		connChan:    make(chan birpc.ClientConnector, 1),
		anz:         anz,
		srvDep:      srvDep,
	}
	if !ipS2.IsRunning() {
		t.Errorf("Expected service to be running")
	}
	serviceName := ipS2.ServiceName()
	if !reflect.DeepEqual(serviceName, utils.IPs) {
		t.Errorf("\nExpecting <%+v>,\n Received <%+v>", utils.IPs, serviceName)
	}
	shouldRun := ipS2.ShouldRun()
	if !reflect.DeepEqual(shouldRun, true) {
		t.Errorf("\nExpecting <true>,\n Received <%+v>", shouldRun)
	}
	cacheSrv, err := engine.NewService(chS)
	if err != nil {
		t.Fatal(err)
	}
	ipS2.srvDep[utils.DataDB].Add(1)
	ipS2.connChan <- cacheSrv
	ipS2.Shutdown()
	if ipS2.IsRunning() {
		t.Errorf("Expected service to be down")
	}
}

func TestIPSReload(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	ipS := &IPService{
		ips: engine.NewIPService(nil, cfg, nil, nil),
	}
	ipS.ips.StartLoop()
	err := ipS.Reload()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
