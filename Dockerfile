FROM golang:1.10-alpine as build-env

# Change when glide updates
ENV GLIDE_VER v0.13.1

# Install glide
RUN \
	apk add -q --no-cache --virtual .build-deps git && \
	go get github.com/Masterminds/glide                 && \
	cd ${GOPATH}/src/github.com/Masterminds/glide       && \
	TAG=$(git tag | grep ${GLIDE_VER} | sort | tail -1) && \
	git checkout -q ${TAG}                                 && \
	go install

# First copy glide dep files
# so that we don't reinstall our deps
# on each build
# changes in either of these files will trigger dep rebuild
WORKDIR $GOPATH/src/github.com/aurieh/ddg-ng/
COPY ./glide.yaml .
COPY ./glide.lock .

RUN glide --no-color install

# Copy over source
COPY . .

# Build and install the binary.
RUN go install github.com/aurieh/ddg-ng

# Make the minimal runtime container from the binary.
FROM golang:1.10-alpine

COPY --from=build-env /go/bin/ddg-ng /go/bin/
ENTRYPOINT ["ddg-ng"]
