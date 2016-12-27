package inmem

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/patrickmn/sortutil"
)

func (d *Datastore) OptionByName(name string) (*kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	for _, opt := range d.options {
		if opt.Name == name {
			result := *opt
			return &result, nil
		}
	}
	return nil, notFound("options")
}

type optPair struct {
	newOpt      kolide.Option
	existingOpt *kolide.Option
}

func (d *Datastore) SaveOptions(opts []kolide.Option) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	validPairs := []optPair{}
	for _, opt := range opts {
		if opt.ReadOnly {
			return fmt.Errorf("readonly option can't be changed")
		}
		existing, ok := d.options[opt.ID]
		if !ok {
			return notFound("option")
		}
		if existing.Type != opt.Type {
			return fmt.Errorf("type mismatch for option")
		}
		validPairs = append(validPairs, optPair{opt, existing})
	}
	// if all the options to be modified pass validation copy values over to
	// existing options
	if len(validPairs) == len(opts) {
		for _, pair := range validPairs {
			if pair.newOpt.Value == nil {
				pair.existingOpt.Value = nil
				continue
			}
			pair.existingOpt.Value = pair.newOpt.Value
		}
	}
	return nil
}

func (d *Datastore) Option(id uint) (*kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	saved, ok := d.options[id]
	if !ok {
		return nil, notFound("Option").WithID(id)
	}
	result := *saved
	return &result, nil
}

func (d *Datastore) Options() ([]kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	result := []kolide.Option{}
	for _, opt := range d.options {
		result = append(result, *opt)
	}
	sortutil.AscByField(result, "Name")
	return result, nil
}
