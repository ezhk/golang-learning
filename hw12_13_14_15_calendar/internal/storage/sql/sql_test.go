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

		user, err := db.CreateUser("vasya@pupkin.com", "Vasya", "Pupkin")
		require.NoError(t, err)

		user.Email = "винни@пух.рф"
		user.FirstName = "Винни"
		user.LastName = "Пух"

		err = db.UpdateUser(user)
		require.NoError(t, err)

		u, err := db.GetUserByEmail("винни@пух.рф")
		require.NoError(t, err)

		require.Equal(t, "Винни", u.FirstName)
		require.Equal(t, "Пух", u.LastName)

		_, err = db.CreateUser("винни@пух.рф", "Винни", "Пух")
		require.NotNil(t, err)
		require.Equal(t, err, storage.ErrUserExists)

		err = db.DeleteUser(user)
		require.NoError(t, err)

		// secord run not cause error "user not exist" in SQL
		err = db.DeleteUser(user)
		require.Nil(t, err)
	})

	t.Run("calendar operations", func(t *testing.T) {
		db := NewDatatabase()

		err := db.Connect("user=postgres password=postgres dbname=postgres sslmode=disable")
		require.NoError(t, err)
		defer db.Close()

		// calendar operations must contain user
		user, err := db.CreateUser("vinny@pooh.com", "Винни", "Пух")
		require.Nil(t, err)
		defer db.DeleteUser(user)

		// userID int64, title, content string, dateFrom, dateTo time.Time
		recordStartDate := time.Date(2020, time.September, 1, 12, 0, 0, 0, time.UTC)
		recordEndDate := time.Now()
		event, err := db.CreateEvent(user.ID, "Встреча", "Кофе в кафе", recordStartDate, recordEndDate)
		require.Nil(t, err)

		rec, err := db.GetEventsByUserID(user.ID)
		require.Nil(t, err)

		require.Equal(t, "Встреча", rec[0].Title)
		require.Equal(t, "Кофе в кафе", rec[0].Content)

		require.True(t, recordStartDate.Equal(rec[0].DateFrom))
		require.True(t, recordEndDate.Equal(rec[0].DateTo))

		event.Content = "Coffee time"
		err = db.UpdateEvent(event)
		require.Nil(t, err)

		rec, _ = db.GetEventsByUserID(user.ID)
		require.Equal(t, 1, len(rec))
		require.Equal(t, "Coffee time", rec[0].Content)

		err = db.DeleteEvent(event)
		require.Nil(t, err)

		rec, _ = db.GetEventsByUserID(user.ID)
		require.Equal(t, 0, len(rec))
	})

	t.Run("range events", func(t *testing.T) {
		db := NewDatatabase()

		err := db.Connect("user=postgres password=postgres dbname=postgres sslmode=disable")
		require.NoError(t, err)
		defer db.Close()

		// calendar operations must contain user
		user, err := db.CreateUser("test@user.com", "", "")
		require.Nil(t, err)
		defer db.DeleteUser(user)

		// 2020.09.14 is monday.
		timeToday := time.Date(2020, 9, 14, 1, 19, 31, 0, time.UTC)
		timeThisWeek := timeToday.Add(3 * time.Hour * 24)
		timeThisMonth := timeToday.Add(8 * time.Hour * 24)
		timeNextMonth := timeToday.Add(31 * time.Hour * 24)

		todayEvent, err := db.CreateEvent(user.ID, "Сегодня", "Событие на сегодня", timeToday, timeToday.Add(1*time.Hour))
		require.Nil(t, err)
		defer db.DeleteEvent(todayEvent)

		thisWeekEvent, err := db.CreateEvent(user.ID, "На неделе", "Случайное событие на неделе", timeThisWeek, timeThisWeek.Add(1*time.Hour))
		require.Nil(t, err)
		defer db.DeleteEvent(thisWeekEvent)

		thisMonthEvent, err := db.CreateEvent(user.ID, "Раз в месяц", "Проверить счета", timeThisMonth, timeThisMonth.Add(1*time.Hour))
		require.Nil(t, err)
		defer db.DeleteEvent(thisMonthEvent)

		nextMonthEvent, err := db.CreateEvent(user.ID, "Раз в месяц", "Проверить счета", timeNextMonth, timeNextMonth.Add(1*time.Hour))
		require.Nil(t, err)
		defer db.DeleteEvent(nextMonthEvent)

		// Daily events check.
		dailyEvents, err := db.DailyEvents(user.ID, timeToday)
		require.Nil(t, err)
		require.Equal(t, todayEvent.ID, dailyEvents[0].ID)
		require.Equal(t, 1, len(dailyEvents))

		// Weekly events check.
		weeklyEvents, err := db.WeeklyEvents(user.ID, timeToday)
		require.Nil(t, err)
		weeklyIDs := make([]int64, 0)
		for _, event := range weeklyEvents {
			weeklyIDs = append(weeklyIDs, event.ID)
		}
		require.Contains(t, weeklyIDs, todayEvent.ID)
		require.Contains(t, weeklyIDs, thisWeekEvent.ID)
		require.NotContains(t, weeklyIDs, thisMonthEvent.ID)

		// Monthly events check.
		monthlyEvents, err := db.MonthlyEvents(user.ID, timeToday)
		require.Nil(t, err)
		monthlyIDs := make([]int64, 0)
		for _, event := range monthlyEvents {
			monthlyIDs = append(monthlyIDs, event.ID)
		}
		require.Contains(t, weeklyIDs, todayEvent.ID)
		require.Contains(t, weeklyIDs, thisWeekEvent.ID)
		require.Contains(t, monthlyIDs, thisMonthEvent.ID)
		require.NotContains(t, monthlyIDs, nextMonthEvent.ID)
	})

	t.Run("notify events", func(t *testing.T) {
		db := NewDatatabase()

		err := db.Connect("user=postgres password=postgres dbname=postgres sslmode=disable")
		require.NoError(t, err)
		defer db.Close()

		// calendar operations must contain user
		user, err := db.CreateUser("test@user.com", "", "")
		require.Nil(t, err)
		defer db.DeleteUser(user)

		timeNow := time.Now()
		timeThreeWeek := timeNow.Add(3 * time.Hour * 24 * 7)

		todayEvent, err := db.CreateEvent(user.ID, "Сегодня", "Событие на сегодня", timeNow, timeNow.Add(1*time.Hour))
		require.Nil(t, err)
		defer db.DeleteEvent(todayEvent)

		threeWeekEvent, err := db.CreateEvent(user.ID, "Сегодня", "Событие на сегодня", timeThreeWeek, timeThreeWeek.Add(1*time.Hour))
		require.Nil(t, err)
		defer db.DeleteEvent(threeWeekEvent)

		events, err := db.GetNotifyReadyEvents()
		require.Nil(t, err)
		IDs := make([]int64, 0)
		for _, event := range events {
			IDs = append(IDs, event.ID)
		}
		require.Contains(t, IDs, todayEvent.ID)
		require.NotContains(t, IDs, threeWeekEvent.ID)

		// Mark event.
		require.False(t, todayEvent.Notified)
		err = db.MarkEventAsNotified(&todayEvent)
		require.Nil(t, err)
		require.True(t, todayEvent.Notified)

		// Empty notified events in two weeks.
		events, err = db.GetNotifyReadyEvents()
		require.Nil(t, err)
		require.Equal(t, 0, len(events))
	})
}
