// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/kolide/kolide/server/kolide"

var _ kolide.FileIntegrityMonitoringStore = (*FileIntegrityMonitoringStore)(nil)

type NewFIMSectionFunc func(path *kolide.FIMSection) (*kolide.FIMSection, error)

type FIMSectionsFunc func() (kolide.FIMSections, error)

type FileIntegrityMonitoringStore struct {
	NewFIMSectionFunc        NewFIMSectionFunc
	NewFIMSectionFuncInvoked bool

	FIMSectionsFunc        FIMSectionsFunc
	FIMSectionsFuncInvoked bool
}

func (s *FileIntegrityMonitoringStore) NewFIMSection(path *kolide.FIMSection) (*kolide.FIMSection, error) {
	s.NewFIMSectionFuncInvoked = true
	return s.NewFIMSectionFunc(path)
}

func (s *FileIntegrityMonitoringStore) FIMSections() (kolide.FIMSections, error) {
	s.FIMSectionsFuncInvoked = true
	return s.FIMSectionsFunc()
}
