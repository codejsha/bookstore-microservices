package com.codejsha.bookstore.settlement.infrastructure.support

import org.springframework.batch.core.configuration.annotation.EnableBatchProcessing
import org.springframework.batch.core.configuration.annotation.EnableJdbcJobRepository
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.jdbc.support.JdbcTransactionManager
import org.springframework.transaction.PlatformTransactionManager
import javax.sql.DataSource

@Configuration
@EnableBatchProcessing
@EnableJdbcJobRepository
class BatchConfig {
    @Bean
    fun transactionManager(dataSource: DataSource): PlatformTransactionManager =
        JdbcTransactionManager(dataSource)
}
