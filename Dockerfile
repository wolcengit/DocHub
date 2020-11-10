#FROM wolcen/golang-compile:1.10.3-alpine3.7 AS build
#ADD . /go/src/github.com/TruthHun/DocHub
#WORKDIR /go/src/github.com/TruthHun/DocHub
#RUN	CGO_ENABLE=1 go build -v -a -o DocHub_linux_amd64 -ldflags="-w -s -X main.VERSION=$TAG -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'"

FROM truthhun/dochub:env
#COPY --from=build /go/src/github.com/TruthHun/DocHub/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
#COPY --from=build /go/src/github.com/TruthHun/DocHub/LICENSE /www/dochub/LICENSE
#COPY --from=build /go/src/github.com/TruthHun/DocHub/conf/app.conf.example /www/dochub/conf/app.conf.example
#COPY --from=build /go/src/github.com/TruthHun/DocHub/dictionary /www/dochub/dictionary
#COPY --from=build /go/src/github.com/TruthHun/DocHub/static /www/dochub/static
#COPY --from=build /go/src/github.com/TruthHun/DocHub/views /www/dochub/views
#COPY --from=build /go/src/github.com/TruthHun/DocHub/DocHub_linux_amd64 /www/dochub/DocHub_linux_amd64

COPY zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY LICENSE /www/dochub/LICENSE
COPY conf/app.conf.example /www/dochub/conf/app.conf.example
COPY dictionary /www/dochub/dictionary
COPY static /www/dochub/static
COPY views /www/dochub/views
COPY DocHub_linux_amd64 /www/dochub/DocHub_linux_amd64
WORKDIR /www/dochub
RUN chmod +x /www/dochub/DocHub_linux_amd64
CMD ["/www/dochub/DocHub_linux_amd64"]