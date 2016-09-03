package kitserver

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	return getUserRequest{ID: uid}, nil
}

func decodeChangePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.UserID = uid
	return req, nil
}

func decodeUpdateAdminRoleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	var req updateAdminRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.UserID = uid
	return req, nil
}

func decodeUpdateUserStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	var req updateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.UserID = uid
	return req, nil
}

func decodeModifyUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	var req modifyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	req.ID = uid
	return req, nil
}
