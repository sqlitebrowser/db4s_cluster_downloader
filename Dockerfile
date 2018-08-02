FROM golang

# This should grab the real cert and key, which shouldn't be stored in the repo
ADD . /go/src/github.com/justinclift/db4s_cluster_downloader

RUN go install github.com/justinclift/db4s_cluster_downloader

ENTRYPOINT /go/bin/db4s_cluster_downloader

EXPOSE 10443
