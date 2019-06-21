package probes

import (
	"os"
	"time"

	"github.com/cirocosta/slirunner/runnable"
)

func NewLogin(target, username, password, concourseUrl string) runnable.Runnable {
	var (
		config = Config{
			Target:       target,
			Username:     username,
			Password:     password,
			ConcourseUrl: concourseUrl,
		}
		timeout = 60 * time.Second
	)

	return runnable.NewWithLogging("login",
		runnable.NewWithMetrics("login",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`

	fly -t {{ .Target }} login -u {{ .Username }} -p {{ .Password }} -c {{ .ConcourseUrl }}

				`, config), os.Stderr),
				timeout,
			),
		),
	)
}

func NewCreateAndRunNewPipeline(target, prefix string) runnable.Runnable {
	var (
		config = Config{
			Target:   target,
			Pipeline: prefix + "create-and-run-new-pipeline",
		}
		timeout = 60 * time.Second
	)

	return runnable.NewWithLogging("create-and-run-new-pipeline",
		runnable.NewWithMetrics("create-and-run-new-pipeline",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`

	set -o errexit
	set -o xtrace

	fly -t {{ .Target }} destroy-pipeline -n -p {{ .Pipeline }} || true
	fly -t {{ .Target }} set-pipeline -n -p {{ .Pipeline }} -c <(echo '`+pipelineContents+`')
	fly -t {{ .Target }} unpause-pipeline -p {{ .Pipeline }}

	wait_for_build () {
		fly -t {{ .Target }} builds -j {{ .Pipeline }}/auto-triggering | \
			grep -v pending | \
			wc -l
	}

	until [ "$(wait_for_build)" -gt 0 ]; do
		echo 'waiting for job to automatically trigger...'
		sleep 1
	done

	fly -t {{ .Target }} watch -j {{ .Pipeline }}/auto-triggering
	fly -t {{ .Target }} destroy-pipeline -n -p {{ .Pipeline }}

				`, config), os.Stderr),
				timeout,
			),
		),
	)
}

func NewHijackFailingBuild(target, prefix string) runnable.Runnable {
	var (
		config = Config{
			Target:   target,
			Pipeline: prefix + "hijack-failing-build",
		}
		timeout = 60 * time.Second
	)

	return runnable.NewWithLogging("hijack-failing-build",
		runnable.NewWithMetrics("hijack-failing-build",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`

	set -o errexit
	set -o xtrace

	fly -t {{ .Target }} set-pipeline -n -p {{ .Pipeline }} -c <(echo '`+pipelineContents+`')
	fly -t {{ .Target }} unpause-pipeline -p {{ .Pipeline }}

	job_name={{ .Pipeline }}/failing
	fly -t {{ .Target }} trigger-job -j "$job_name" -w || true

	build=$(fly -t {{ .Target }} builds -j "$job_name" | head -1 | awk '{print $3}')
	fly -t {{ .Target }} hijack -j "$job_name" -b $build echo Hello World

				`, config), os.Stderr),
				timeout,
			),
		),
	)
}

func NewRunExistingPipeline(target, prefix string) runnable.Runnable {
	var (
		config = Config{
			Target:   target,
			Pipeline: prefix + "run-existing-pipeline",
		}
		timeout = 60 * time.Second
	)

	return runnable.NewWithLogging("run-existing-pipeline",
		runnable.NewWithMetrics("run-existing-pipeline",
			runnable.NewWithTimeout(
				runnable.NewShellCommand(FormatProbe(`

	set -o xtrace
	set -o errexit

	fly -t {{ .Target }} set-pipeline -n -p {{ .Pipeline }} -c <(echo '`+pipelineContents+`')
	fly -t {{ .Target }} unpause-pipeline -p {{ .Pipeline }}

	fly -t {{ .Target }} trigger-job -w -j "{{ .Pipeline }}/simple-job"

				`, config), os.Stderr),
				timeout,
			),
		),
	)
}

func NewAll(target, username, password, concourseUrl, prefix string) runnable.Runnable {
	return runnable.NewSequentially([]runnable.Runnable{

		NewLogin(target, username, password, concourseUrl),

		runnable.NewConcurrently([]runnable.Runnable{
			NewCreateAndRunNewPipeline(target, prefix),
			NewHijackFailingBuild(target, prefix),
			NewRunExistingPipeline(target, prefix),
		}),
	})

}
