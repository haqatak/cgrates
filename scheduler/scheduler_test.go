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
package scheduler

import (
	"testing"
	"time"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func TestSchedulerUpdateActStats(t *testing.T) {
	sched := &Scheduler{actStatsInterval: time.Millisecond, actSuccessStats: make(map[string]map[time.Time]bool)}
	sched.updateActStats(&engine.Action{Id: "REMOVE_1", ActionType: utils.MetaRemoveAccount}, false)
	if len(sched.actSuccessStats[utils.MetaRemoveAccount]) != 1 {
		t.Errorf("Wrong stats: %+v", sched.actSuccessStats[utils.MetaRemoveAccount])
	}
	sched.updateActStats(&engine.Action{Id: "REMOVE_2", ActionType: utils.MetaRemoveAccount}, false)
	if len(sched.actSuccessStats[utils.MetaRemoveAccount]) != 2 {
		t.Errorf("Wrong stats: %+v", sched.actSuccessStats[utils.MetaRemoveAccount])
	}
	sched.updateActStats(&engine.Action{Id: "LOG1", ActionType: utils.MetaLog}, false)
	if len(sched.actSuccessStats[utils.MetaLog]) != 1 ||
		len(sched.actSuccessStats[utils.MetaRemoveAccount]) != 2 {
		t.Errorf("Wrong stats: %+v", sched.actSuccessStats)
	}
	time.Sleep(sched.actStatsInterval)
	sched.updateActStats(&engine.Action{Id: "REMOVE_3", ActionType: utils.MetaRemoveAccount}, false)
	if len(sched.actSuccessStats[utils.MetaRemoveAccount]) != 1 || len(sched.actSuccessStats) != 1 {
		t.Errorf("Wrong stats: %+v", sched.actSuccessStats)
	}
}

func TestNewScheduler(t *testing.T) {
	// Create mock dependencies
	cfg := config.NewDefaultCGRConfig()

	// Create DataDB mock, engine needs DataDB
	dataDBMock := &engine.DataDBMock{}
	dm := engine.NewDataManager(dataDBMock, cfg.CacheCfg(), nil)

	fltrS := engine.NewFilterS(cfg, nil, dm)

	// Call NewScheduler
	s := NewScheduler(dm, cfg, fltrS)

	// Verify it returns a valid non-nil instance
	if s == nil {
		t.Fatal("NewScheduler returned nil")
	}

	// Verify fields are correctly assigned
	if s.dm != dm {
		t.Errorf("expected dm to be assigned to scheduler, got %v", s.dm)
	}
	if s.cfg != cfg {
		t.Errorf("expected cfg to be assigned to scheduler, got %v", s.cfg)
	}
	if s.fltrS != fltrS {
		t.Errorf("expected fltrS to be assigned to scheduler, got %v", s.fltrS)
	}
	if s.restartLoop == nil {
		t.Error("expected restartLoop channel to be initialized, but it was nil")
	}
}
