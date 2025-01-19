FROM golang:1.22.5-alpine3.20
WORKDIR /app


COPY . .

# download dependencies
RUN go get -d -v ./...

RUN go build -o main .

EXPOSE 8000

CMD [ "./main" ]