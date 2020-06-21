package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
	ShowElements() []listItem          // показать все элементы списка
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	// using on delete validation
	elementsSet map[*listItem]struct{}

	// Length stay for back compatibility when we cannot use map
	Length int

	First *listItem
	Last  *listItem
}

func NewList() List {
	return &list{
		elementsSet: make(map[*listItem]struct{}),
		Length:      0,
	}
}

func (l *list) Len() int {
	// also we can use len(l.elementsSet)
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
	l.elementsSet[l.First] = struct{}{}

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
	l.elementsSet[l.Last] = struct{}{}

	return l.Last
}

func (l *list) Remove(i *listItem) {
	if _, ok := l.elementsSet[i]; !ok {
		return
	}

	switch l.Len() {
	case 0:
		// back compatibility for disabled elementsSet
		return
	case 1:
		l.First = nil
		l.Last = nil
	case 2:
		if i == l.First {
			l.First = l.Last
		} else {
			l.Last = l.First
		}
		l.First.Next = nil
		l.Last.Prev = nil
	default:
		switch {
		// remove first element
		case i == l.First:
			i.Next.Prev = nil
			l.First = i.Next
		// remove last element
		case i == l.Last:
			i.Prev.Next = nil
			l.Last = i.Prev
		// remove middle element
		default:
			i.Prev.Next, i.Next.Prev = i.Next, i.Prev
		}
	}

	delete(l.elementsSet, i)
	l.Length--
}

func (l *list) MoveToFront(i *listItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l *list) ShowElements() []listItem {
	elements := make([]listItem, 0, l.Len())
	for key := range l.elementsSet {
		elements = append(elements, *key)
	}
	return elements
}
