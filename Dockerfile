FROM alpine as build
# install compiler
RUN apk add go git
RUN go env -w GOPROXY=direct
# cache dependencies
ADD go.mod go.sum ./
RUN go mod download
# install aws-nuke
RUN wget https://github.com/rebuy-de/aws-nuke/releases/download/v2.15.0-rc.3/aws-nuke-v2.15.0.rc.3-linux-amd64.tar.gz
RUN tar zxpf aws-nuke-v2.15.0.rc.3-linux-amd64.tar.gz
RUN rm aws-nuke-v2.15.0.rc.3-linux-amd64.tar.gz
# build
ADD . .
RUN go build -o /main
# copy artifacts to a clean image
FROM alpine
COPY --from=build /main /main
COPY --from=build /aws-nuke-v2.15.0.rc.3-linux-amd64 /usr/local/bin/aws-nuke
RUN chmod 755 /usr/local/bin/aws-nuke
RUN mkdir /configs
COPY *config.yaml /configs
ENTRYPOINT [ "/main" ]
WORKDIR /
