package queue

import "sync"

type Queue struct {
	items []any
	mutex *sync.Mutex
}

func New() *Queue {
	q := &Queue{
		items: []any{},
		mutex: &sync.Mutex{},
	}

	return q
}

func (q *Queue) Push(item any) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, item)
}

func (q *Queue) Pop() any {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	idx := len(q.items) - 1
	item := q.items[idx]
	q.items = q.items[:idx]
	return item
}

func (q *Queue) NotEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items) > 0
}

func (q *Queue) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items) == 0
}

func (q *Queue) Length() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items)
}

func (q *Queue) GetAll() any {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	items := q.items
	q.items = []any{}
	return items
}
