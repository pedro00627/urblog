package infrastructure

type InMemoryqueue struct {
	messages []string
}

func NewInMemoryqueue() *InMemoryqueue {
	return &InMemoryqueue{
		messages: []string{},
	}
}

func (w *InMemoryqueue) WriteMessage(message []byte) error {
	w.messages = append(w.messages, string(message))
	return nil
}
