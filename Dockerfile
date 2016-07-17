FROM alpine
ADD ssltest /
ENTRYPOINT ["/ssltest"]
