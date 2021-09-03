FROM golang as base
WORKDIR /go/src/firestorekv
COPY . /go/src/firestorekv
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o firestorekv

FROM scratch 
COPY --from=base /go/src/firestorekv/firestorekv /usr/local/bin/
CMD ["/usr/local/bin/firestorekv"]
