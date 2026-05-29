1. **Config additions:** Add a new `ERateAgentCfg` struct in `config` to control `enabled` flag and base `url` for the agent. Register the agent inside `config_defaults.go` and `config_json.go`.
2. **ERate Agent Service (`services/erateagent.go`):** Implement `servmanager.Service` similar to `JanusAgent` which wraps the `agents.ERateAgent`. It registers the specific routes using `ja.server.RegisterHttpFunc`.
3. **ERate Agent logic (`agents/erateagent.go`):** Implement `agents.ERateAgent` with handler functions for the three events:
  - `POST /api/network-events/v1/users/{network_user_id}/bon-voyage-sms`
  - `POST /api/network-events/v1/users/{network_user_id}/fraud-alert`
  - `POST /api/network-events/v1/users/{network_user_id}/data-cost-control`
These handlers will validate input (e.g., regex `^\d+$` for `network_user_id`), check payload structures, and convert the event to an internal CGRateS event (using `engine.ConnManager` or just returning `200 OK` + JSON as a stub for now if no internal mapping is specified).
4. **Registration:** Update `services/engine.go` to construct and start `NewERateAgent`.
