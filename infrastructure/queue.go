package infrastructure

type Queue interface {
	WriteMessage(message []byte) error
}
