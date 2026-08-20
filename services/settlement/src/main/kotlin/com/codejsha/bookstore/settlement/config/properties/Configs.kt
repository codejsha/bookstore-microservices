package com.codejsha.bookstore.settlement.config.properties

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
