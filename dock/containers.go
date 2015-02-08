package dock

import (
	"github.com/dinp/common/model"
	"github.com/fsouza/go-dockerclient"
	"log"
	"strings"
)

func Containers(endpoint string) ([]*model.ContainerDto, error) {

	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	cs, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}

	containers := []*model.ContainerDto{}

	for _, c := range cs {
		// inspect every container to get env['APP_NAME']
		ic, err := client.InspectContainer(c.ID)
		if err != nil {
			log.Println("[ERROR] inspect container error:", err)
			continue
		}

		appName, existent := retrieveAppName(ic.Config)
		if !existent {
			continue
		}

		containers = append(containers, &model.ContainerDto{
			Id:      c.ID,
			Image:   c.Image,
			AppName: appName,
			Ports:   buildPorts(c.Ports),
			Status:  c.Status,
		})
	}

	return containers, nil
}

func retrieveAppName(cfg *docker.Config) (string, bool) {
	if cfg == nil || cfg.Env == nil || len(cfg.Env) == 0 {
		return "", false
	}

	for _, env := range cfg.Env {
		if strings.HasPrefix(env, "APP_NAME=") {
			return env[9:], true
		}
	}

	return "", false
}

func buildPorts(ports []docker.APIPort) (ret []*model.Port) {
	if len(ports) == 0 {
		return
	}

	for _, p := range ports {
		ret = append(ret, &model.Port{PublicPort: int(p.PublicPort)})
	}

	return
}
