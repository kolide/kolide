// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/kolide/fleet/server/kolide"

var _ kolide.InviteStore = (*InviteStore)(nil)

type NewInviteFunc func(i *kolide.Invite) (*kolide.Invite, error)

type ListInvitesFunc func(opt kolide.ListOptions) ([]*kolide.Invite, error)

type InviteFunc func(id uint) (*kolide.Invite, error)

type InviteByEmailFunc func(email string) (*kolide.Invite, error)

type InviteByTokenFunc func(token string) (*kolide.Invite, error)

type SaveInviteFunc func(i *kolide.Invite) error

type DeleteInviteFunc func(id uint) error

type InviteStore struct {
	NewInviteFunc        NewInviteFunc
	NewInviteFuncInvoked bool

	ListInvitesFunc        ListInvitesFunc
	ListInvitesFuncInvoked bool

	InviteFunc        InviteFunc
	InviteFuncInvoked bool

	InviteByEmailFunc        InviteByEmailFunc
	InviteByEmailFuncInvoked bool

	InviteByTokenFunc        InviteByTokenFunc
	InviteByTokenFuncInvoked bool

	SaveInviteFunc        SaveInviteFunc
	SaveInviteFuncInvoked bool

	DeleteInviteFunc        DeleteInviteFunc
	DeleteInviteFuncInvoked bool
}

func (s *InviteStore) NewInvite(i *kolide.Invite) (*kolide.Invite, error) {
	s.NewInviteFuncInvoked = true
	return s.NewInviteFunc(i)
}

func (s *InviteStore) ListInvites(opt kolide.ListOptions) ([]*kolide.Invite, error) {
	s.ListInvitesFuncInvoked = true
	return s.ListInvitesFunc(opt)
}

func (s *InviteStore) Invite(id uint) (*kolide.Invite, error) {
	s.InviteFuncInvoked = true
	return s.InviteFunc(id)
}

func (s *InviteStore) InviteByEmail(email string) (*kolide.Invite, error) {
	s.InviteByEmailFuncInvoked = true
	return s.InviteByEmailFunc(email)
}

func (s *InviteStore) InviteByToken(token string) (*kolide.Invite, error) {
	s.InviteByTokenFuncInvoked = true
	return s.InviteByTokenFunc(token)
}

func (s *InviteStore) SaveInvite(i *kolide.Invite) error {
	s.SaveInviteFuncInvoked = true
	return s.SaveInviteFunc(i)
}

func (s *InviteStore) DeleteInvite(id uint) error {
	s.DeleteInviteFuncInvoked = true
	return s.DeleteInviteFunc(id)
}
