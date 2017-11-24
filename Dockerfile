FROM golang:alpine


RUN \
	apk add --no-cache --virtual .build-deps git && \
	go-wrapper download github.com/yhat/scrape && \
	go-wrapper download github.com/lib/pq && \
	go-wrapper download github.com/spf13/viper && \
	go-wrapper download github.com/spf13/pflag && \
	go-wrapper download github.com/Sirupsen/logrus && \
	go-wrapper download github.com/bwmarrin/discordgo && \
	go-wrapper download github.com/dustin/go-humanize && \
	go-wrapper download github.com/olekukonko/tablewriter

COPY . $GOPATH/src/github.com/aurieh/ddg-ng/

RUN \
	go-wrapper install github.com/aurieh/ddg-ng && \
	apk del .build-deps

WORKDIR $GOPATH/src/github.com/aurieh/ddg-ng/
ENTRYPOINT ["go-wrapper", "run"]
