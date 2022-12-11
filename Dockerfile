FROM golang:latest
RUN apt update
RUN apt install vim tree -y
WORKDIR /code

