package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockOsqueryResultHandler struct{}

func (h *MockOsqueryResultHandler) HandleResultLog(log OsqueryResultLog, nodeKey string) error {
	return nil
}

type MockOsqueryStatusHandler struct{}

func (h *MockOsqueryStatusHandler) HandleStatusLog(log OsqueryStatusLog, nodeKey string) error {
	return nil
}

func TestGetAllQueries(t *testing.T) {
	ds := createTestUsers(t, createTestPacksAndQueries(t, createTestDatastore(t)))
	server := createTestServer(ds)

	////////////////////////////////////////////////////////////////////////////
	// try to get queries while logged out
	////////////////////////////////////////////////////////////////////////////

	response := makeRequest(
		t,
		server,
		"GET",
		"/api/v1/kolide/queries",
		nil,
		"",
	)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	////////////////////////////////////////////////////////////////////////////
	// log-in with a user
	////////////////////////////////////////////////////////////////////////////

	// log in with admin test user
	response = makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "user1",
			Password: "foobar",
		},
		"",
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a non-empty cookie was in-fact set
	userCookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, userCookie)

	////////////////////////////////////////////////////////////////////////////
	// get queries from a user account
	////////////////////////////////////////////////////////////////////////////

	response = makeRequest(
		t,
		server,
		"GET",
		"/api/v1/kolide/queries",
		nil,
		userCookie,
	)
	assert.Equal(t, http.StatusOK, response.Code)
	var queries GetAllQueriesResponseBody
	err := json.NewDecoder(response.Body).Decode(&queries)
	assert.Nil(t, err)
	assert.Len(t, queries.Queries, 3)
}

func TestGetQuery(t *testing.T) {
	ds := createTestUsers(t, createTestPacksAndQueries(t, createTestDatastore(t)))
	server := createTestServer(ds)
	queries, err := ds.Queries()
	assert.Nil(t, err)
	assert.NotEmpty(t, queries)
	query := queries[0]

	////////////////////////////////////////////////////////////////////////////
	// try to get query while logged out
	////////////////////////////////////////////////////////////////////////////

	response := makeRequest(
		t,
		server,
		"GET",
		fmt.Sprintf("/api/v1/kolide/queries/%d", query.ID),
		nil,
		"",
	)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	////////////////////////////////////////////////////////////////////////////
	// log-in with a user
	////////////////////////////////////////////////////////////////////////////

	// log in with admin test user
	response = makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "user1",
			Password: "foobar",
		},
		"",
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a non-empty cookie was in-fact set
	userCookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, userCookie)

	////////////////////////////////////////////////////////////////////////////
	// get query from a user account
	////////////////////////////////////////////////////////////////////////////

	response = makeRequest(
		t,
		server,
		"GET",
		fmt.Sprintf("/api/v1/kolide/queries/%d", query.ID),
		nil,
		userCookie,
	)
	assert.Equal(t, http.StatusOK, response.Code)
	var q GetQueryResponseBody
	err = json.NewDecoder(response.Body).Decode(&q)
	assert.Nil(t, err)
	assert.Equal(t, q.Name, query.Name)
}

func TestGetAllPacks(t *testing.T) {
	ds := createTestUsers(t, createTestPacksAndQueries(t, createTestDatastore(t)))
	server := createTestServer(ds)

	////////////////////////////////////////////////////////////////////////////
	// try to get packs while logged out
	////////////////////////////////////////////////////////////////////////////

	response := makeRequest(
		t,
		server,
		"GET",
		"/api/v1/kolide/packs",
		nil,
		"",
	)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	////////////////////////////////////////////////////////////////////////////
	// log-in with a user
	////////////////////////////////////////////////////////////////////////////

	// log in with admin test user
	response = makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "user1",
			Password: "foobar",
		},
		"",
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a non-empty cookie was in-fact set
	userCookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, userCookie)

	////////////////////////////////////////////////////////////////////////////
	// get queries from a user account
	////////////////////////////////////////////////////////////////////////////

	response = makeRequest(
		t,
		server,
		"GET",
		"/api/v1/kolide/packs",
		nil,
		userCookie,
	)
	assert.Equal(t, http.StatusOK, response.Code)
	var packs GetAllPacksResponseBody
	err := json.NewDecoder(response.Body).Decode(&packs)
	assert.Nil(t, err)
	assert.Len(t, packs.Packs, 2)
}
