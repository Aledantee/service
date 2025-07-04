package service

import (
	"context"
	"github.com/aledantee/ae"
)

func RunAndWait(ctx context.Context, service ...Service) (err error) {
	if len(service) == 0 {
		return nil
	}

	m := NewManager()

	var errs []error
	for _, s := range service {
		_, e := m.Run(ctx, s)

		if e != nil {
			errs = append(errs, ae.New().
				Attr("name", s.Name()).
				Attr("version", s.Version()).
				Cause(e).
				Msg("service failed"),
			)
		}
	}

	return m.WaitAll()
}
