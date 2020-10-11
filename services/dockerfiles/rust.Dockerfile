FROM alpine

ARG EXE_NAME

COPY ${EXE_NAME} /rust_binary

ENTRYPOINT [ "/rust_binary" ]