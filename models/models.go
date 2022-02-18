package models

// Event ...
type Event struct {
	Address      string `json:"address"`
	OwnerAddress string `json:"ownerAddress"`
	ID           uint64 `json:"id"`
}
