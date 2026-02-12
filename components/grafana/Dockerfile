FROM grafana/grafana:10.3.1

ENV GF_SECURITY_ADMIN_USER=admin \
    GF_SECURITY_ADMIN_PASSWORD=admin \
    GF_USERS_ALLOW_SIGN_UP=false \
    GF_SERVER_ROOT_URL=http://localhost:3000 \
    GF_INSTALL_PLUGINS=""

COPY provisioning /etc/grafana/provisioning

COPY dashboards /var/lib/grafana/dashboards

EXPOSE 3000

USER grafana
