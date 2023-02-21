FROM cgr.dev/chainguard/go:1.20.1@sha256:402786fb655a09632e3e0d468e27001339e59fd68086cb9aa5b5b4388773aaa8 as build

WORKDIR /work

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd /work/cmd/
COPY pkg /work/pkg

RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o /miflora-go /work/cmd/miflora-go/

FROM cgr.dev/chainguard/static:latest@sha256:9a2320c5820ba0e75cc3a84397dd02b0e36787f24bcdd4a36bb9af5c3a37ec7e

COPY --from=build /miflora-go /miflora-go

CMD ["./miflora-go"]