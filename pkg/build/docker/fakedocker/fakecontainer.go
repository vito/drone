package fakedocker

import (
	"fmt"
	"io"
	"sync"

	"github.com/drone/drone/pkg/build/docker"
)

type FakeContainerService struct {
	Created     map[string]*docker.Config
	CreateError error

	Started    map[string]*docker.HostConfig
	StartError error

	WaitResults map[string]int
	Waited      []string
	WaitError   error

	Attached    map[string]io.Writer
	AttachError error

	Stopped   []string
	StopError error

	Removed     []string
	RemoveError error

	runCount int

	sync.Mutex
}

func NewFakeContainerService() *FakeContainerService {
	return &FakeContainerService{
		Created:     make(map[string]*docker.Config),
		Started:     make(map[string]*docker.HostConfig),
		WaitResults: make(map[string]int),
		Attached:    make(map[string]io.Writer),
	}
}

func (c *FakeContainerService) List() ([]*docker.Containers, error) {
	panic("NOOP")
	return nil, nil
}

func (c *FakeContainerService) ListAll() ([]*docker.Containers, error) {
	panic("NOOP")
	return nil, nil
}

func (c *FakeContainerService) Create(conf *docker.Config) (*docker.Run, error) {
	if c.CreateError != nil {
		return nil, c.CreateError
	}

	c.Lock()
	defer c.Unlock()

	id := c.nextRunID()

	c.Created[id] = conf

	return &docker.Run{
		ID: id,
	}, nil
}

func (c *FakeContainerService) Start(id string, conf *docker.HostConfig) error {
	if c.StartError != nil {
		return c.StartError
	}

	c.Started[id] = conf

	return nil
}

func (c *FakeContainerService) Stop(id string, timeoutInSeconds int) error {
	if c.StopError != nil {
		return c.StopError
	}

	c.Stopped = append(c.Stopped, id)

	return nil
}

func (c *FakeContainerService) Remove(id string) error {
	if c.RemoveError != nil {
		return c.RemoveError
	}

	c.Removed = append(c.Removed, id)

	return nil
}

func (c *FakeContainerService) Wait(id string) (*docker.Wait, error) {
	if c.WaitError != nil {
		return nil, c.WaitError
	}

	c.Waited = append(c.Waited, id)

	return &docker.Wait{StatusCode: c.WaitResults[id]}, nil
}

func (c *FakeContainerService) Attach(id string, out io.Writer) error {
	if c.AttachError != nil {
		return c.AttachError
	}

	c.Attached[id] = out

	return nil
}

func (c *FakeContainerService) Inspect(id string) (*docker.Container, error) {
	panic("NOOP")
	return nil, nil
}

func (c *FakeContainerService) Run(conf *docker.Config, host *docker.HostConfig, out io.Writer) (*docker.Wait, error) {
	panic("NOOP")
	return nil, nil
}

func (c *FakeContainerService) RunDaemon(conf *docker.Config, host *docker.HostConfig) (*docker.Run, error) {
	panic("NOOP")
	return nil, nil
}

func (c *FakeContainerService) RunDaemonPorts(image string, ports ...string) (*docker.Run, error) {
	panic("NOOP")
	return nil, nil
}

func (c *FakeContainerService) nextRunID() string {
	id := c.runCount

	c.runCount++

	return fmt.Sprintf("run-%d", id)
}
