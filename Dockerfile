FROM alpine:3.13

RUN apk add --no-cache bash curl git

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY ndiag_*.apk /tmp/
RUN apk add --allow-untrusted /tmp/ndiag_*.apk
