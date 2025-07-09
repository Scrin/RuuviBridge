FROM golang:1.24-alpine as build

ARG VERSION="unknown-docker"

WORKDIR /go/src/github.com/Scrin/RuuviBridge/
COPY . ./
RUN go install -v -ldflags "-X github.com/Scrin/RuuviBridge/common/version.Version=${VERSION}" ./cmd/ruuvibridge

FROM alpine

COPY --from=build /go/bin/ruuvibridge /usr/local/bin/ruuvibridge

USER 1337:1337

CMD ["ruuvibridge"]
