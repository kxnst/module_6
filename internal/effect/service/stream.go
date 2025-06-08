package service

type Stream interface {
	Start() error
	Stop() error
	Close() error
}
