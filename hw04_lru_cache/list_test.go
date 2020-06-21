package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("zero or one element", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		require.Equal(t, l.Len(), 1)
		// validate internal struct: when last and first are equals
		require.Equal(t, l.Front(), l.Back())

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, elems, []int{10})
		require.Equal(t, l.ShowElements(), []listItem{*l.Front()})

		l.Remove(l.Front())
		require.Equal(t, l.Len(), 0)
		require.Equal(t, l.Front(), (*listItem)(nil))
		require.Equal(t, l.Back(), (*listItem)(nil))
	})

	t.Run("two elements", func(t *testing.T) {
		l := NewList()
		l.PushFront("begin")
		l.PushBack("end")
		require.Equal(t, l.Len(), 2)

		elems := make([]string, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(string))
		}
		require.Equal(t, elems, []string{"begin", "end"})

		l.Remove(l.Front())
		require.Equal(t, l.Len(), 1)

		// validate last value
		require.Equal(t, l.Back().Value.(string), "end")
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, l.Len(), 2)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
