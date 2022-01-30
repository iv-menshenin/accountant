package memory

import (
	"sync"

	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	Memory struct {
		mux  sync.RWMutex
		data map[uuid.UUID]interface{}
	}
)

func New() *Memory {
	return &Memory{
		data: make(map[uuid.UUID]interface{}),
	}
}

func (m *Memory) Create(id uuid.UUID, data interface{}) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[id]; ok {
		return ErrDuplicate
	}
	m.data[id] = data
	return nil
}

func (m *Memory) Replace(id uuid.UUID, data interface{}) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[id]; !ok {
		return ErrNotFound
	}
	m.data[id] = data
	return nil
}

func (m *Memory) Delete(id uuid.UUID) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[id]; !ok {
		return ErrNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *Memory) Lookup(id uuid.UUID) (interface{}, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if data, ok := m.data[id]; ok {
		return data, nil
	}
	return nil, ErrNotFound
}

func (m *Memory) Find(filter func(interface{}) bool) []interface{} {
	m.mux.RLock()
	defer m.mux.RUnlock()
	var result = make([]interface{}, 0)
	for _, v := range m.data {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}
