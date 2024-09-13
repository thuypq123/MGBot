FROM golang:1.23.1-alpine

WORKDIR /app
COPY . .

# RUN apt update
# RUN  apt-get update &&  apt-get -y install golang-go 
# RUN apk add --no-cache bash

# RUN apt-get update && apt-get install -y bash
# Update CA certificates after adding the new one
# RUN update-ca-certificates
# RUN cp local-ca.crt /usr/local/share/ca-certificates
RUN apk add --no-cache ca-certificates
RUN apk update 
RUN apk add --no-cache opus opus-dev alsa-lib alsa-lib-dev 
RUN apk add --no-cache ffmpeg
RUN apk add --no-cache bash
RUN go clean -modcache  
RUN go mod download 
RUN go mod tidy
RUN go mod verify
# Enable CGO and set necessary environment variables for cross-compilation
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Install build dependencies for gopus
RUN apk add --no-cache build-base
RUN go build -ldflags "-w" -o main cmd/main.go

# FROM ubuntu

# COPY --from=GOLANG /app .

# COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD [ "./main" ]