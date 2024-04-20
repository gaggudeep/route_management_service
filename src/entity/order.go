package entity

type Order struct {
	Restaurant Restaurant `json:"restaurant"`
	Consumer   Consumer   `json:"consumer"`
}
