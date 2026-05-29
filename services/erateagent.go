package services

import (
	"fmt"
	"sync"

	"github.com/cgrates/cgrates/agents"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/servmanager"
	"github.com/cgrates/cgrates/utils"
)

// NewERateAgent returns the ERate Agent
func NewERateAgent(cfg *config.CGRConfig,
	server *cores.Server, connMgr *engine.ConnManager, caps *engine.Caps,
	srvDep map[string]*sync.WaitGroup) servmanager.Service {
	return &ERateAgent{
		cfg:         cfg,
		server:      server,
		connMgr:     connMgr,
		caps:        caps,
		srvDep:      srvDep,
	}
}

// ERateAgent implements Service interface
type ERateAgent struct {
	sync.RWMutex
	cfg         *config.CGRConfig
	server      *cores.Server
	eA          *agents.ERateAgent

	started bool
	connMgr *engine.ConnManager
	caps    *engine.Caps
	srvDep  map[string]*sync.WaitGroup
}

// Start should handle the service start
func (ea *ERateAgent) Start() (err error) {
	ea.Lock()
	if ea.started {
		ea.Unlock()
		return utils.ErrServiceAlreadyRunning
	}
	ea.eA, err = agents.NewERateAgent(ea.cfg, ea.connMgr, ea.caps)
	if err != nil {
		return
	}

	// Register the endpoints using go 1.22+ ServeMux routing capabilities
	// that cores.Server.RegisterHttpFunc delegates to
	ea.server.RegisterHttpFunc("POST /api/network-events/v1/users/{network_user_id}/bon-voyage-sms", ea.eA.BonVoyageSMSHandler)
	ea.server.RegisterHttpFunc("POST /api/network-events/v1/users/{network_user_id}/fraud-alert", ea.eA.FraudAlertHandler)
	ea.server.RegisterHttpFunc("POST /api/network-events/v1/users/{network_user_id}/data-cost-control", ea.eA.DataCostControlHandler)

	ea.started = true
	ea.Unlock()
	utils.Logger.Info(fmt.Sprintf("<%s> successfully started.", "ERateAgent"))
	return
}

// Reload handles the change of config
func (ea *ERateAgent) Reload() (err error) {
	return // no reload
}

// Shutdown stops the service
func (ea *ERateAgent) Shutdown() (err error) {
	ea.Lock()
	ea.started = false
	ea.Unlock()
	return
}

// IsRunning returns if the service is running
func (ea *ERateAgent) IsRunning() bool {
	ea.RLock()
	defer ea.RUnlock()
	return ea.started
}

// ServiceName returns the service name
func (ea *ERateAgent) ServiceName() string {
	return "ERateAgent"
}

// ShouldRun returns if the service should be running
func (ea *ERateAgent) ShouldRun() bool {
	return ea.cfg.ERateAgentCfg().Enabled
}
