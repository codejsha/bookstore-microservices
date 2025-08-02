package com.codejsha.bookstore.order.infrastructure.support

import io.r2dbc.spi.ConnectionFactory
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.r2dbc.config.EnableR2dbcAuditing
import org.springframework.data.r2dbc.core.R2dbcEntityTemplate
import org.springframework.data.r2dbc.repository.config.EnableR2dbcRepositories

@Configuration
@EnableR2dbcRepositories(
    basePackages = [
        "com.codejsha.bookstore.order.application.port.repo",
        "com.codejsha.bookstore.order.infrastructure.adapter.mysql"
    ]
)
@EnableR2dbcAuditing
class PersistenceManager {
    @Bean
    fun r2dbcEntityTemplate(connectionFactory: ConnectionFactory): R2dbcEntityTemplate =
        R2dbcEntityTemplate(connectionFactory)
}
