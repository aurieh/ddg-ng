FROM golang:1.16-alpine as build-env

# First copy glide dep files
# so that we don't reinstall our deps
# on each build
# changes in either of these files will trigger dep rebuild
WORKDIR $GOPATH/src/github.com/aurieh/ddg-ng/

COPY . .

RUN go install ./...

# Make the minimal runtime container from the binary.
FROM golang:1.16-alpine

RUN adduser \
	    --home=/var/empty/ \
	    --shell=/sbin/nologin \
	    --disabled-password \
	    --no-create-home \
	    ddg-ng
USER ddg-ng

COPY --from=build-env /go/bin/ddg-ng /go/bin/
ENTRYPOINT ["ddg-ng"]
