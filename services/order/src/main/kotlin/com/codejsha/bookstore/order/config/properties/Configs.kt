package com.codejsha.bookstore.order.config.properties

import org.springframework.boot.context.properties.ConfigurationProperties
import java.time.Duration

@ConfigurationProperties(prefix = "app")
class AppConfig

@ConfigurationProperties(prefix = "spring")
class SpringConfig(
    val application: ApplicationConfig
)

class ApplicationConfig(
    val name: String,
    val version: String
)

@ConfigurationProperties(prefix = "grpc")
class GrpcConfig(
    val server: GrpcServerConfig
)

class GrpcServerConfig(
    val host: String,
    val port: Int
)

@ConfigurationProperties(prefix = "temporal.auth")
class TemporalAuthConfig(
    val enabled: Boolean = false,
    val tokenUrl: String = "",
    val clientId: String = "",
    val clientSecret: String = "",
    val refreshRatio: Double = 0.75,
    val connectTimeout: Duration = Duration.ofSeconds(2),
    val readTimeout: Duration = Duration.ofSeconds(5),
)
