FROM golang:1.14.4 AS builder
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN go get github.com/sirupsen/logrus && go get github.com/mitchellh/mapstructure && go get github.com/Akagi201/utilgo/jobber && go get github.com/Akagi201/utilgo/conflag && go get github.com/jessevdk/go-flags && go get github.com/Shopify/go-lua
WORKDIR /go/src/esalert
ADD . ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o esalert cmd/esalert/main.go
############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /go/src/esalert/esalert /go/bin/esalert
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/localtime /etc/localtime
ENTRYPOINT ["/go/bin/esalert"]