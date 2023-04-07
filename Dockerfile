FROM golang:1.20

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o tgbot ./cmd/main.go

CMD [ "./tgbot" ]