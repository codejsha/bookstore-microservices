from pydantic import BaseModel
from pydantic_settings import BaseSettings, PydanticBaseSettingsSource, SettingsConfigDict

from internal.config.cloud_config import CloudConfigSettingsSource

CLOUD_CONFIG_APP_NAME = "support"


class AppConfig(BaseModel):
    logging_level: str = "info"


class ServerConfig(BaseModel):
    port: int = 8080
    mode: str = "prod"


class DatabaseConfig(BaseModel):
    host: str = "localhost"
    port: int = 3306
    db_name: str = "support_db"
    username: str = "root"
    password: str = ""
    vault_env_prefix: str = ""

    @property
    def url(self) -> str:
        return (
            f"mysql+aiomysql://{self.username}:{self.password}@{self.host}:{self.port}/{self.db_name}?charset=utf8mb4"
        )


class TelemetryConfig(BaseModel):
    enabled: bool = False
    endpoint: str = "localhost:4317"
    profiling_endpoint: str = "http://localhost:9095"
    service_name: str = "support"
    service_version: str = "1.0.0"
    sampling_ratio: float = 0.1


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_prefix="SUPPORT__",
        env_nested_delimiter="__",
        env_file="/vault/secrets/db.env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    app: AppConfig = AppConfig()
    server: ServerConfig = ServerConfig()
    database: DatabaseConfig = DatabaseConfig(vault_env_prefix="SUPPORT")
    telemetry: TelemetryConfig = TelemetryConfig()

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls: type[BaseSettings],
        init_settings: PydanticBaseSettingsSource,
        env_settings: PydanticBaseSettingsSource,
        dotenv_settings: PydanticBaseSettingsSource,
        file_secret_settings: PydanticBaseSettingsSource,
    ) -> tuple[PydanticBaseSettingsSource, ...]:
        return (
            init_settings,
            env_settings,
            dotenv_settings,
            CloudConfigSettingsSource(settings_cls, CLOUD_CONFIG_APP_NAME),
            file_secret_settings,
        )
