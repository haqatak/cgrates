package config

// ERateAgentCfg is the ERateAgent config section
type ERateAgentCfg struct {
	Enabled bool
}

// Clone returns a deep copy of ERateAgentCfg
func (e *ERateAgentCfg) Clone() *ERateAgentCfg {
	if e == nil {
		return nil
	}
	return &ERateAgentCfg{
		Enabled: e.Enabled,
	}
}
