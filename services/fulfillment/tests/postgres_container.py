from collections.abc import Iterator

POSTGRES_IMAGE = "postgres:16"


def container_runtime_available() -> bool:
    try:
        import docker
    except Exception:
        return False
    try:
        client = docker.from_env()
        client.ping()
        return True
    except Exception:
        return False


def postgres_container_url() -> Iterator[str]:
    from testcontainers.postgres import PostgresContainer

    with PostgresContainer(POSTGRES_IMAGE, driver=None) as postgres:
        _, _, tail = postgres.get_connection_url().partition("://")
        yield f"postgresql+asyncpg://{tail}"
