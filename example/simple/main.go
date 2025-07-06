package main

import (
	"context"
	"github.com/aledantee/service"
)

func main() {
	svc := service.Service{
		Name:    "simple",
		Version: "0.0.1",
		Init: func(ctx context.Context) error {
			service.Logger(ctx).Info("init phase")
			return nil
		},
		Run: func(svc *service.Service) error {
			service.Logger(svc.Context()).Info("run phase")
			return nil
		},
		Shutdown: func(ctx context.Context) error {
			service.Logger(ctx).Info("shutdown phase")
			return nil
		},
	}

	svc.ExecuteExit()
}
