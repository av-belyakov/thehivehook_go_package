FROM golang:1.23.2-alpine AS temporary_image
WORKDIR /go/src/
ENV PATH /usr/local/go/bin:$PATH
RUN apk update && \
    apk add --no-cache git && \
    git clone https://github.com/av-belyakov/thehivehook_go_package.git /go/src/
RUN go build -C cmd/ -o ../thehivehook_go_package

FROM alpine
LABEL author="Artemij Belyakov"
WORKDIR /opt/thehivehook_go_package
RUN mkdir /opt/thehivehook_go_package/configs && \
    mkdir /opt/thehivehook_go_package/logs
COPY --from=temporary_image /go/src/thehivehook_go_package /opt/thehivehook_go_package/ 
COPY --from=temporary_image /go/src/README.md /opt/thehivehook_go_package/
COPY configs/* /opt/thehivehook_go_package/configs/
#EXPOSE 13000
ENTRYPOINT [ "./thehivehook_go_package" ]