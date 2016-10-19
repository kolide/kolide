package service

import (
	"testing"
	"time"

	"github.com/kolide/kolide-ose/server/datastore"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

const bcryptCost = 6

func TestAuthenticate(t *testing.T) {
	ds, err := datastore.New("inmem", "")
	require.Nil(t, err)
	svc, err := newTestService(ds)
	require.Nil(t, err)
	users := createTestUsers(t, ds)

	var loginTests = []struct {
		username string
		password string
		user     kolide.User
		wantErr  error
	}{
		{
			user:     users["admin1"],
			username: testUsers["admin1"].Username,
			password: testUsers["admin1"].PlaintextPassword,
		},
		{
			user:     users["user1"],
			username: testUsers["user1"].Email,
			password: testUsers["user1"].PlaintextPassword,
		},
	}

	for _, tt := range loginTests {
		t.Run(tt.username, func(st *testing.T) {
			user := tt.user
			ctx := context.Background()
			loggedIn, token, err := svc.Login(ctx, tt.username, tt.password)
			require.Nil(st, err, "login unsuccesful")
			assert.Equal(st, user.ID, loggedIn.ID)
			assert.NotEmpty(st, token)

			sessions, err := svc.GetInfoAboutSessionsForUser(ctx, user.ID)
			require.Nil(st, err)
			require.Len(st, sessions, 1, "user should have one session")
			session := sessions[0]
			assert.Equal(st, user.ID, session.UserID)
			assert.WithinDuration(st, time.Now(), session.AccessedAt, 3*time.Second,
				"access time should be set with current time at session creation")
		})
	}

}
