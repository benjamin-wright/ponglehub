FROM alpine AS builder

RUN apk add --no-cache gcc musl-dev openssl-dev curl cargo
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- --default-toolchain none -y
RUN /root/.cargo/bin/rustup toolchain install nightly --allow-downgrade --profile minimal --component clippy

WORKDIR /home/rust/app
