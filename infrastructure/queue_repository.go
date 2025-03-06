package infrastructure

//go:generate mockgen -destination=./mocks/mock_queue.go -package=mocks github.com/pedro00627/urblog/infrastructure Queue
type Queue interface {
	WriteMessage(message []byte) error
}
