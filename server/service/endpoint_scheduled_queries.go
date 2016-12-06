package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Get Scheduled Query
////////////////////////////////////////////////////////////////////////////////

type getScheduledQueryRequest struct {
	ID uint
}

type scheduledQueryResponse struct {
	kolide.PackQuery
}

type getScheduledQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r getScheduledQueryResponse) error() error { return r.Err }

func makeGetScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getScheduledQueryRequest)

		// TODO: call service
		_ = req

		return getScheduledQueryResponse{
			Scheduled: scheduledQueryResponse{},
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Scheduled Queries In Pack
////////////////////////////////////////////////////////////////////////////////

type getScheduledQueriesInPackRequest struct {
	ID uint
}

type getScheduledQueriesInPackResponse struct {
	Scheduled []scheduledQueryResponse `json:"scheduled"`
	Err       error                    `json:"error,omitempty"`
}

func (r getScheduledQueriesInPackResponse) error() error { return r.Err }

func makeGetScheduledQueriesInPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getScheduledQueriesInPackRequest)
		resp := getScheduledQueriesInPackResponse{Scheduled: []scheduledQueryResponse{}}

		// TODO: call service
		_ = req

		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Schedule Queries
////////////////////////////////////////////////////////////////////////////////

type scheduleQueriesRequest struct {
	// TODO: add fields
}

type scheduleQueriesResponse struct {
	Scheduled []scheduledQueryResponse `json:"scheduled"`
	Err       error                    `json:"error,omitempty"`
}

func (r scheduleQueriesResponse) error() error { return r.Err }

func makeScheduleQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(scheduleQueriesRequest)
		resp := getScheduledQueriesInPackResponse{Scheduled: []scheduledQueryResponse{}}

		// TODO: call service
		_ = req

		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify Scheduled Query
////////////////////////////////////////////////////////////////////////////////

type modifyScheduledQueryRequest struct {
	ID uint
	// payload kolide.PackQueryPayload ??
	// TODO: add fields
}

type modifyScheduledQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r modifyScheduledQueryResponse) error() error { return r.Err }

func makeModifyScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyScheduledQueryRequest)

		// TDOD: call service
		_ = req

		return modifyScheduledQueryResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Scheduled Query
////////////////////////////////////////////////////////////////////////////////

type deleteScheduledQueryRequest struct {
	ID uint
}

type deleteScheduledQueryResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteScheduledQueryResponse) error() error { return r.Err }

func makeDeleteScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteScheduledQueryRequest)

		// TODO: call service
		_ = req

		return deleteScheduledQueryResponse{}, nil
	}
}
