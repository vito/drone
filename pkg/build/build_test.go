package build

import (
	"github.com/drone/drone/pkg/build/script"
	"testing"

	"github.com/drone/drone/pkg/build/docker"
	"github.com/drone/drone/pkg/build/docker/fakedocker"
	"github.com/drone/drone/pkg/build/repo"
)

func TestPrivilegedBuilds(t *testing.T) {
	fakeContainers := fakedocker.NewFakeContainerService()
	fakeImages := fakedocker.NewFakeImageService()

	dockerClient := docker.New()
	dockerClient.Containers = fakeContainers
	dockerClient.Images = fakeImages

	builder := New(dockerClient)

	builder.Build = &script.Build{
		Image: "some-image",
	}

	builder.Repo = &repo.Repo{
		Path:       "https://github.com/drone/drone",
		Dir:        "/var/cache/drone/src",
		Privileged: true,
	}

	// just have to return something so it's found
	fakeImages.InspectResult["some-image"] = nil

	err := builder.Run()
	if err != nil {
		t.Error(err)
	}

	if len(fakeContainers.Created) > 1 {
		t.Error("created too many containers?!")
	}

	if len(fakeContainers.Created) < 1 {
		t.Error("did not create a container")
	}

	for id, _ := range fakeContainers.Created {
		started, ok := fakeContainers.Started[id]
		if !ok {
			t.Error("did not start the container")
		}

		if !started.Privileged {
			t.Error("container should have been privileged")
		}
	}
}

func TestPrivilegedBuildsWithPullRequests(t *testing.T) {
	fakeContainers := fakedocker.NewFakeContainerService()
	fakeImages := fakedocker.NewFakeImageService()

	dockerClient := docker.New()
	dockerClient.Containers = fakeContainers
	dockerClient.Images = fakeImages

	builder := New(dockerClient)

	builder.Build = &script.Build{
		Image: "some-image",
	}

	builder.Repo = &repo.Repo{
		Path:       "https://github.com/drone/drone",
		Dir:        "/var/cache/drone/src",
		PR:         "some-dangerous-pr",
		Privileged: true,
	}

	// just have to return something so it's found
	fakeImages.InspectResult["some-image"] = nil

	err := builder.Run()
	if err != nil {
		t.Error(err)
	}

	if len(fakeContainers.Created) > 1 {
		t.Error("created too many containers?!")
	}

	if len(fakeContainers.Created) < 1 {
		t.Error("did not create a container")
	}

	for id, _ := range fakeContainers.Created {
		started, ok := fakeContainers.Started[id]
		if !ok {
			t.Error("did not start the container")
		}

		if started.Privileged {
			t.Error("container should NOT have been privileged")
		}
	}
}
