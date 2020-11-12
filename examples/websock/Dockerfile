FROM golang:1.13-alpine3.11 as builder
RUN apk add git 

RUN mkdir /build
ADD *.go /build/

WORKDIR /build
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -a -o webtty .


# generate clean, final image for end users
FROM alpine:3.11.3
COPY --from=builder /build/webtty .

# executable
ENTRYPOINT [ "./webtty" ]