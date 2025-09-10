package com.codejsha.bookstore.settlement.infrastructure.support

import org.jooq.ConnectionProvider
import org.jooq.DSLContext
import org.jooq.SQLDialect
import org.jooq.conf.RenderKeywordCase
import org.jooq.conf.RenderQuotedNames
import org.jooq.conf.Settings
import org.jooq.impl.DSL
import org.jooq.impl.DataSourceConnectionProvider
import org.jooq.impl.DefaultConfiguration
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.ComponentScan
import org.springframework.context.annotation.Configuration
import org.springframework.jdbc.datasource.TransactionAwareDataSourceProxy
import org.springframework.transaction.annotation.EnableTransactionManagement
import javax.sql.DataSource

@Configuration
@EnableTransactionManagement
@ComponentScan(basePackages = ["com.codejsha.bookstore.generated.infrastructure.adapter.jooq"])
class PersistenceManager {
    @Bean
    fun connectionProvider(dataSource: DataSource): DataSourceConnectionProvider =
        DataSourceConnectionProvider(TransactionAwareDataSourceProxy(dataSource))

    @Bean
    fun jooqConfiguration(connectionProvider: ConnectionProvider): org.jooq.Configuration {
        System.setProperty("org.jooq.no-logo", "true")
        System.setProperty("org.jooq.no-tips", "true")
        val settings =
            Settings()
                .withExecuteWithOptimisticLocking(true)
                .withExecuteLogging(false)
                .withRenderFormatted(true)
                .withRenderKeywordCase(RenderKeywordCase.UPPER)
                .withRenderQuotedNames(RenderQuotedNames.EXPLICIT_DEFAULT_UNQUOTED)
                .withRenderSchema(false)
        return DefaultConfiguration()
            .set(SQLDialect.MYSQL)
            .set(settings)
            .set(connectionProvider)
    }

    @Bean
    fun dslContext(jooqConfiguration: org.jooq.Configuration): DSLContext =
        DSL.using(jooqConfiguration)
}
