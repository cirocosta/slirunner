package exporter

import (
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	Address string `long:"prometheus-address" default:":9001" description:"address to listen for prometheus scrapes"`
	Path    string `long:"prometheus-path"    default:"/"     description:"path to serve prometheus metrics"`

	listener net.Listener
}

// Listen initiates the HTTP server using the configurations
// provided via ExporterConfig.
//
// This is a blocking method - make sure you either make use of
// goroutines to not block if needed.
func (e *Exporter) Listen() (err error) {
	http.Handle(e.Path, promhttp.Handler())

	e.listener, err = net.Listen("tcp", e.Address)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on address %s", e.Address)
		return
	}

	err = http.Serve(e.listener, nil)
	if err != nil {
		err = errors.Wrapf(err,
			"failed listening on address %s",
			e.Address)
		return
	}

	return
}

// Stop closes the tcp listener (if exists).
func (e *Exporter) Close() (err error) {
	if e.listener == nil {
		return
	}

	err = e.listener.Close()
	return
}
