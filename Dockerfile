FROM quay.io/cdis/golang:1.16-bullseye as build-deps

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR $GOPATH/src/github.com/uc-cdis/sower/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /ssjdispatcher

FROM scratch
COPY --from=build-deps /ssjdispatcher /ssjdispatcher
CMD ["/sower"]
