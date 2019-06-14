SLI RUNNER

	runs your Concourse SLI probes


WHY

	Even tough Concourse emits many metrics that are useful for an operator,
	it might still be hard to have a quick grasp of how high-level 
	funcitionality is performing.

	With SLIs (see [1]), one is able to better reason about what's broken
	on user-facing functionality that the service exposes.


	[1]: https://landing.google.com/sre/sre-book/chapters/service-level-objectives/
	

PRIOR

	Before `slirunner`, oxygenmask[1] was the solution for running 
	high-level probes against Concourse installation.

	It has few problems though:

		- requires another Concourse installation to run those probes
		  from

		- tightly coupled to datadog


	[1]: https://github.com/concourse/oxygen-mask



WHAT

	`slirunner` is a single Go binary that, once "installed" (run 
	somewhere), periodically executes several probes against Concourse,
	keeping track of the success and failures.

	A consumer of `slirunner` can consume the reports from two mediums:

		- Prometheus exposed metrics
		- structured logs

	It also supports:

		- single runs
		- worker-related probing against multiple tags and teams


USAGE

	`slirunner $fly_login_flags`
		|
		*--> 	1. authenticates against the target `concourse-url`

			2. sets the "existing pipeline"

			3. initiates periodic timeout'able execution

				- fly hijack failing build

				- run existing pipeline

				- view public pipeline

				- view build history 

				- create and run new pipeline


		--> 	iniates a prometheus server exposing how the execution 
			is going


	FAILURES

		task failures have the content of their failure 
		reported to `stderr` with a bit of wrapping around, but
		not structured as `json` or other machine-readable format.




// ConcurrentRunner concurrently runs a set of runnables, possibly
// cancelling them all at once.
//
// Implements the `Runnable` interface.
// 
type Concurrently struct {
	runnables []Runnable
	func run(ctx) -> bool
}


// WithTimeout wraps a Runnable with a timeout, ensuring that it doesn't
// run forever.
//
type WithTimeout struct {
	runnable Runnable
	func run(ctx) -> bool
}


// Periodically runs a set of registered Runnables forever (until cancelled).
//
type Periodically struct {
	runnables []Runnable
	func run(ctx) -> bool
}

