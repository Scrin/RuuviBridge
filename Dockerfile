FROM golang:1.17-alpine as build

RUN apk add git

WORKDIR /go/src/github.com/Scrin/RuuviBridge/
COPY . ./
RUN go install -v -ldflags "-X github.com/Scrin/RuuviBridge/common/version.Version=git-$(git rev-list -1 HEAD 2>/dev/null || echo unknown)" ./cmd/ruuvibridge

FROM alpine

COPY --from=build /go/bin/ruuvibridge /usr/local/bin/ruuvibridge

USER 1337:1337

CMD ["ruuvibridge"]
