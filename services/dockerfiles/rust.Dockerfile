FROM rust:1.48.0

ARG EXE_NAME

COPY ${EXE_NAME} /rust_binary
COPY ./static* /static

ENTRYPOINT [ "/rust_binary" ]