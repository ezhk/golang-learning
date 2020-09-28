// +build integration

package sqlstorage

import (
	"testing"
	"time"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	t.Run("user operations", func(t *testing.T) {
		db := NewDatatabase()

		err := db.Connect("user=postgres password=postgres dbname=postgres sslmode=disable")
		require.NoError(t, err)
		defer db.Close()

		id, err := db.CreateUser("vasya@pupkin.com", "Vasya", "Pupkin")
		require.NoError(t, err)

		err = db.UpdateUser(id, "винни@пух.рф", "Винни", "Пух")
		require.NoError(t, err)

		u, err := db.GetUserByEmail("винни@пух.рф")
		require.NoError(t, err)

		require.Equal(t, id, u.ID)
		require.Equal(t, "Винни", u.FirstName)
		require.Equal(t, "Пух", u.LastName)

		_, err = db.CreateUser("винни@пух.рф", "Винни", "Пух")
		require.NotNil(t, err)
		require.Equal(t, err, storage.ErrUserExists)

		err = db.DeleteUser(id)
		require.NoError(t, err)

		// secord run not cause error "user not exist" in SQL
		err = db.DeleteUser(id)
		require.Nil(t, err)
	})

	t.Run("calendar operations", func(t *testing.T) {
		db := NewDatatabase()

		err := db.Connect("user=postgres password=postgres dbname=postgres sslmode=disable")
		require.NoError(t, err)
		defer db.Close()

		// calendar operations must contain user
		userID, err := db.CreateUser("vinny@pooh.com", "Винни", "Пух")
		require.Nil(t, err)
		defer db.DeleteUser(userID)

		// userID int64, title, content string, dateFrom, dateTo time.Time
		recordStartDate := time.Date(2020, time.September, 1, 12, 0, 0, 0, time.UTC)
		recordEndDate := time.Now()
		_, err = db.CreateRecord(userID, "Встреча", "Кофе в кафе", recordStartDate, recordEndDate)
		require.Nil(t, err)

		rec, err := db.GetRecordsByUserID(userID)
		require.Nil(t, err)

		require.Equal(t, "Встреча", rec[0].Title)
		require.Equal(t, "Кофе в кафе", rec[0].Content)

		require.True(t, recordStartDate.Equal(rec[0].DateFrom))
		require.True(t, recordEndDate.Equal(rec[0].DateTo))

		err = db.UpdateRecord(rec[0].ID, rec[0].UserID, "Встреча", "Coffee time", rec[0].DateFrom, rec[0].DateTo)
		require.Nil(t, err)

		rec, _ = db.GetRecordsByUserID(userID)
		require.Equal(t, 1, len(rec))
		require.Equal(t, "Coffee time", rec[0].Content)

		err = db.DeleteRecord(rec[0].ID)
		require.Nil(t, err)

		rec, _ = db.GetRecordsByUserID(userID)
		require.Equal(t, 0, len(rec))
	})
}
