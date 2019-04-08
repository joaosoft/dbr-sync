package profile

import (
	"github.com/joaosoft/logger"
)

type IStorageDB interface {
	GetSections() (Sections, error)
	GetSection(sectionKey string) (*Section, error)
	GetSectionContents(sectionKey string) (Contents, error)
}

type Interactor struct {
	logger  logger.ILogger
	storage IStorageDB
}

func NewInteractor(logger logger.ILogger, storageDB IStorageDB) *Interactor {
	return &Interactor{
		logger:  logger,
		storage: storageDB,
	}
}

func (i *Interactor) GetSections() (Sections, error) {
	i.logger.WithFields(map[string]interface{}{"method": "GetSections"})
	i.logger.Info("getting sections")
	sections, err := i.storage.GetSections()
	if err != nil {
		i.logger.WithFields(map[string]interface{}{"error": err.Error()}).Errorf("error getting sections %s", err).ToError()
		return nil, err
	}

	return sections, err
}

func (i *Interactor) GetSection(request *GetSectionRequest) (*Section, error) {
	i.logger.WithFields(map[string]interface{}{"method": "GetSection"})

	i.logger.Infof("getting section [section key: %s]", request.SectionKey)

	section, err := i.storage.GetSection(request.SectionKey)
	if err != nil {
		i.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting section [section key: %s] storage database %s", request.SectionKey, err).ToError()
		return nil, err
	}

	return section, err
}

func (i *Interactor) GetSectionContents(request *GetSectionContentsRequest) (Contents, error) {
	i.logger.WithFields(map[string]interface{}{"method": "GetSection"})

	i.logger.Infof("getting section contents [section key: %s]", request.SectionKey)

	sections, err := i.storage.GetSectionContents(request.SectionKey)
	if err != nil {
		i.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting section contents [section key: %s] storage database %s", request.SectionKey, err).ToError()
		return nil, err
	}

	return sections, err
}
