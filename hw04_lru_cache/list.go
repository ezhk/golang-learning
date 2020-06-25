package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	// Length stay for back compatibility when we cannot use map
	Length int

	First *listItem
	Last  *listItem
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *listItem {
	return l.First
}

func (l *list) Back() *listItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *listItem {
	switch l.Len() {
	case 0:
		l.First = &listItem{Value: v}
		l.Last = l.First
	default:
		value := &listItem{
			Value: v,
			Next:  l.First,
		}
		l.First.Prev = value
		l.First = value
	}

	l.Length++
	return l.First
}

func (l *list) PushBack(v interface{}) *listItem {
	switch l.Len() {
	case 0:
		l.Last = &listItem{Value: v}
		l.First = l.Last
	default:
		value := &listItem{
			Value: v,
			Prev:  l.Last,
		}
		l.Last.Next = value
		l.Last = value
	}

	l.Length++
	return l.Last
}

func (l *list) Remove(i *listItem) {
	switch {
	// element in the middle
	case i.Prev != nil && i.Next != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	// element in the right/end
	case i.Prev != nil:
		l.Last = i.Prev
		i.Prev.Next = nil
	// element in the left/begin
	case i.Next != nil:
		l.First = i.Next
		i.Next.Prev = nil
	// stay only one element
	default:
		l.First = nil
		l.Last = nil
	}

	l.Length--
}

func (l *list) MoveToFront(i *listItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
