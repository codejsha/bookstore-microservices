FROM nginx:1.27-alpine

ENV TZ="Asia/Seoul"

ARG IMAGE_VERSION=1.0.0
ARG IMAGE_REVISION

LABEL org.opencontainers.image.title="admin-web" \
      org.opencontainers.image.description="Admin console SPA served by nginx" \
      org.opencontainers.image.source="https://git.example.com/example-corp/admin-web-source" \
      org.opencontainers.image.vendor="example-corp" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}"

COPY ./ci/nginx.conf /etc/nginx/conf.d/default.conf
COPY ./dist /usr/share/nginx/html

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
