FROM scratch

ARG EXE_NAME

COPY ${EXE_NAME} /usr/bin/rust_binary

ENTRYPOINT [ "/usr/bin/rust_binary" ]