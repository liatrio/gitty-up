FROM golang:1.13-alpine AS build

RUN apk add git build-base

RUN addgroup -g 1000 jenkins && adduser -h /home/jenkins -G jenkins -u 1000 -D jenkins
USER jenkins

RUN go get gopkg.in/src-d/go-git.v4/... github.com/hashicorp/hcl golang.org/x/oauth2 github.com/google/go-github/github github.com/stretchr/testify/assert gopkg.in/yaml.v3

COPY --chown=jenkins:jenkins ./sample /go/src/gitty-up/sample/
COPY --chown=jenkins:jenkins ./*.go /go/src/gitty-up/
WORKDIR /go/src/gitty-up

RUN go test gitty-up -cover -v
RUN go install gitty-up

FROM alpine

RUN addgroup -g 1000 jenkins && adduser -h /home/jenkins -G jenkins -u 1000 -D jenkins
USER jenkins

COPY --from=build /go/bin/gitty-up /usr/bin/gitty-up

WORKDIR /home/jenkins

ENTRYPOINT [ "/usr/bin/gitty-up" ]
