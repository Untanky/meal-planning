FROM golang:1.22-alpine AS build

RUN apk add --update nodejs npm make g++

RUN npm install -G pnpm

ENV PNPM_HOME=/usr/local/bin
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .

RUN npx pnpm install

RUN make clean build

FROM alpine as release

WORKDIR /app

COPY --from=build /app/dist /app

VOLUME /data

EXPOSE 8080

ENTRYPOINT /app/meal-planner
