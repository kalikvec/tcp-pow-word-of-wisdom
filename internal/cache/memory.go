package cache

type Memory struct {
	storage map[int]int64
}

func NewMemory() *Memory {
	return &Memory{storage: make(map[int]int64)}
}

func (m *Memory) Add(key int, val int64) error {
	m.storage[key] = val
	return nil
}

func (m *Memory) Get(key int) (bool, error) {
	if _, ok := m.storage[key]; ok {
		return true, nil
	}
	return false, nil
}

func (m *Memory) Delete(key int) {
	delete(m.storage, key)
}
