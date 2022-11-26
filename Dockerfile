FROM gcr.io/distroless/static-debian11
COPY istio-ca-shim-step /
CMD ["/istio-ca-shim-step"]