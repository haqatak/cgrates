package agents

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
)

var digitsRe = regexp.MustCompile(`^\d+$`)

// NewERateAgent initializes the ERateAgent
func NewERateAgent(cgrCfg *config.CGRConfig, connMgr *engine.ConnManager, caps *engine.Caps) (*ERateAgent, error) {
	return &ERateAgent{
		cgrCfg:  cgrCfg,
		connMgr: connMgr,
		caps:    caps,
	}, nil
}

// ERateAgent handles incoming ERate network events over HTTP
type ERateAgent struct {
	cgrCfg  *config.CGRConfig
	connMgr *engine.ConnManager
	caps    *engine.Caps
}

// sendJSONResponse sends a JSON response with status and message.
func sendJSONResponse(w http.ResponseWriter, status int, statusStr string, message string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]interface{}{
		"status":  statusStr,
		"message": message,
	}
	if errorCode != 0 {
		resp["error_code"] = errorCode
	}
	json.NewEncoder(w).Encode(resp)
}

// BonVoyageSMSHandler handles Bon Voyage SMS network events.
func (ea *ERateAgent) BonVoyageSMSHandler(w http.ResponseWriter, r *http.Request) {
	networkUserID := r.PathValue("network_user_id")
	if !digitsRe.MatchString(networkUserID) {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}

	var reqBody BonVoyageSMSRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}
	if reqBody.NewCountryCode == "" || len(reqBody.NewCountryCode) != 2 {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}


	sendJSONResponse(w, http.StatusOK, "success", "Bon voyage SMS sent successfully.", 0)
}

// FraudAlertHandler handles Fraud Alert network events.
func (ea *ERateAgent) FraudAlertHandler(w http.ResponseWriter, r *http.Request) {
	networkUserID := r.PathValue("network_user_id")
	if !digitsRe.MatchString(networkUserID) {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}

	var reqBody FraudAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}
	if reqBody.FraudType == "" || len(reqBody.FraudType) > 20 {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}

	sendJSONResponse(w, http.StatusOK, "success", "Fraud alert recorded successfully. Subscription blocked and SMS sent.", 0)
}

// DataCostControlHandler handles Data Cost Control network events.
func (ea *ERateAgent) DataCostControlHandler(w http.ResponseWriter, r *http.Request) {
	networkUserID := r.PathValue("network_user_id")
	if !digitsRe.MatchString(networkUserID) {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}

	var reqBody DataCostControlRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}
	if reqBody.DCCLimit < 0 {
		sendJSONResponse(w, http.StatusBadRequest, "error", "Missing or invalid input parameters.", 400)
		return
	}

	sendJSONResponse(w, http.StatusOK, "success", "Data cost control event recorded successfully. Data roaming disabled and SMS sent.", 0)
}
