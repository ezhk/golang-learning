// generated by generate-table-methods; DO NOT EDIT
package storage

import "github.com/ezhk/golang-learning/banners-rotation/internal/structs"

func (s *Storage) CreateGroup(groupName string) (structs.Group, error) {
	group := structs.Group{Name: groupName}
	result := s.db.Create(&group)

	return group, result.Error
}

func (s *Storage) ReadGroups() ([]*structs.Group, error) {
	groups := []*structs.Group{}
	result := s.db.Find(&groups)

	return groups, result.Error
}

func (s *Storage) UpdateGroup(b structs.Group) (structs.Group, error) {
	group := structs.Group{ID: b.ID}
	result := s.db.Model(&group).Updates(b)

	return group, result.Error
}

func (s *Storage) DeleteGroup(id uint) error {
	result := s.db.Delete(&structs.Group{}, id)

	return result.Error
}
