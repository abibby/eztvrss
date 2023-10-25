FROM golang:1.21 as go-build

WORKDIR /go/src/github.com/abibby/eztvrss

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /eztvrss

# Now copy it into our base image.
FROM alpine

RUN apk update && \
    apk add ca-certificates && \
    update-ca-certificates

# RUN apt-get update && apt-get install -y ca-certificates
# RUN update-ca-certificates

COPY --from=go-build /eztvrss /eztvrss

CMD ["/eztvrss"]
