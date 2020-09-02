package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	t.Run("user operations", func(t *testing.T) {
		db := NewDatatabase()
		id, err := db.CreateUser("Vasya", "Pupkin")

		require.Nil(t, err)
		require.Equal(t, 1, id)
		// require.Equal(t, 1, len(db.Users))

		err = db.UpdateUser(id, "Винни", "Пух")
		require.Nil(t, err)
		// require.Equal(t, 1, len(db.Users))

		u := db.GetUser("Винни", "Пух")
		require.Equal(t, id, u.ID)
		require.Equal(t, "Винни", u.FirstName)
		require.Equal(t, "Пух", u.LastName)

		_, err = db.CreateUser("Винни", "Пух")
		require.NotNil(t, err)

		err = db.DeleteUser(id)
		require.Nil(t, err)

		// secord run cause error "user not exist"
		err = db.DeleteUser(id)
		require.NotNil(t, err)
	})

	t.Run("calendar operations", func(t *testing.T) {
		db := NewDatatabase()

		// calendar operations must contain user
		userID, err := db.CreateUser("Винни", "Пух")
		require.Nil(t, err)

		id, err := db.CreateRecord(userID, "Встреча", "Кофе в кафе")
		require.Nil(t, err)
		require.Equal(t, 1, id)

		rec := db.GetRecords(userID)
		require.Equal(t, 1, len(rec))
		require.Equal(t, "Встреча", rec[0].Title)
		require.Equal(t, "Кофе в кафе", rec[0].Content)
		require.True(t, time.Now().After(rec[0].UpdatedAt))

		err = db.UpdateRecord(id, "Встреча", "Coffee time")
		require.Nil(t, err)

		rec = db.GetRecords(userID)
		require.Equal(t, 1, len(rec))
		require.Equal(t, "Coffee time", rec[0].Content)

		err = db.DeleteRecord(rec[0].ID)
		require.Nil(t, err)

		rec = db.GetRecords(userID)
		require.Equal(t, 0, len(rec))
	})
}
