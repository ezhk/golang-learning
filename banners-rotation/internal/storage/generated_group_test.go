// generated by generate-table-tests; DO NOT EDIT
// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/suite"
)

type GroupTestSuite struct {
	suite.Suite
	db *Storage
}

func TestGroupSuite(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}

func (s *GroupTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	s.NoError(err)

	// Define storage.
	s.db = db

	// Clean previous values.
	s.TearDownTest()
}

func (s *GroupTestSuite) TearDownTest() {
	// Clean exists groups.
	groups, err := s.db.ReadGroups()
	s.NoError(err)
	for _, group := range groups {
		err = s.db.DeleteGroup(group.ID)
		s.NoError(err)
	}
}

func (s *GroupTestSuite) TestGroupOperations() {
	// Create new group.
	group, err := s.db.CreateGroup("test group", "test description")
	s.NoError(err)
	s.Equal("test group", group.Name)

	group.Name = "updated test group"
	updatedGroup, err := s.db.UpdateGroup(group)
	s.NoError(err)
	s.Equal("updated test group", updatedGroup.Name)

	groups, err := s.db.ReadGroups()
	s.NoError(err)
	s.Greater(len(groups), 0)
	s.Equal("updated test group", groups[0].Name)

	// Call "duplicate key value violates unique constraint".
	_, err = s.db.CreateGroup("updated test group", "empty")
	s.Error(err)

	err = s.db.DeleteGroup(groups[0].ID)
	s.NoError(err)
}
