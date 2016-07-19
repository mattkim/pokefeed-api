FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/pokefeed/pokefeed-api

ENV USER mattkim
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET Avxrhb9PI1uJTAb0

# Replace this with actual PostgreSQL DSN.
# TODO: need to hook this up to heroku env var
ENV DSN postgres://mattkim@localhost:5432/pokefeed-api?sslmode=disable

WORKDIR /go/src/github.com/pokefeed/pokefeed-api

RUN godep go build

EXPOSE 8888
CMD ./pokefeed-api
