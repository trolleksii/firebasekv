FROM golang as base
WORKDIR /go/src/firestorekv
COPY . /go/src/firestorekv
RUN go build

FROM scratch
COPY --from=base /go/src/firestorekv/firestorekv /usr/bin/firestorekv
CMD ["/usr/bin/firestorekv"]