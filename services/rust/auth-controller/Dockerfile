FROM alpine AS builder

RUN apk add --no-cache gcc musl-dev openssl-dev curl
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- --default-toolchain none -y
RUN /root/.cargo/bin/rustup toolchain install nightly --allow-downgrade --profile minimal --component clippy
ENV PATH="/root/.cargo/bin:${PATH}"

WORKDIR /home/rust

RUN USER=root cargo new auth-controller

WORKDIR /home/rust/auth-controller

COPY Cargo.toml Cargo.lock ./

RUN cargo install --path .

COPY src ./src

RUN cargo build --release
FROM scratch

COPY --from=builder /home/rust/auth-controller/target/release/auth-controller .
USER 1000
CMD ["./auth-controller"]