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
		-o /usr/bin/slirunner \
		-ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\""

FROM base AS test

	RUN go test -v ./...


FROM ubuntu AS release

	RUN apt update && apt install -y ca-certificates

	COPY \
		--from=build \
		/usr/bin/slirunner \
		/usr/bin/slirunner

	ENTRYPOINT [ "/usr/bin/slirunner" ]
