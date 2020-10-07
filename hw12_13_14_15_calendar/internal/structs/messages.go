package structs

// ErrorMessage - JSON errors message view.
type ErrorMessage struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UserMessage - JSON user message view.
// gRPC gateway; proto3 -> JSON mapping: int64, fixed64, uint64 -> string
// https://developers.google.com/protocol-buffers/docs/proto3#json
type UserMessage struct {
	ID        string `json:"ID"`
	Email     string `json:"Email"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

// EventMessage - JSON event message.
type EventMessage struct {
	ID       string `json:"ID"`
	UserID   string `json:"UserID"`
	Title    string `json:"Title"`
	Content  string `json:"Content"`
	DateFrom string `json:"DateFrom"`
	DateTo   string `json:"DateTo"`
	Notified bool   `json:"Notified"`
}

// ManyEventsMessage - multiple EventMessage struct.
type ManyEventsMessage struct {
	Events []EventMessage `json:"Events"`
}
