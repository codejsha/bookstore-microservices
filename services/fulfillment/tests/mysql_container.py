from collections.abc import Iterator

MYSQL_IMAGE = "mysql:8.0"


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


def mysql_container_url() -> Iterator[str]:
    from testcontainers.mysql import MySqlContainer

    with MySqlContainer(MYSQL_IMAGE) as mysql:
        _, _, tail = mysql.get_connection_url().partition("://")
        yield f"mysql+aiomysql://{tail}"
