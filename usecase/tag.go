package usecase

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type Tag struct {
	tagRepository repository.Tag
}

func NewTag(tagRepository repository.Tag) *Tag {
	return &Tag{
		tagRepository: tagRepository,
	}
}

func (a *Tag) List(ws model.WSName) ([]*model.Tag, error) {
	// FIXME
	if err := a.tagRepository.Init(ws); err != nil {
		return nil, err
	}
	return a.tagRepository.ListAll(ws)
}

// SetTags set tag list to workspace
func (a *Tag) SetTags(ws model.WSName, tags []*model.Tag) error {
	// FIXME
	if err := a.tagRepository.Init(ws); err != nil {
		return err
	}
	if err := a.tagRepository.RecreateBucket(ws); err != nil {
		return fmt.Errorf("failed to set tags: %w", err)
	}
	for _, tag := range tags {
		if err := a.tagRepository.Add(ws, tag); err != nil {
			return fmt.Errorf("failed to set tags: %w", err)
		}
	}
	return nil
}
