// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/suite"
)

type SlotsTestSuite struct {
	suite.Suite
	db *Storage
}

func TestSlotsSuite(t *testing.T) {
	suite.Run(t, new(SlotsTestSuite))
}

func (s *SlotsTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	s.NoError(err)

	// Define Storage.
	s.db = db

	s.TearDownTest()
}

func (s *SlotsTestSuite) TearDownTest() {
	// Clean exists slots.
	slots, err := s.db.ReadSlots()
	s.NoError(err)
	for _, slot := range slots {
		err = s.db.DeleteSlot(slot.ID)
		s.NoError(err)
	}
}

func (s *SlotsTestSuite) TestSlotOperations() {
	// Create new slot.
	slot, err := s.db.CreateSlot("test slot")
	s.NoError(err)
	s.Equal("test slot", slot.Name)

	slot.Name = "updated test slot"
	updatedSlot, err := s.db.UpdateSlot(slot)
	s.NoError(err)
	s.Equal("updated test slot", updatedSlot.Name)

	slots, err := s.db.ReadSlots()
	s.NoError(err)
	s.Greater(len(slots), 0)
	s.Equal("updated test slot", slots[0].Name)

	// Call "duplicate key value violates unique constraint".
	_, err = s.db.CreateSlot("updated test slot")
	s.Error(err)

	err = s.db.DeleteSlot(slots[0].ID)
	s.NoError(err)
}
