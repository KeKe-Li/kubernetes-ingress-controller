FROM golang:1.13-alpine

RUN apk add --update \
    ca-certificates \
  && rm -rf /var/cache/apk/*

RUN echo "patent:x:65534:65534:Patent:/:" > /etc_passwd

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY="https://goproxy.io"

WORKDIR /go/src/github.com/KeKe-Li/kubernetes-ingress-controller
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go install -ldflags='-d -s -w' -tags netgo -installsuffix netgo -v ./...

FROM scratch

COPY --from=0 /go/bin/kubernetes-ingress-controller /bin/kubernetes-ingress-controller
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /etc_passwd /etc/passwd

CMD ["/bin/kubernetes-ingress-controller"]
