package storage

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewStore_ReturnsStoreInstance(t *testing.T) {
	assert.IsType(t, &Store{}, NewStore())
}

func TestStore_getID_returnsValidId(t *testing.T)  {
	type test struct {
		mem map[int]Todo
		expected int
	}

	tests := []test{
		{
			mem: make(map[int]Todo),
			expected: 1,
		},
		{
			mem: map[int]Todo{
				5: Todo{
					ID: 5,
					Name: "5",
				}	,
				3: Todo{
					ID: 3,
					Name: "3",
				},
			},
			expected: 6,
		},
	}

	for _, ts := range tests {
		s := &Store{
			mem: ts.mem,
		}

		if s.newID() != ts.expected {
			t.Errorf("failed expecting: %v to equal %v", s.newID(), ts.expected)
		}
	}
}

func TestStore_Put(t *testing.T) {
	type test struct {
		mem map[int]Todo
		input Todo
		expectedIndex int
		expectedLen int
	}

	tests := []test{
		{
			mem: make(map[int]Todo),
			input: Todo{
				Name: "h",
			},
			expectedIndex: 1,
			expectedLen: 1,
		},
		{
			mem: map[int]Todo{
				1: Todo{
					Name: "x",
				},
			},
			input: Todo{
				Name: "y",
			},
			expectedIndex: 2,
			expectedLen: 2,
		},
		{
			mem: map[int]Todo{
				1: Todo{
					Name: "x",
				},
			},
			input: Todo{
				ID: 1,
				Name: "y",
			},
			expectedIndex: 1,
			expectedLen: 1,
		},
	}

	for _, ts := range tests {
		s := NewStore()
		s.mem = ts.mem

		id := s.Put(ts.input)

		if v, ok := s.mem[ts.expectedIndex]; ok  {
			assert.Equal(t, ts.input.Name, v.Name)
			assert.Len(t, s.mem, ts.expectedLen)
			assert.Equal(t, ts.expectedIndex, id)
			continue
		}

		t.Errorf("failed expecting Put method to store at index: %v map given: %+v", ts.expectedIndex, s.mem)
	}
}

func TestStore_Del(t *testing.T) {
	type test struct {
		mem map[int]Todo
		input int
		expectedMem map[int]Todo
	}

	tests := []test{
		{
			mem: make(map[int]Todo),
			input: 1,
			expectedMem: make(map[int]Todo),
		},
		{
			mem: map[int]Todo{
				1: Todo{
					Name: "h",
				},
			},
			input: 1,
			expectedMem: make(map[int]Todo),
		},
		{
			mem: map[int]Todo{
				1: Todo{
					Name: "h",
				},
				2: Todo{
					Name: "x",
				},
				3: Todo{
					Name: "y",
				},
			},
			input: 2,
			expectedMem: map[int]Todo{
				1: Todo{
					Name: "h",
				},
				3: Todo{
					Name: "y",
				},
			},
		},
	}

	for _, ts := range tests {
		s := NewStore()
		s.mem = ts.mem

		s.Del(ts.input)
		assert.EqualValues(t, ts.expectedMem, s.mem)
	}
}

func TestStore_List(t *testing.T) {
	type test struct {
		mem map[int]Todo
		expectedList []Todo
	}

	tests := []test{
		{
			mem: make(map[int]Todo),
			expectedList: []Todo{},
		},
		{
			mem: map[int]Todo{
				5: Todo{
					Name: "h",
				},
				1: Todo{
					Name: "x",
				},
			},
			expectedList: []Todo{
				{
					Name: "x",
				},
				{
					Name: "h",
				},
			},
		},
	}

	for _, ts := range tests {
		s := NewStore()
		s.mem = ts.mem

		assert.EqualValues(t, ts.expectedList,s.List())
	}
}

func TestStore_Flush(t *testing.T) {
	s := NewStore()
	s.mem = map[int]Todo{
		1: Todo{
			Name: "h",
		},
	}

	s.Flush()

	assert.Len(t, s.mem, 0)
}