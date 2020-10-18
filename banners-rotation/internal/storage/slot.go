package storage

import "github.com/ezhk/golang-learning/banners-rotation/internal/structs"

func (s *Storage) CreateSlot(slotName string) (structs.Slot, error) {
	slot := structs.Slot{Name: slotName}
	result := s.db.Create(&slot)

	return slot, result.Error
}

func (s *Storage) ReadSlots() ([]*structs.Slot, error) {
	slots := []*structs.Slot{}
	result := s.db.Find(&slots)

	return slots, result.Error
}

func (s *Storage) UpdateSlot(sl structs.Slot) (structs.Slot, error) {
	slot := structs.Slot{ID: sl.ID}
	result := s.db.Model(&slot).Updates(sl)

	return slot, result.Error
}

func (s *Storage) DeleteSlot(id uint) error {
	result := s.db.Delete(&structs.Slot{}, id)

	return result.Error
}
