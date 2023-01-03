# リリース用のビルドを行うコンテナイメージを作成するステージ
FROM golang:1.18.2-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# ---------------------------------------------------
# リリース用のコンテナイメージを作成するステージ

FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# ---------------------------------------------------
# 開発時に利用するコンテナイメージを作成するステージ

FROM golang:latest as dev
RUN apt update
# Module ModeをONにする
ENV GO111MODULE on
WORKDIR /app
RUN go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \ 
    go install -v github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/ramya-rao-a/go-outline@latest && \
    go install github.com/stamblerre/gocode@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install github.com/cosmtrek/air@latest
CMD ["air"]
