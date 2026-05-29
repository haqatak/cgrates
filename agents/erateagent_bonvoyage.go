package agents

import (
)

// BonVoyageSMSRequest represents the request body for Bon Voyage SMS
type BonVoyageSMSRequest struct {
	NewCountryCode string `json:"new_country_code"`
}

// FraudAlertRequest represents the request body for Fraud Alert
type FraudAlertRequest struct {
	FraudType string `json:"fraud_type"`
}

// DataCostControlRequest represents the request body for Data Cost Control
type DataCostControlRequest struct {
	DCCLimit int `json:"dcc_limit"`
}
