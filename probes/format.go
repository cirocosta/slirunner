package probes

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type Config struct {
	Username     string
	Password     string
	ConcourseUrl string

	Target           string
	ExistingPipeline string
	Pipeline         string
}

func FormatProbe(formatting string, c Config) (res string) {
	var buf = new(bytes.Buffer)

	tmpl, err := template.New("").Parse(formatting)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse template"))
	}

	err = tmpl.Execute(buf, c)
	if err != nil {
		panic(errors.Wrapf(err, "failed to execute template"))
	}

	res = buf.String()
	return
}
