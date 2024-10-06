package server

import "container/list"

type GenericList[T any] struct {
	data list.List
}

func (l *GenericList[T]) Len() int {
	return l.data.Len()
}

func (l *GenericList[T]) PushBack(t T) {
	l.data.PushBack(t)
}

func (l *GenericList[T]) PopFront() T {
	f := l.data.Front()
	v := f.Value.(T)
	l.data.Remove(f)
	return v
}
