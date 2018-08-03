FROM golang

ADD favicon.ico main.go /go/src/github.com/justinclift/db4s_cluster_downloader/

# Only used for development on my local PC
#ADD DB.Browser.for.SQLite-3.10.1-win64.exe DB.Browser.for.SQLite-3.10.1-win32.exe cert.pem key.pem /data/

RUN go install github.com/justinclift/db4s_cluster_downloader

ENTRYPOINT /go/bin/db4s_cluster_downloader

EXPOSE 443
