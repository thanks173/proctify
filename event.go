package proctify

const (
	OperationDestroyed = iota
	OperationCreated
	OperationModified
)

type Event struct {
	Process   Process
	Operation int
}
