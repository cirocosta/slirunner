FROM golang:alpine AS golang


FROM golang AS base

	ENV CGO_ENABLED=0
	RUN apk add --update git gcc musl-dev

	ADD . /src
	WORKDIR /src

	RUN go mod download


FROM base AS build

	RUN go build \
		-tags netgo -v -a \
		-o /usr/bin/slirunner \
		-ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\""

FROM base AS test

	RUN go test -v ./...


FROM alpine AS release

	COPY \
		--from=build \
		/usr/bin/slirunner \
		/usr/bin/slirunner

	ENTRYPOINT [ "/usr/bin/slirunner" ]
