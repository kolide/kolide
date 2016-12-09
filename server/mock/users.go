package mock

import "github.com/kolide/kolide-ose/server/kolide"
import kolide_errors "github.com/kolide/kolide-ose/server/errors"

var _ kolide.UserStore = (*UserStore)(nil)

type UserByEmailFunc func(email string) (*kolide.User, error)
type UserByIDFunc func(id uint) (*kolide.User, error)

type UserStore struct {
	NewUserFunc        func(user *kolide.User) (*kolide.User, error)
	NewUserFuncInvoked bool

	UserFunc        func(username string) (*kolide.User, error)
	UserFuncInvoked bool

	ListUsersFunc        func(opt kolide.ListOptions) ([]*kolide.User, error)
	ListUsersFuncInvoked bool

	UserByEmailFunc        UserByEmailFunc
	UserByEmailFuncInvoked bool

	UserByIDFunc        UserByIDFunc
	UserByIDFuncInvoked bool

	SaveUserFunc        func(user *kolide.User) error
	SaveUserFuncInvoked bool
}

func (s *UserStore) NewUser(user *kolide.User) (*kolide.User, error) {
	s.NewUserFuncInvoked = true
	return s.NewUserFunc(user)
}

func (s *UserStore) User(username string) (*kolide.User, error) {
	s.UserFuncInvoked = true
	return s.UserFunc(username)
}

func (s *UserStore) ListUsers(opt kolide.ListOptions) ([]*kolide.User, error) {
	s.ListUsersFuncInvoked = true
	return s.ListUsersFunc(opt)
}

func (s *UserStore) UserByEmail(email string) (*kolide.User, error) {
	s.UserByEmailFuncInvoked = true
	return s.UserByEmailFunc(email)
}

func (s *UserStore) UserByID(id uint) (*kolide.User, error) {
	s.UserByIDFuncInvoked = true
	return s.UserByIDFunc(id)
}

func (s *UserStore) SaveUser(user *kolide.User) error {
	s.SaveUserFuncInvoked = true
	return s.SaveUserFunc(user)
}

// helpers

func UserByEmailWithUser(u *kolide.User) UserByEmailFunc {
	return func(email string) (*kolide.User, error) {
		return u, nil
	}
}

func UserWithEmailNotFound() UserByEmailFunc {
	return func(email string) (*kolide.User, error) {
		return nil, kolide_errors.ErrNotFound
	}
}

func UserWithID(u *kolide.User) UserByIDFunc {
	return func(id uint) (*kolide.User, error) {
		return u, nil
	}
}
