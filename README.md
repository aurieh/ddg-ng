# DDG-NG
![](https://travis-ci.org/aurieh/ddg-ng.svg?branch=master)

A [DuckDuckGo](https://duckduckgo.com) bot for [Discord](https://discordapp.com).

# Building with Docker
`git clone https://github.com/aurieh/ddg-ng` then `docker build -t aurieh/ddg-ng .`.

To run, you can use several ways to provide the [required token](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token).
```sh
# Using a command line argument
docker run aurieh/ddg-ng --token=asdf1234

# Using an ENV var
docker run -e token=asdf1234 aurieh/ddg-ng

# Using a yaml settings file
echo 'token: asdf1234' > ddg.yaml
docker run -v $(pwd)/ddg.yaml:/etc/ddg/ddg.yaml aurieh/ddg-ng
```

# Building manually
`git clone https://github.com/aurieh/ddg-ng`, `glide install` and `go build`.

# License
[BSD 3-Clause](./LICENSE)
