FROM golang:1.22.5

WORKDIR /app

RUN apt-get update && apt-get upgrade -y \
    apt-get install -y \
    curl \
    gnupg2 \
    vim

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

RUN go install github.com/air-verse/air@latest
CMD ["air"]