FROM rust:1.48.0

ARG EXE_NAME

COPY ${EXE_NAME} /rust_binary

ENTRYPOINT [ "/rust_binary" ]