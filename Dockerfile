FROM golang as build

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/isovalent/rebel-base
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

ENV \
  CGO_ENABLED=0 \
  GOOS=linux

# Install the package and create test binary
RUN go install -v ./... && \
    go test -c


FROM scratch
COPY --from=build /go/bin/rebel-base /rebel-base

# This container exposes port 8080 to the outside world
EXPOSE 8000

USER 1000

# Run the executable
ENTRYPOINT ["/rebel-base"]
