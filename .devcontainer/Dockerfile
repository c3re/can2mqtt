FROM debian:latest AS base

ARG GO_VERSION="<VERSION>"
ENV GO_VERSION=${GO_VERSION}

RUN apt-get update \
	&& apt-get upgrade -y
RUN apt-get install -y ca-certificates openssl wget git git-lfs tar gnupg2 gcc

RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
	&& tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH /usr/local/go/bin:$PATH

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"
WORKDIR $GOPATH
RUN go install -v golang.org/x/tools/gopls@latest \
	&& go install github.com/go-delve/delve/cmd/dlv@latest \
	&& go install honnef.co/go/tools/cmd/staticcheck@latest

FROM base AS dev

RUN go install github.com/pact-foundation/pact-go/v2@v2.0.2 \
	&& go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest \
	&& go install github.com/google/go-licenses@latest \
	&& go install golang.org/x/vuln/cmd/govulncheck@latest
