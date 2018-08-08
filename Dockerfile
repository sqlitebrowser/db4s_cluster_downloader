FROM golang

ADD favicon.ico main.go template.html /go/src/github.com/justinclift/db4s_cluster_downloader/

RUN go get github.com/BurntSushi/toml
RUN go get github.com/jackc/pgx
RUN go get github.com/jackc/pgx/pgtype
RUN go install github.com/justinclift/db4s_cluster_downloader

ENTRYPOINT /go/bin/db4s_cluster_downloader

EXPOSE 443
