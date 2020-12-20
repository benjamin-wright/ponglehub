FROM rust:1.48.0 as builder

WORKDIR /usr/src
RUN USER=root cargo new myapp

WORKDIR /usr/src/myapp

COPY Cargo.toml Cargo.lock ./
RUN cargo install --path .

COPY src ./src
RUN cargo install --path .


FROM debian:buster-slim
ARG EXE_NAME
RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/${EXE_NAME} /usr/local/bin/myapp
COPY static /static
CMD ["myapp"]