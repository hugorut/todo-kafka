package storage

import (
	"sort"
	"sync"
)

// Todo is a struct which holds information about the
// the name and id of the todo
type Todo struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

// Store defines a wrapper around a simple in memory
// data store of Todos
type Store struct {
	mem map[int]Todo
	lock *sync.Mutex
}

// NewStore returns a pointer to a Store instance with a
// clean memory map
func NewStore() *Store  {
	return &Store{
		mem: make(map[int]Todo),
		lock: &sync.Mutex{},
	}
}

// S is the default Store which is to be used by the application
// this can be overridden for test by simply changing the reference
var S = NewStore()

// Store.newID method iterates through the stored Todos and returns
// a unused ID to create a new Todo with
func (s *Store) newID() int {
	if len(s.mem) == 0 {
		return 1
	}

	ids := todoIds(s.mem)
	return ids[len(ids)-1]+1
}

// Store.Put method either updates or creates a reference to a Todo
// in the memory map. If there is no given ID in the passed Todo
// a new ID is generated and then added to the memory
func (s *Store) Put(todo Todo) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	if todo.ID == 0 {
		todo.ID = s.newID()
	}

	s.mem[todo.ID] = todo

	return todo.ID
}

// Store.Del deletes a given Todo in the mem map
func (s *Store) Del(id int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.mem, id)
}

// Store.List returns a ordered slice of Todos
func (s *Store) List() []Todo {
	ids := todoIds(s.mem)
	var sorted []Todo = make([]Todo, len(ids))

	var k int
	for _, i := range ids {
		sorted[k] = s.mem[i]
		k++
	}

	return sorted
}

// Store.Flush resets the underlying memory
func (s *Store) Flush()  {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.mem = make(map[int]Todo)
}

// todoIds takes a given map and returns a list of sorted
// ids using its keys
func todoIds(todos map[int]Todo) []int {
	var ids []int = make([]int, len(todos))

	var k int
	for i := range todos {
		ids[k] = i
		k++
	}

	sort.Ints(ids)
	return ids
}

