package api

import "context"

type Server interface {
	Run(ctx context.Context) error
}

type API struct {
	server Server
}

func New(server Server) (*API, error) {
	return &API{server: server}, nil
}

func (a *API) RunAPI(ctx context.Context) error {
	return a.server.Run(ctx)
}
