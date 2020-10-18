// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/suite"
)

type GroupsTestSuite struct {
	suite.Suite
	db *Storage
}

func TestGroupsSuite(t *testing.T) {
	suite.Run(t, new(GroupsTestSuite))
}

func (s *GroupsTestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	s.NoError(err)

	// Define Storage.
	s.db = db

	// Clean exists Groups.
	groups, err := s.db.ReadGroups()
	s.NoError(err)
	for _, group := range groups {
		err = s.db.DeleteGroup(group.ID)
		s.NoError(err)
	}
}

func (s *GroupsTestSuite) TestGroupOperations() {
	// Create new Group.
	group, err := s.db.CreateGroup("test group")
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
	_, err = s.db.CreateGroup("updated test group")
	s.Error(err)

	err = s.db.DeleteGroup(groups[0].ID)
	s.NoError(err)
}
