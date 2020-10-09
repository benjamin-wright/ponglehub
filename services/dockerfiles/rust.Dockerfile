FROM alpine

ARG EXE_NAME

COPY ${EXE_NAME} /usr/bin/rust_binary

RUN ls -lah /usr/bin

ENTRYPOINT [ "/usr/bin/rust_binary" ]