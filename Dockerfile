FROM golang:1.21 as build-env
WORKDIR /go/src/github.com/adelowo/sdump

COPY ./go.mod /go/src/github.com/adelowo/sdump
COPY ./go.sum /go/src/github.com/adelowo/sdump

RUN go mod download && go mod verify
COPY . .

ENV CGO_ENABLED=0 
RUN go install ./cmd

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/cmd /
CMD ["/cmd"]

