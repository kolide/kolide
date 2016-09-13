package datastore

import "github.com/kolide/kolide-ose/kolide"

func (orm *inmem) NewPasswordResetRequest(req *kolide.PasswordResetRequest) (*kolide.PasswordResetRequest, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	req.ID = uint(len(orm.passwordResets) + 1)
	orm.passwordResets[req.ID] = req
	return req, nil
}

func (orm *inmem) SavePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if _, ok := orm.passwordResets[req.ID]; !ok {
		return ErrNotFound
	}

	orm.passwordResets[req.ID] = req
	return nil
}

func (orm *inmem) DeletePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if _, ok := orm.passwordResets[req.ID]; !ok {
		return ErrNotFound
	}

	delete(orm.passwordResets, req.ID)
	return nil
}

func (orm *inmem) FindPassswordResetByID(id uint) (*kolide.PasswordResetRequest, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if req, ok := orm.passwordResets[id]; ok {
		return req, nil
	}

	return nil, ErrNotFound
}

func (orm *inmem) FindPassswordResetsByUserID(userID uint) (*kolide.PasswordResetRequest, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, pr := range orm.passwordResets {
		if pr.UserID == userID {
			return pr, nil
		}
	}

	return nil, ErrNotFound

}

func (orm *inmem) FindPassswordResetByToken(token string) (*kolide.PasswordResetRequest, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, pr := range orm.passwordResets {
		if pr.Token == token {
			return pr, nil
		}
	}

	return nil, ErrNotFound
}

func (orm *inmem) FindPassswordResetByTokenAndUserID(token string, userID uint) (*kolide.PasswordResetRequest, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, pr := range orm.passwordResets {
		if pr.Token == token && pr.UserID == userID {
			return pr, nil
		}
	}

	return nil, ErrNotFound
}
