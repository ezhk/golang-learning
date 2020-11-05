package structs

type QueueEvent struct {
	PlacementID uint64 `json:"placement_id"`
	EventType   string `json:"event_type"`
}
