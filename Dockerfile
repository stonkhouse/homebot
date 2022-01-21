#TODO: Make the image leaner
FROM golang:1.17-alpine
#FROM scratch
#FROM debian
WORKDIR /go/stonkhouse/src/
COPY . .

# -d separates getting and installing
# -v to make it verbose
RUN go get -d -v ./...
RUN go mod tidy
RUN go mod vendor
RUN go install -v ./...
RUN go build

RUN mv ./config.yml ../
RUN mv ./stonkbot ../

WORKDIR /go/stonkhouse
RUN rm -rf ./src
# Application Entrypoint
ENTRYPOINT ["./stonkbot"]
#ENTRYPOINT ["echo 'hello'"]