package orcareaper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig = &Config{
	Region:      "us-west-2",
	Credentials: "mgmt",
}

func TestBuildReapTask_NoDisabledGroups(t *testing.T) {
	// given
	orcaServerGroups := []serverGroup{
		{
			Account:  "mgmt",
			Cluster:  "orca-main",
			Name:     "orca-main-v000",
			Region:   "us-west-2",
			Disabled: false,
			Instances: []instance{
				{
					ID:          "v000-1",
					HealthState: "UP",
				},
			},
		},
	}
	executionsStatus := map[string]int{
		"v000-1": 2,
	}

	// when
	task := buildReapTask(orcaServerGroups, executionsStatus, testConfig)

	// then
	assert.Equal(t, task.Application, "orca")
	assert.Nil(t, task.Job)
}

func TestBuildReapTask_NoInactive(t *testing.T) {
	// given
	orcaServerGroups := []serverGroup{
		{
			Account:  "mgmt",
			Cluster:  "orca-main",
			Name:     "orca-main-v000",
			Region:   "us-west-2",
			Disabled: true,
			Instances: []instance{
				{
					ID:          "v000-1",
					HealthState: "UP",
				},
			},
		},
		{
			Account:  "mgmt",
			Cluster:  "orca-main",
			Name:     "orca-main-v001",
			Region:   "us-west-2",
			Disabled: false,
			Instances: []instance{
				{
					ID:          "v001-1",
					HealthState: "UP",
				},
			},
		},
	}
	executionsStatus := map[string]int{
		"v000-1": 2,
		"v001-1": 1,
	}

	// when
	task := buildReapTask(orcaServerGroups, executionsStatus, testConfig)

	// then
	assert.Equal(t, task.Application, "orca")
	assert.Nil(t, task.Job)
}

func TestBuildReapTask_InactiveInstances(t *testing.T) {
	// given
	orcaServerGroups := []serverGroup{
		{
			Account:  "mgmt",
			Cluster:  "orca-main",
			Name:     "orca-main-v000",
			Region:   "us-west-2",
			Disabled: true,
			Instances: []instance{
				{
					ID:          "v000-1",
					HealthState: "UP",
				},
			},
		},
		{
			Account:  "mgmt",
			Cluster:  "orca-main",
			Name:     "orca-main-v001",
			Region:   "us-west-2",
			Disabled: false,
			Instances: []instance{
				{
					ID:          "v001-1",
					HealthState: "UP",
				},
			},
		},
	}
	executionsStatus := map[string]int{
		"v000-1": 0,
		"v001-1": 1,
	}

	// when
	task := buildReapTask(orcaServerGroups, executionsStatus, testConfig)

	// then
	assert.Equal(t, task.Application, "orca")
	assert.NotNil(t, task.Job)

	assert.Len(t, task.Job, 1)
	j, ok := task.Job[0].(terminateInstancesJob)
	assert.True(t, ok)
	assert.Equal(t, j.Type, "terminateInstances")
	assert.Equal(t, j.Credentials, "mgmt")
	assert.Equal(t, j.Region, "us-west-2")
	assert.Contains(t, j.InstanceIDs, "v000-1")
}

func TestBuildReapTask_EmptyServerGroup(t *testing.T) {
	// given
	orcaServerGroups := []serverGroup{
		{
			Account:   "mgmt",
			Cluster:   "orca-main",
			Name:      "orca-main-v000",
			Region:    "us-west-2",
			Disabled:  true,
			Instances: []instance{},
		},
		{
			Account:  "mgmt",
			Cluster:  "orca-main",
			Name:     "orca-main-v001",
			Region:   "us-west-2",
			Disabled: false,
			Instances: []instance{
				{
					ID:          "v001-1",
					HealthState: "UP",
				},
			},
		},
	}
	executionsStatus := map[string]int{
		"v001-1": 1,
	}

	// when
	task := buildReapTask(orcaServerGroups, executionsStatus, testConfig)

	// then
	assert.Equal(t, task.Application, "orca")
	assert.NotNil(t, task.Job)

	assert.Len(t, task.Job, 1)
	j, ok := task.Job[0].(destroyServerGroupJob)
	assert.True(t, ok)
	assert.Equal(t, j.Type, "destroyServerGroup")
	assert.Equal(t, j.Credentials, "mgmt")
	assert.Equal(t, j.Region, "us-west-2")
	assert.Equal(t, j.CloudProvider, "aws")
	assert.Contains(t, j.ServerGroupName, "orca-main-v000")
}
