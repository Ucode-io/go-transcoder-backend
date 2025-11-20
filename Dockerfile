FROM golang:1.21.1 as builder

#
RUN mkdir -p $GOPATH/src/gitlab.udevs.io/ucode/ucode_go_transcoder_backend
WORKDIR $GOPATH/src/gitlab.udevs.io/ucode/ucode_go_transcoder_backend

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make build && \
    mv ./bin/ucode_go_transcoder_backend / && \
    mv ./config /

FROM alpine
COPY --from=builder ucode_go_transcoder_backend .
COPY --from=builder ./config ./config

ENTRYPOINT ["/ucode_go_transcoder_backend"]