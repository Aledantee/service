package main

import (
	"context"
	"github.com/aledantee/ae"
	"github.com/aledantee/service"
)

func main() {
	svc := service.Service{
		Name:    "error",
		Version: "0.0.1",
		Init: func(ctx context.Context) error {
			cause := ae.New().
				Tag("component1").
				Hint("do the thing").
				ExitCode(200).
				Msg("failed")

			return ae.Wrap("some error", cause)
		},
		Run: func(svc *service.Service) error {
			panic("unreachable")
		},
		Shutdown: func(ctx context.Context) error {
			panic("unreachable")
		},
	}

	svc.ExecuteExit()
}
