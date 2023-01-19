FROM golang:latest
RUN apt update
RUN apt install vim tree -y
WORKDIR /code

# EXPOSE 3000

CMD ["sh", "run.sh"]
