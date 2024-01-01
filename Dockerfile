# Build the application from source
ARG BUILDARCH
FROM --platform=$BUILDPLATFORM golang:1.21 AS build-stage

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /meal-planning ./cmd

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /meal-planning /meal-planning
COPY index.html /index.html
COPY days.template.html /days.template.html

EXPOSE 8080

ENTRYPOINT ["/meal-planning"]
