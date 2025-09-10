package com.codejsha.bookstore.payment.config.properties

import org.springframework.boot.context.properties.ConfigurationProperties

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
)
