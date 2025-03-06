package in_memory

type InMemoryQueue struct {
	messages []string
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		messages: []string{},
	}
}

func (w *InMemoryQueue) WriteMessage(message []byte) error {
	w.messages = append(w.messages, string(message))
	return nil
}
