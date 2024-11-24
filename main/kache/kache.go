package kache

// StringSet тип для хранения набора строк с ограничением строк
type StringSet struct {
	strings      []string
	maxSize      int
	currentIndex int
}

// NewStringSet создает новый набор строк с указанным максимальным размером
func NewStringSet(maxSize int) *StringSet {
	return &StringSet{
		strings:      make([]string, maxSize),
		maxSize:      maxSize,
		currentIndex: 0,
	}
}

// Add добавляет новую строку в набор, поддерживая ограничение в максимальное количество строк
func (s *StringSet) Add(str string) {
	if s.currentIndex > s.maxSize {
		s.currentIndex = 0
	}
	s.strings[s.currentIndex] = str
	// Иначе просто добавляем новую строку
	s.currentIndex++
}

// GetAll возвращает все строки в наборе
func (s *StringSet) GetAll() []string {
	var filteredStrings []string
	for _, str := range s.strings {
		if str != "" {
			filteredStrings = append(filteredStrings, str)
		}
	}
	return filteredStrings
}
