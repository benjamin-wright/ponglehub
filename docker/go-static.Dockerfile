FROM scratch
ARG EXECUTABLE
COPY ${EXECUTABLE} /app
COPY html /html
ENTRYPOINT [ "/app" ]
