package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("generate table tests", func(t *testing.T) {
		data, err := generate("slot")
		require.NoError(t, err)

		require.Equal(t, `// generated by generate-table-tests; DO NOT EDIT
// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/suite"
)

type SlotTestSuite struct {
	suite.Suite
	db DatabaseInterface
}

func TestSlotSuite(t *testing.T) {
	suite.Run(t, new(SlotTestSuite))
}

func (s *SlotTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	s.NoError(err)

	// Define storage.
	s.db = db
}

func (s *SlotTestSuite) TestSlotOperations() {
	// Create new slot.
	slot, err := s.db.CreateSlot("test slot", "test description")
	s.NoError(err)
	defer s.db.DeleteSlot(slot.ID)

	s.Equal("test slot", slot.Name)

	slot.Name = "updated test slot"
	updatedSlot, err := s.db.UpdateSlot(slot)
	s.NoError(err)
	s.Equal("updated test slot", updatedSlot.Name)

	slots, err := s.db.ReadSlots()
	s.NoError(err)
	s.Greater(len(slots), 0)
	for _, obj := range slots {
		if obj.ID != slot.ID {
			continue
		}

		s.Equal("updated test slot", obj.Name)
	}

	// Call "duplicate key value violates unique constraint".
	_, err = s.db.CreateSlot("updated test slot", "empty")
	s.Error(err)

	err = s.db.DeleteSlot(slot.ID)
	s.NoError(err)
}
`, string(data))
	})
}
