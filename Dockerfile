FROM golang:alpine

WORKDIR /anti-bruteforce

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN go build -o /anti-bruteforce/build/anti_bruteforce_app/anti_bruteforce ./cmd/server

EXPOSE 8080

ENTRYPOINT [ "/anti-bruteforce/build/anti_bruteforce_app/anti_bruteforce" ]