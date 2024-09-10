FROM ubuntu

WORKDIR /app
COPY . .

RUN apt update
RUN  apt-get update &&  apt-get -y install golang-go 

RUN apt install libopus0 libopus-dev
RUN go mod tidy
RUN  go build -o main cmd/main.go 

# FROM ubuntu:20.04

# COPY --from=GOLANG /app .

# COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD [ "/main" ]