package core

type TaskInterFace interface {
	GetName() string
	GetConnType() ConnType
	Handel(string) error
	Enable() bool
	GetConsumerNumber() int
}
