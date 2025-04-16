package com.codejsha.bookstore.order.config.properties

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(prefix = "app")
class AppConfig(
    val serviceName: String,
    val env: String
)

@ConfigurationProperties(prefix = "spring")
class SpringConfig(
    val application: ApplicationConfig
)

class ApplicationConfig(
    val name: String
)

@ConfigurationProperties(prefix = "grpc")
class GrpcConfig(
    val server: GrpcServerConfig,
    val bookServer: GrpcServerConfig,
    val paymentServer: GrpcServerConfig,
    val userServer: GrpcServerConfig
)

class GrpcServerConfig(
    val host: String,
    val port: Int
)

@ConfigurationProperties(prefix = "telemetry")
class TelemetryConfig(
    val collector: CollectorConfig
)

class CollectorConfig(
    val url: String
)
