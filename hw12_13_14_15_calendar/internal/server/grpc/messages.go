package internalgrpc

type ErrorMessage struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// gRPC gateway; proto3 -> JSON mapping: int64, fixed64, uint64 -> string
// https://developers.google.com/protocol-buffers/docs/proto3#json
type UserMessage struct {
	ID        string `db:"id" json:"ID"`
	Email     string `db:"email" json:"Email"`
	FirstName string `db:"first_name" json:"FirstName"`
	LastName  string `db:"last_name" json:"LastName"`
}

type EventMessage struct {
	ID       string `db:"id" json:"ID"`
	UserID   string `db:"user_id" json:"UserID"`
	Title    string `db:"title" json:"Title"`
	Content  string `db:"content" json:"Content"`
	DateFrom string `db:"date_from" json:"DateFrom"`
	DateTo   string `db:"date_to" json:"DateTo"`
	Notified bool   `db:"notified" json:"Notified"`
}

type ManyEventsMessage struct {
	Events []EventMessage `json:"Events"`
}
