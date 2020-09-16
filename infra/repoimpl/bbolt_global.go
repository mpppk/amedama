package repoimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
	bolt "go.etcd.io/bbolt"
)

const globalBucketName = "Global"

var workSpacesKey = []byte("WorkSpaces")

var globalBucketNames = []string{globalBucketName}

type BBoltGlobal struct {
	base           *boltRepository
	pathRepository *pathRepository
}

func NewBBoltGlobal(b *bolt.DB) repository.Global {
	return &BBoltGlobal{
		base:           newBoltRepository(b),
		pathRepository: newPathRepository(b),
	}
}

func (b *BBoltGlobal) globalBucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.base.bucketFunc(globalBucketNames, f)
}

func (b *BBoltGlobal) getWorkSpacesFromBucket(bucket *bolt.Bucket) (workspaces []*model.WorkSpace, err error) {
	workSpacesBytes := bucket.Get(workSpacesKey)
	if workSpacesBytes == nil {
		return nil, nil
	}
	err = json.Unmarshal(workSpacesBytes, &workspaces)
	return
}

func (b *BBoltGlobal) setWorkSpaces(bucket *bolt.Bucket, workspaces []*model.WorkSpace) error {
	workspacesBytes, err := json.Marshal(workspaces)
	if err != nil {
		return fmt.Errorf("failed to marshal workspaces: %w", err)
	}
	return bucket.Put(workSpacesKey, workspacesBytes)
}

func (b *BBoltGlobal) updateWorkSpaces(f func(workspaces []*model.WorkSpace) ([]*model.WorkSpace, error)) error {
	return b.globalBucketFunc(func(bucket *bolt.Bucket) error {
		workspaces, err := b.getWorkSpacesFromBucket(bucket)
		if err != nil {
			return err
		}
		newWorkspaces, err := f(workspaces)
		if err != nil {
			return err
		}
		return b.setWorkSpaces(bucket, newWorkspaces)
	})
}

func (b *BBoltGlobal) ListWorkSpace() (workspaces []*model.WorkSpace, err error) {
	err = b.globalBucketFunc(func(bucket *bolt.Bucket) error {
		workspaces, err = b.getWorkSpacesFromBucket(bucket)
		return err
	})
	return
}

func (b *BBoltGlobal) AddWorkSpace(ws *model.WorkSpace) error {
	return b.updateWorkSpaces(func(workspaces []*model.WorkSpace) ([]*model.WorkSpace, error) {
		return append(workspaces, ws), nil
	})
}
