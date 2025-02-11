FROM golang:1.21.0-alpine AS builder

COPY . /github.com/Snake1-1eyes/Test_Ozon
WORKDIR /github.com/Snake1-1eyes/Test_Ozon
RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go

FROM scratch AS runner
WORKDIR /YouNote/

COPY --from=builder /github.com/Snake1-1eyes/Test_Ozon/.bin .
COPY --from=builder /github.com/Snake1-1eyes/Test_Ozon/.env ./

EXPOSE 8080

ENTRYPOINT ["./.bin"]