FROM golang
RUN apt-get update \
    && apt-get install -y curl git zip unzip bzip2 xvfb default-jre default-jdk \
    && rm -rf /var/lib/apt/lists/*
RUN go get -u github.com/remiflavien1/recon-archy \
    && cd /go/src/github.com/remiflavien1/recon-archy \
    && go get -u ./... \
    && cd /go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor \
    && go run init.go --alsologtostderr  --download_browsers --download_latest
