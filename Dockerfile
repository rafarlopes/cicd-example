FROM golang:1.17.5-stretch as base

WORKDIR /cicd-example

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod ./
# My sample add doesn't have any external deps so the go.sum doesn't exist. In case you have, use the commented COPY command below
# COPY go.mod go.sum ./

ARG VERBOSE

# https://github.com/golang/go/issues/27719
# Better caching of dependencies between layers
RUN go mod graph | awk '$1 !~ /@/ { print $2 }' | xargs -r go get

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 go build $VERBOSE -o myservice .
#########################

FROM base as test
# Later used only for testing. Added in case we want to add something more in that step

#########################

FROM gcr.io/distroless/static-debian11 as final

# Copy both binary and source code for sentry.
COPY --from=base /cicd-example/ /cicd-example

EXPOSE 8080

CMD ["/cicd-example/myservice"]
