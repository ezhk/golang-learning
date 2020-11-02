package storage

// Generate base methods for tables.
//go:generate go run ./generate-table-methods/... -table banner -file generated_banner.go
//go:generate go run ./generate-table-methods/... -table slot -file generated_slot.go
//go:generate go run ./generate-table-methods/... -table group -file generated_group.go

// Generate base tests for table methods.
//go:generate go run ./generate-table-tests/... -table banner -file generated_banner_test.go
//go:generate go run ./generate-table-tests/... -table slot -file generated_slot_test.go
//go:generate go run ./generate-table-tests/... -table group -file generated_group_test.go
