package orcareaper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	orcaServerGroupsPath     = "/applications/orca/serverGroups"
	orcaTasksPath            = "/applications/orca/tasks"
	orcaExecutionsStatusPath = "/executions/activeByInstance"
)

func Run(config *Config) error {
	if err := initHTTPClient(config); err != nil {
		return errors.Wrap(err, "initializing http client")
	}

	orcaServerGroups, err := getOrcaServerGroups(config.APIBaseURL)
	if err != nil {
		return errors.Wrap(err, "could not get orca server groups")
	}

	executionsStatus, err := getOrcaExecutionsStatus(config.APIBaseURL)
	if err != nil {
		return errors.Wrap(err, "could not get orca executions status")
	}

	reapTask := buildReapTask(orcaServerGroups, executionsStatus, config)

	if len(reapTask.Job) == 0 {
		log.Println("nothing to cleanup")
		return nil
	}

	log.Printf("%#v", reapTask)

	if config.DryRun {
		log.Println("dry-run mode enabled, skipping action")
		return nil
	}

	body, err := json.Marshal(reapTask)
	if err != nil {
		return errors.Wrap(err, "marshaling reap task to json")
	}
	req, err := http.NewRequest("POST", orcaTasksURL(config.APIBaseURL), bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "creating new post request")
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "performing post reap request")
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Body:", string(body))

	return nil
}

func getOrcaServerGroups(baseURL string) ([]serverGroup, error) {
	r, err := httpClient.Get(orcaServerGroupsURL(baseURL))
	if err != nil {
		return nil, errors.Wrap(err, "requesting gate server groups")
	}
	defer r.Body.Close()

	var orcaServerGroups []serverGroup
	if err := json.NewDecoder(r.Body).Decode(&orcaServerGroups); err != nil {
		raw, _ := ioutil.ReadAll(r.Body)
		return nil, errors.Wrapf(err, "unmarshaling server groups (body: %s)", raw)
	}

	return orcaServerGroups, nil
}

func getOrcaExecutionsStatus(baseURL string) (map[string]int, error) {
	r, err := httpClient.Get(orcaExecutionsStatusURL(baseURL))
	if err != nil {
		return nil, errors.Wrap(err, "requesting orca tasks status")
	}
	defer r.Body.Close()

	var status map[string]int
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		return nil, errors.Wrap(err, "unmarshaling tasks status")
	}

	return status, nil
}

func buildReapTask(orcaServerGroups []serverGroup, tasksStatus map[string]int, config *Config) task {
	reapTask := task{}
	reapTask.Application = "orca"
	reapTask.Description = "Terminate instances and server groups"

	termInstancesJob := terminateInstancesJob{}
	termInstancesJob.Type = "terminateInstances"
	termInstancesJob.Region = config.Region
	termInstancesJob.Credentials = config.Credentials
	for _, serverGroup := range orcaServerGroups {
		if !serverGroup.Disabled {
			continue
		}
		if config.Cluster != "" && config.Cluster != serverGroup.Cluster {
			continue
		}

		if len(serverGroup.Instances) == 0 {
			destroyJob := destroyServerGroupJob{}
			destroyJob.Type = "destroyServerGroup"
			destroyJob.ServerGroupName = serverGroup.Name
			destroyJob.Region = serverGroup.Region
			destroyJob.Credentials = serverGroup.Account
			// TODO rz - Update to allow k8s, etc
			destroyJob.CloudProvider = "aws"
			reapTask.Job = append(reapTask.Job, destroyJob)
			continue
		}
		for _, i := range serverGroup.Instances {
			// if unhealthy, term
			if strings.ToLower(i.HealthState) != "up" {
				termInstancesJob.InstanceIDs = append(termInstancesJob.InstanceIDs, i.ID)
				continue
			}

			// unknown or inactive, term
			execs, ok := tasksStatus[i.ID]
			if !ok || execs == 0 {
				termInstancesJob.InstanceIDs = append(termInstancesJob.InstanceIDs, i.ID)
			}
		}
	}

	if len(termInstancesJob.InstanceIDs) > 0 {
		reapTask.Job = append(reapTask.Job, termInstancesJob)
	}

	return reapTask
}

func orcaServerGroupsURL(baseURL string) string {
	if strings.HasSuffix(baseURL, "/") {
		return baseURL[1:] + orcaServerGroupsPath
	}
	return baseURL + orcaServerGroupsPath
}

func orcaTasksURL(baseURL string) string {
	if strings.HasSuffix(baseURL, "/") {
		return baseURL[1:] + orcaTasksPath
	}
	return baseURL + orcaTasksPath
}

func orcaExecutionsStatusURL(baseURL string) string {
	if strings.HasSuffix(baseURL, "/") {
		return baseURL[1:] + orcaExecutionsStatusPath
	}
	return baseURL + orcaExecutionsStatusPath
}
