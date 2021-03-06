package opts

import (
	"os"
	"strings"

	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/jenkins-x/jx/pkg/log"
)

// DockerRegistryOrg parses the docker registry organisation from various places
func (o *CommonOptions) DockerRegistryOrg(repository *gits.GitRepository) string {
	answer := ""
	teamSettings, err := o.TeamSettings()
	if err != nil {
		log.Warnf("Could not load team settings %s\n", err.Error())
	} else {
		answer = teamSettings.DockerRegistryOrg
	}
	if answer == "" {
		answer = os.Getenv("DOCKER_REGISTRY_ORG")
	}
	if answer == "" && repository != nil {
		answer = repository.Organisation
	}
	return strings.ToLower(answer)
}

// DockerRegistry parses the docker registry from various places
func (o *CommonOptions) DockerRegistry() string {
	dockerRegistry := os.Getenv("DOCKER_REGISTRY")
	if dockerRegistry == "" {
		kubeClient, ns, err := o.KubeClientAndDevNamespace()
		if err != nil {
			log.Warnf("failed to create kube client: %s\n", err.Error())
		} else {
			name := kube.ConfigMapJenkinsDockerRegistry
			data, err := kube.GetConfigMapData(kubeClient, name, ns)
			if err != nil {
				log.Warnf("failed to load ConfigMap %s in namespace %s: %s\n", name, ns, err.Error())
			} else {
				dockerRegistry = data["docker.registry"]
			}
		}
	}
	return dockerRegistry
}
