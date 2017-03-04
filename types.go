package orcareaper

type serverGroup struct {
	Account   string     `json:"account"`
	Cluster   string     `json:"cluster"`
	Name      string     `json:"name"`
	Region    string     `json:"region"`
	Disabled  bool       `json:"isDisabled"`
	Instances []instance `json:"instances"`
}

type instance struct {
	ID          string `json:"id"`
	HealthState string `json:"healthState"`
}

type task struct {
	Application string        `json:"application"`
	Description string        `json:"description"`
	Job         []interface{} `json:"job"`
}

type job struct {
	Type        string `json:"type"`
	Region      string `json:"region"`
	Credentials string `json:"credentials"`
}

type terminateInstancesJob struct {
	job

	InstanceIDs []string `json:"instanceIds"`
}

type destroyServerGroupJob struct {
	job

	ServerGroupName string `json:"serverGroupName"`
	CloudProvider   string `json:"cloudProvider"`
}

type orcaInstance struct {
	Overdue    bool                `json:"overdue"`
	Count      int                 `json:"count"`
	Executions orcaExecutionDetail `json:"executions"`
}

type orcaExecutionDetail struct {
	Application string `json:"application"`
	URL         string `json:"url"`
}
