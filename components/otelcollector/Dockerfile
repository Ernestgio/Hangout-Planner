FROM otel/opentelemetry-collector-contrib:0.96.0

COPY config.yaml /etc/otelcol-contrib/config.yaml

EXPOSE 4317 4318 8888 8889

# Run collector with custom config
CMD ["--config=/etc/otelcol-contrib/config.yaml"]
