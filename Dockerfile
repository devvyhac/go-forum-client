FROM golang:1.25.0-alpine

WORKDIR /app

COPY . .

RUN go build -o /client .

CMD [ "/client" ]