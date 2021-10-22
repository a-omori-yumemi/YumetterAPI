FROM golang:1.17 AS init

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.17 AS dev
WORKDIR /go/src/app
COPY --from=init /go /go
COPY . .
RUN go get github.com/google/wire/cmd/wire
RUN wire
CMD [ "go", "run", "main.go", "wire_gen.go" ]


FROM golang:1.17 AS build
WORKDIR /go/src/app
COPY --from=init /go /go
COPY . .
RUN go install -v
RUN go get github.com/google/wire/cmd/wire
RUN wire
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a

FROM scratch AS prod
COPY --from=build /go/src/app/YumetterAPI /bin/YumetterAPI
CMD ["YumetterAPI"]
