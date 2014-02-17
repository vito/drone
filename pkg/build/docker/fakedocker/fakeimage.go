package fakedocker

import (
	"github.com/drone/drone/pkg/build/docker"
)

type FakeImageService struct {
	InspectResult map[string]*docker.Image

	Pulled    []string
	PullError error

	Removed     []string
	RemoveError error

	Built      map[string]*docker.Image
	BuildError error
}

func NewFakeImageService() *FakeImageService {
	return &FakeImageService{
		InspectResult: make(map[string]*docker.Image),
		Built:         make(map[string]*docker.Image),
	}
}

func (i *FakeImageService) List() ([]*docker.Images, error) {
	panic("NOOP")
	return nil, nil
}

func (i *FakeImageService) Create(image string) error {
	panic("NOOP")
	return nil
}

func (i *FakeImageService) Pull(image string) error {
	if i.PullError != nil {
		return i.PullError
	}

	i.Pulled = append(i.Pulled, image)

	return nil
}

func (i *FakeImageService) PullTag(name, tag string) error {
	panic("NOOP")
	return nil
}

func (i *FakeImageService) Remove(image string) ([]*docker.Delete, error) {
	if i.RemoveError != nil {
		return nil, i.RemoveError
	}

	i.Removed = append(i.Removed, image)

	return []*docker.Delete{{}}, nil
}

func (i *FakeImageService) Inspect(image string) (*docker.Image, error) {
	injected, ok := i.InspectResult[image]
	if ok {
		return injected, nil
	}

	built, ok := i.Built[image]
	if ok {
		return built, nil
	}

	return nil, docker.ErrNotFound
}

func (i *FakeImageService) Build(tag, dir string) error {
	if i.BuildError != nil {
		return i.BuildError
	}

	i.Built[tag] = &docker.Image{ID: tag}

	return nil
}
