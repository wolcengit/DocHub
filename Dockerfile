FROM truthhun/dochub:env

WORKDIR /www/dochub

#RUN  wget https://github.com/TruthHun/DocHub/releases/download/v2.0/DocHub.V2.0_linux_amd64.zip \
#    && apt install unzip -y \
#    && unzip DocHub.V2.0_linux_amd64.zip -d /www/dochub/ \
#    && rm -rf /www/dochub/__MACOSX \
#    && chmod 0777 -R /www/dochub

COPY zoneinfo.zip  /usr/local/go/lib/time/zoneinfo.zip
COPY LICENSE  /www/dochub/LICENSE
COPY dictionary  /www/dochub/dictionary
COPY conf/app.conf.example  /www/dochub/conf/app.conf.example
COPY static  /www/dochub/static
COPY views  /www/dochub/views
COPY DocHub /www/dochub/DocHub

RUN chmod 0777 -R /www/dochub/

CMD [ "./DocHub" ]