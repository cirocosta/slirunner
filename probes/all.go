package probes

import (
	"os"
	"time"

	"github.com/cirocosta/slirunner/runnable"
)

const samplePipeline = `
resources:
- name: time-trigger
  type: time
  source: {interval: 24h}

jobs:
- name: simple-job
  build_logs_to_retain: 20
  public: true
  plan:
  - &say-hello
    task: say-hello
    config:
      platform: linux
      image_resource:
        type: registry-image
        source: {repository: busybox}
      run:
        path: echo
        args: ["Hello, world!"]

- name: failing
  build_logs_to_retain: 20
  public: true
  plan:
  - task: fail
    config:
      platform: linux
      image_resource:
        type: registry-image
        source: {repository: busybox}
      run:
        path: false

- name: auto-triggering
  build_logs_to_retain: 20
  public: true
  plan:
  - get: time-trigger
    trigger: true
  - *say-hello
`

var CreateAndRunNewPipeline = runnable.NewWithTimeout(runnable.NewShellCommand(`
	set -o errexit

	fly -t local destroy-pipeline -n -p new-pipeline
	fly -t local set-pipeline -n -p new-pipeline -c <("`+samplePipeline+`")
	fly -t local unpause-pipeline -p new-pipeline

	until [ "$(fly -t local builds -j new-pipeline/auto-triggering | grep -v pending | wc -l)" -gt 0 ]; do
		echo 'waiting for job to trigger...'
		sleep 1
	done
	fly -t local watch -j new-pipeline/auto-triggering
	fly -t local destroy-pipeline -n -p new-pipeline
`, os.Stderr), 60*time.Second)

var All = runnable.NewConcurrently([]runnable.Runnable{
	CreateAndRunNewPipeline,
})
