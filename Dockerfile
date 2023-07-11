FROM golang:1.20.4-alpine3.18

LABEL base.name='onnuribackend'

WORKDIR /app

COPY . .

# Build the Go app
RUN go build -o main .

EXPOSE 80

ENTRYPOINT [ "./main" ]