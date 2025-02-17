FROM quay.io/cdis/golang:1.23-bookworm as build-deps

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR $GOPATH/src/github.com/uc-cdis/ssjdispatcher/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GITCOMMIT=$(git rev-parse HEAD) \
    GITVERSION=$(git describe --always --tags) \
    && go build \
    -ldflags="-X 'github.com/uc-cdis/ssjdispatcher/handlers/version.GitCommit=${GITCOMMIT}' -X 'github.com/uc-cdis/ssjdispatcher/handlers/version.GitVersion=${GITVERSION}'" \
    -o /ssjdispatcher

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
USER nobody
COPY --from=build-deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-deps /ssjdispatcher /ssjdispatcher
COPY --from=build-deps /etc_passwd /etc/passwd
CMD ["/ssjdispatcher"]
