# Dockerfile собирать с аргументом --build-arg
# sudo docker build -t gitlab.cloud.gcm:5050/a.belyakov/thehivehook_go_package:test_image --build-arg VERSION=v0.3.2 .

FROM golang:1.23.4-alpine AS packages_image
WORKDIR /go/src/app
ENV PATH /usr/local/go/bin:$PATH
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.23.4-alpine AS build_image
WORKDIR /go/
RUN apk update && \
    apk add --no-cache git && \
    git clone https://github.com/av-belyakov/thehivehook_go_package.git /go/src/
COPY --from=packages_image /go ./
RUN go build -C ./src/cmd/ -o app

FROM alpine
LABEL author="Artemij Belyakov"
ARG US_DIR=/opt/application_${VERSION}
ARG USERNAME=dockeruser
ARG VERSION
RUN adduser -D ${USERNAME} --home ${US_DIR}
USER ${USERNAME}
WORKDIR ${US_DIR}
RUN mkdir ./configs && \
    mkdir ./logs
COPY --from=build_image /go/src/app ./
COPY --from=build_image /go/src/README.md ./ 
COPY config/* ./config/

ENTRYPOINT [ "./app" ]
