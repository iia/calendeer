FROM golang:1.16.0-buster

USER root
WORKDIR /root
COPY ./src/ ./calendeer
WORKDIR /root/calendeer/
RUN go get .
RUN go mod tidy
RUN go install .
CMD calendeer
