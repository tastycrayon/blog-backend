FROM golang:latest

WORKDIR /usr/src/app

# to hot reload build
RUN go install github.com/cosmtrek/air@latest

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify


COPY . .

RUN go mod tidy
# RUN go build -v -o /usr/local/bin/app ./...

# CMD ["app"]