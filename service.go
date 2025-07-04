package service

import "context"

type Service interface {
	Name() string
	Version() string
	Init(context.Context) error
	Run(context.Context) error
	Shutdown(context.Context) error
}
