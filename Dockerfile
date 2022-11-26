FROM scratch
COPY istio-ca-shim-step /
ENTRYPOINT ["/istio-ca-shim-step"]