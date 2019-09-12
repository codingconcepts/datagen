package runner

import (
	"fmt"
	"math/rand"
	"sync"
)

type rowKey struct {
	groupType interface{}
	groupID   int
}

// store holds row data that comes out of the database during runtime.
type store struct {
	mu          sync.RWMutex
	data        map[string][]map[string]interface{}
	group       map[rowKey]map[string]interface{}
	eachContext string
	eachGroup   int
	eachRow     int
}

func newStore() *store {
	return &store{
		data:  map[string][]map[string]interface{}{},
		group: map[rowKey]map[string]interface{}{},
	}
}

func (s *store) set(groupName string, rows map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[groupName] = append(s.data[groupName], rows)
}

func (s *store) reference(key string, column string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("data not found key=%q", key)
	}

	index := rand.Intn(len(rows))
	value, ok := rows[index][column]
	if !ok {
		return nil, fmt.Errorf("data not found key=%q column=%q index=%d", key, column, index)
	}

	return value, nil
}

func (s *store) row(key, column string, group int) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	groupKey := rowKey{groupType: key, groupID: group}

	// Check if we've scanned this row before.
	row, ok := s.group[groupKey]
	if ok {
		value, ok := row[column]
		if !ok {
			return nil, fmt.Errorf("data not found key=%q column=%q group=%d", key, column, group)
		}
		return value, nil
	}

	// Get a random item from the row context and cache it for the next read.
	randomValue := s.data[key][rand.Intn(len(s.data[key]))]

	s.group[groupKey] = randomValue

	value, ok := randomValue[column]
	if !ok {
		return nil, fmt.Errorf("data not found key=%q column=%q", key, column)
	}

	return value, nil
}

func (s *store) each(key, column string, group int) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	groupKey := rowKey{groupType: key, groupID: group}

	// Check if we've scanned this row before.
	row, ok := s.group[groupKey]
	if ok {
		value, ok := row[column]
		if !ok {
			return nil, fmt.Errorf("data not found key=%q column=%q group=%d", key, column, group)
		}
		return value, nil
	}

	// Get a random item from the row context and cache it for the next read.
	randomValue := s.data[key][s.eachRow]

	s.group[groupKey] = randomValue

	value, ok := randomValue[column]
	if !ok {
		return nil, fmt.Errorf("data not found key=%q column=%q", key, column)
	}

	// Increment the row if the group has changed.
	if s.eachGroup != group {
		s.eachRow++
		s.eachGroup = group
	}

	return value, nil
}
