FROM golang:latest

RUN go version

ENV GOPATH=/
COPY ./ ./

RUN go mod download
RUN go build -o CompClubProject ./cmd/main.go

CMD ["./CompClubProject", "test_file.txt"]