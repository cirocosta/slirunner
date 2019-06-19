
# sli runner

> probes your Concourse installation, generating Service Level Indicators (SLIs)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [why](#why)
- [prior](#prior)
- [what](#what)
- [license](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## why

Even tough Concourse emits many metrics that are useful for an operator,
it might still be hard to have a quick grasp of how high-level 
funcitionality is performing.

With SLIs (see [1]), one is able to better reason about what's broken
on user-facing functionality that the service exposes.


[1]: https://landing.google.com/sre/sre-book/chapters/service-level-objectives/
	

## prior

Before `slirunner`, [oxygen-mask][oxygen-mask] was the solution for running 
high-level probes against Concourse installation.

It has a few quirks that I don't think are necessary to have:

- requires another Concourse installation to run those probes
- tightly coupled to datadog
- performs some basic UI testing


[oxygen-mask]: https://github.com/concourse/oxygen-mask



## what

`slirunner` is a single Go binary that, once "installed" (run 
somewhere), periodically executes several probes against Concourse,
keeping track of the success and failures.

A consumer of `slirunner` can consume the reports from two mediums:

- Prometheus exposed metrics
- structured logs


It also supports:

- single runs
- worker-related probing against multiple tags and teams

## license

See [./LICENSE](./LICENSE)
