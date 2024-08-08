FROM golang:1.22.4 as build-env
WORKDIR /go

COPY ./go.mod /go/
COPY ./go.sum /go/

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
RUN go mod verify
# COPY the source code as the last step
COPY . .

RUN CGO_ENABLED=0
RUN go install ./cmd

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/cmd /
CMD ["/cmd"]

