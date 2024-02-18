package agent

import "time"

type Agent struct {
	ID           string    `json:"id"`
	Address      string    `json:"address"`
	Status       string    `json:"status_code"`
	LastHearBeat time.Time `json:"last_heartbeat"`
}
