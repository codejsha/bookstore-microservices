FROM python:3.14-slim AS builder

WORKDIR /app

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-dev --no-install-project

COPY . .
RUN uv sync --frozen --no-dev

FROM python:3.14-slim AS runtime

ENV TZ="Asia/Seoul"

ARG IMAGE_VERSION=1.0.0
ARG IMAGE_REVISION

LABEL org.opencontainers.image.title="notification" \
      org.opencontainers.image.description="Notification service: shipment notifications via Temporal" \
      org.opencontainers.image.source="https://git.example.com/example-corp/notification-source" \
      org.opencontainers.image.vendor="example-corp" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}"

RUN groupadd --gid 65532 nonroot \
    && useradd --uid 65532 --gid nonroot --create-home nonroot

WORKDIR /app

COPY --from=builder /app/.venv /app/.venv
COPY --from=builder /app/internal /app/internal
COPY --from=builder /app/main.py /app/main.py
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/db /app/db
COPY --from=builder /app/alembic.ini /app/alembic.ini

ENV PATH="/app/.venv/bin:$PATH"
ENV PYTHONPATH="/app"

USER nonroot

EXPOSE 8080

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
