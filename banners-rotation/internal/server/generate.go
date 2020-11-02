package server

// Generate object methods.
//go:generate go run ./generate-server-methods/... -object banner -file generated_banner.go
//go:generate go run ./generate-server-methods/... -object slot -file generated_slot.go
//go:generate go run ./generate-server-methods/... -object group -file generated_group.go
