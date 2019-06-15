package probes

import (
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/cirocosta/slirunner/runnable"
)

func New(target, prefix string) runnable.Runnable {
	logger := lager.NewLogger("probes")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.INFO))

	var createAndRunNewPipeline = runnable.NewWithLogging(logger.Session("create-and-run-new-pipeline"),
		runnable.NewWithMetrics("create-and-run-new-pipeline",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`

				set -o errexit
				set -o xtrace

				fly -t {{ .Target }} destroy-pipeline -n -p {{ .Pipeline }} || true
				fly -t {{ .Target }} set-pipeline -n -p {{ .Pipeline }} -c <(echo '`+pipelineContents+`')
				fly -t {{ .Target }} unpause-pipeline -p {{ .Pipeline }}

				wait_for_build () {
					fly -t local builds -j {{ .Pipeline }}/auto-triggering | \
						grep -v pending | \
						wc -l
				}

				until [ "$(wait_for_build)" -gt 0 ]; do
					echo 'waiting for job to automatically trigger...'
					sleep 1
				done

				fly -t local watch -j {{ .Pipeline }}/auto-triggering
				fly -t local destroy-pipeline -n -p {{ .Pipeline }}

				`, Config{Target: target, Pipeline: prefix + "create-and-run-new-pipeline"}), os.Stderr),
				60*time.Second,
			),
		),
	)

	var hijackFailingBuild = runnable.NewWithLogging(logger.Session("hijack-failing-build"),
		runnable.NewWithMetrics("hijack-failing-build",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`

				set -o errexit
				set -o xtrace

				fly -t {{ .Target }} set-pipeline -n -p {{ .Pipeline }} -c <(echo '`+pipelineContents+`')
				fly -t {{ .Target }} unpause-pipeline -p {{ .Pipeline }}

				job_name={{ .Pipeline }}/failing
				fly -t local trigger-job -j "$job_name" -w || true

				build=$(fly -t local builds -j "$job_name" | head -1 | awk '{print $3}')
				fly -t local hijack -j "$job_name" -b $build echo Hello World

				`, Config{Target: target, Pipeline: prefix + "hijack-failing-build"}), os.Stderr),
				60*time.Second,
			),
		),
	)

	var runExistingPipeline = runnable.NewWithLogging(logger.Session("run-existing-pipeline"),
		runnable.NewWithMetrics("run-existing-pipeline",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`
				set -o xtrace
				set -o errexit

				fly -t {{ .Target }} set-pipeline -n -p {{ .Pipeline }} -c <(echo '`+pipelineContents+`')
				fly -t {{ .Target }} unpause-pipeline -p {{ .Pipeline }}

				fly -t {{ .Target }} trigger-job -w -j "{{ .Pipeline }}/simple-job"

				`, Config{Target: target, Pipeline: prefix + "run-existing-pipeline"}), os.Stderr),
				60*time.Second,
			),
		),
	)

	return runnable.NewConcurrently([]runnable.Runnable{
		createAndRunNewPipeline,
		hijackFailingBuild,
		runExistingPipeline,
	})
}
