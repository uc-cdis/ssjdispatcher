FROM golang:1.14 as build-deps

# Build static binary
RUN mkdir -p /go/src/github.com/uc-cdis/ssjdispatcher
WORKDIR /go/src/github.com/uc-cdis/ssjdispatcher
ADD . .

RUN go build -ldflags "-linkmode external -extldflags -static" -o bin/ssjdispatcher

# Set up small scratch image, and copy necessary things over
# FROM scratch

# COPY --from=build-deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=build-deps /go/src/github.com/uc-cdis/ssjdispatcher/bin/ssjdispatcher /ssjdispatcher

ENTRYPOINT ["bin/ssjdispatcher"]
