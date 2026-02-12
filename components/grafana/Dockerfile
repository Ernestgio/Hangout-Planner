FROM grafana/grafana:10.3.1

COPY provisioning /etc/grafana/provisioning

COPY dashboards /var/lib/grafana/dashboards

EXPOSE 3000

USER grafana
