FROM golang AS build
ADD . /go/src/github.com/slowteetoe/tidechecker
WORKDIR /go/src/github.com/slowteetoe/tidechecker
RUN go get -d ./... && CGO_ENABLED=0 go build -o tidechecker .

FROM alpine
WORKDIR /app
COPY --from=build /go/src/github.com/slowteetoe/tidechecker/tidechecker /app
COPY ./data /app/data
EXPOSE 10000 
ENTRYPOINT ["/app/tidechecker"]
