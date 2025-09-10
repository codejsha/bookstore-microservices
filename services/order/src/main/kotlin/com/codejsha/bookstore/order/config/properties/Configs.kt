package com.codejsha.bookstore.order.config.properties

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
