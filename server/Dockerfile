FROM golang:1.22 as builder

RUN apt-get update

RUN apt-get -y install zip

RUN mkdir /opt/protoc && \
    cd /opt/protoc && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v27.2/protoc-27.2-linux-x86_64.zip && \
    unzip protoc-27.2-linux-x86_64.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

ENV PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /app

COPY ./server/go.mod ./server/go.sum ./server/
COPY ./lib/go.mod ./lib/

RUN cd ./lib && go mod download
RUN cd ./server && go mod download

COPY ./lib ./lib
COPY ./server ./server

RUN cd ./server/cmd/server && CGO_ENABLED=0 go build -o /app/server/chat-server

FROM scratch
COPY --from=builder /app/server/chat-server /opt/app/
EXPOSE 9876

CMD ["/opt/app/chat-server"]