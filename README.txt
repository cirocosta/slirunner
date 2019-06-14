SLI RUNNER

	runs your Concourse SLI probes


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


		--> iniates a prometheus server exposing how the execution is going


// Runnable represents something that has to be executed and
// cancellable through context cancellation.
//
type Runnable interface {
	func run(ctx) -> bool
}

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



// ShellCommandRunnable runs a command in `/bin/bash` with the possibility of timing it
// aut.
//
type ShellCommandRunnable struct {
	Stdout io.Writer
	Stderr io.Writer

	func run(ctx) -> bool
}


