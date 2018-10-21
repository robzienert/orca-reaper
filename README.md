# orca-reaper

_**This project is not necessary anymore. Orca's architecture is fundamentally different now**_

A script to that can be used to gracefully terminate Spinnaker Orca instances.

Spinnaker's orchestration service, Orca, currently uses Spring Batch which
makes it a stateful service that will lose track of executions in the event of
a node loss or process restart. This can make deployment operations difficult if
your organization doesn't have a reaper in place.

This script is an open source adaptation of the reaper we use at Netflix to
perform red/black deployments for Orca. Internally, we check for the number of
execution invocations on a per-instance basis, and once an Orca instance has 0
invocations, we know it has been drained and can be safely terminated. 

Your Orca deployment should be performed via red/black, keeping disabled Orca
servers until they are terminated by this script. We run our reaper job every
60 minutes, but your milage may vary.

# install

You can either run `build.sh` or download the latest release from the Releases
page.

# usage

```
Usage of orca-reaper:
  -apiBaseURL string
    	Base URL for Gate
  -cluster string
    	Limit scope of reaping to a cluster (useful for multiple Spinnaker deploys)
  -credentials string
    	Account Spinnaker is running in
  -dryRun
    	When set, no reap tasks will be submitted to Orca
  -region string
    	Region where Spinnaker is running (default "us-west-2")
  -x509CertPath string
    	x509 certificate path
  -x509KeyPath string
    	x509 key path
```

# example

In this example, the reaper job is only reaping Orca instances in the staging
environment.

```
$ go run cmd/orcareaper/main.go \
  -apiBaseURL https://api-staging.spinnaker.domain.net \
  -x509CertPath ~/.spinnaker/rzienert.pem.crt \
  -x509KeyPath ~/.spinnaker/rzienert.pem.key \
  -credentials mgmt \
  -cluster "orca-staging" \
  -dryRun
2017/02/10 17:18:38 {
  "application":"orca",
  "description":"Terminate instances and server groups",
  "job":[{
    "type":"terminateInstances",
    "region":"us-west-2",
    "credentials":"mgmt",
    "instanceIds":["i-1234"]
  }]
}
2017/02/10 17:18:38 dry-run mode enabled, skipping action
```
