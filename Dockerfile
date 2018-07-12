# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev

ENV SRC_DIR=/go/src/github.com/windmilleng/blorg-backend

# Add the source code:
# (from current dir, add all files to dockerspace: /go/src...)
ADD . $SRC_DIR

# Build it:
RUN cd $SRC_DIR; go build -o server; cp server /app/
ENTRYPOINT ["./server"]
