package testing

import (
	"fmt"
	"github.com/korovaisdead/go-simple-membership/config"
	"github.com/korovaisdead/go-simple-membership/storage"
	"github.com/korovaisdead/go-simple-membership/utils/docker"
	"os/exec"
)

var (
	mongoImage = "mongo"
)

func Setup() string {
	config.BuildTestConfig()
	storage.BuildTestRedisClient()

	if _, err := exec.LookPath("docker"); err != nil {
		panic("Don't hace docker installed in os")
	}

	if ok, err := docker.DockerHaveImage(mongoImage); !ok || err != nil {
		if err != nil {
			panic(fmt.Sprintf("Error running docker to check for %s: %v", mongoImage, err))
		}
		if err := docker.DockerPull(mongoImage); err != nil {
			panic(fmt.Sprintf("Error pulling %s: %v", mongoImage, err))
		}
	}

	containerID, err := docker.DockerRun("-d", "-p", "27018:27017", mongoImage)
	if err != nil {
		panic("failed to run docker container")
	}

	return containerID
}

func Shutdown(containerID string) {
	docker.DockerKillContainer(containerID)
}
