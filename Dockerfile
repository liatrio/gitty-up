FROM golang:1.13-alpine AS build

RUN apk add git build-base

RUN addgroup -g 1000 jenkins && adduser -h /home/jenkins -G jenkins -u 1000 -D jenkins
USER jenkins

RUN go get gopkg.in/src-d/go-git.v4/... github.com/hashicorp/hcl golang.org/x/oauth2 github.com/google/go-github/github github.com/stretchr/testify/assert

COPY --chown=jenkins:jenkins ./*.go /go/src/gitops/
WORKDIR /go

RUN go test gitops -cover -v
RUN go install gitops

FROM alpine

RUN addgroup -g 1000 jenkins && adduser -h /home/jenkins -G jenkins -u 1000 -D jenkins
USER jenkins

COPY --from=build /go/bin/gitops /gitops

WORKDIR /home/jenkins

ENTRYPOINT [ "/gitops" ]
