FROM scratch
ADD ssltest /
ENTRYPOINT ["/ssltest"]
