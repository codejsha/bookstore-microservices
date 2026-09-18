package com.codejsha.bookstore.payment.config.properties

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

@ConfigurationProperties(prefix = "hyperswitch")
class HyperswitchConfig(
    val baseUrl: String = "http://localhost:8080",
    val apiKey: String = "",
    val profileId: String? = null,
    val webhookSecret: String = "",
    val connectTimeout: Duration = Duration.ofSeconds(2),
    val readTimeout: Duration = Duration.ofSeconds(15),
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

@ConfigurationProperties(prefix = "temporal.worker")
class TemporalWorkerConfig(
    val maxConcurrentActivityExecutions: Int = 16,
    val activityTaskPollers: Int = 2
)
