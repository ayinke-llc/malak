FROM golang:1.23 as build-env
WORKDIR /go/malak

LABEL org.opencontainers.image.description "Open source Investors' relationship hub for Founders"

COPY ./go.mod /go/malak
COPY ./go.sum /go/malak

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
RUN go mod verify
# COPY the source code as the last step
COPY . .

RUN CGO_ENABLED=0
RUN go install ./cmd

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/cmd /
CMD ["/cmd http"]

