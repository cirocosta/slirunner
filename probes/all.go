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

	fly -t ci destroy-pipeline -n -p new-pipeline
	fly -t ci set-pipeline -n -p new-pipeline -c <("`+samplePipeline+`")
	fly -t ci unpause-pipeline -p new-pipeline

	until [ "$(fly -t ci builds -j new-pipeline/auto-triggering | grep -v pending | wc -l)" -gt 0 ]; do
		echo 'waiting for job to trigger...'
		sleep 1
	done
	fly -t ci watch -j new-pipeline/auto-triggering
	fly -t ci destroy-pipeline -n -p new-pipeline
`, os.Stderr), 60*time.Second)
