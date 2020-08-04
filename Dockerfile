FROM golang
RUN apt-get update \
    && apt-get install -y curl git xvfb openjdk-11-jre firefox \
    && rm -rf /var/lib/apt/lists/*
RUN go get github.com/remiflavien1/recon-archy \
    && cd $GOPATH/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/ \
    && go run init.go --alsologtostderr  --download_browsers --download_latest

