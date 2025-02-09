package domain

import "time"

type PingData struct {
	IPAddress          string    `json:"ip_address"`
	PingTime           int       `json:"ping_time"`
	LastSuccessfulPing time.Time `json:"last_successful_ping"`
}
