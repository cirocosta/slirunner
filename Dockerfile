FROM golang AS golang


FROM golang AS base

	ENV CGO_ENABLED=0
	RUN apt update && apt install -y git

	ADD . /src
	WORKDIR /src

	RUN go mod download


FROM base AS build

	RUN go build \
		-tags netgo -v -a \
		-o /usr/local/bin/slirunner \
		-ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\""

FROM base AS test

	RUN go test -v ./...


FROM alpine AS concourse-release

	ARG CONCOURSE_VERSION=6.4.0

	ADD https://github.com/concourse/concourse/releases/download/v${CONCOURSE_VERSION}/concourse-${CONCOURSE_VERSION}-linux-amd64.tgz /tmp
	RUN tar xvzf /tmp/concourse-${CONCOURSE_VERSION}-linux-amd64.tgz -C /usr/local
	RUN tar xvzf /usr/local/concourse/fly-assets/fly-linux-amd64.tgz -C /usr/local/bin/


FROM ubuntu AS release

	RUN apt update && apt install -y ca-certificates

	COPY \
		--from=concourse-release \
		/usr/local/bin/fly \
		/usr/local/bin/fly

	COPY \
		--from=build \
		/usr/local/bin/slirunner \
		/usr/local/bin/slirunner

	ENTRYPOINT [ "/usr/local/bin/slirunner" ]
