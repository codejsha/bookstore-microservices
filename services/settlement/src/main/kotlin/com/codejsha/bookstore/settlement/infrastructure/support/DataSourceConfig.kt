package com.codejsha.bookstore.settlement.infrastructure.support

import com.codejsha.bookstore.settlement.config.properties.PaymentDbProperties
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import org.springframework.boot.jdbc.autoconfigure.DataSourceProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Lazy
import org.springframework.context.annotation.Primary

@Configuration
class DataSourceConfig {

    @Bean
    @Primary
    @org.springframework.boot.context.properties.ConfigurationProperties("spring.datasource")
    fun dataSourceProperties(): DataSourceProperties = DataSourceProperties()

    @Bean
    @Primary
    @org.springframework.boot.context.properties.ConfigurationProperties("spring.datasource.hikari")
    fun dataSource(dataSourceProperties: DataSourceProperties): HikariDataSource =
        dataSourceProperties
            .initializeDataSourceBuilder()
            .type(HikariDataSource::class.java)
            .build()

    @Bean(name = ["paymentDataSource"])
    @Lazy
    fun paymentDataSource(
        settlement: SettlementBatchProperties,
        credentials: PaymentDbProperties,
    ): HikariDataSource {
        val conn = settlement.paymentDb
        val config = HikariConfig().apply {
            jdbcUrl = "jdbc:mysql://${conn.host}:${conn.port}/${credentials.name}?${conn.params}"
            username = credentials.username
            password = credentials.password
            driverClassName = "com.mysql.cj.jdbc.Driver"
            poolName = "settlement-payment-ro"
            maximumPoolSize = 2
            minimumIdle = 0
            isReadOnly = true
            isAutoCommit = true
            connectionTimeout = 10_000
        }
        return HikariDataSource(config)
    }
}
