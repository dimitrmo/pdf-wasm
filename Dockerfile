# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS build

WORKDIR /usr/local

ADD go.mod .
ADD go.sum .
ADD font.ttf .
ADD main.go .
ADD czech.png .

RUN GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o pdf.wasm

FROM scratch AS wasm
COPY --from=build /usr/local/pdf.wasm /
