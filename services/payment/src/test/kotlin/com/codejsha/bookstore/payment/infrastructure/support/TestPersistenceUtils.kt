package com.codejsha.bookstore.payment.infrastructure.support

import org.jooq.Configuration
import org.jooq.DSLContext
import org.jooq.SQLDialect
import org.jooq.conf.RenderKeywordCase
import org.jooq.conf.RenderQuotedNames
import org.jooq.conf.Settings
import org.jooq.impl.DSL
import org.jooq.impl.DefaultConfiguration
import org.jooq.tools.jdbc.MockConnection
import org.jooq.tools.jdbc.MockDataProvider
import java.sql.Connection

class TestPersistenceUtils {

    companion object {
        fun mockConnection(provider: MockDataProvider): MockConnection =
            MockConnection(provider)

        fun jooqConfiguration(connection: Connection): Configuration {
            System.setProperty("org.jooq.no-logo", "true")
            System.setProperty("org.jooq.no-tips", "true")
            val settings =
                Settings()
                    .withExecuteWithOptimisticLocking(true)
                    .withExecuteLogging(false)
                    .withRenderFormatted(false)
                    .withRenderKeywordCase(RenderKeywordCase.UPPER)
                    .withRenderQuotedNames(RenderQuotedNames.EXPLICIT_DEFAULT_UNQUOTED)
                    .withRenderSchema(false)
            return DefaultConfiguration()
                .set(SQLDialect.MYSQL)
                .set(settings)
                .set(connection)
        }

        fun dslContext(jooqConfiguration: Configuration): DSLContext =
            DSL.using(jooqConfiguration)

        fun dslWithProvider(provider: MockDataProvider): DSLContext {
            val connection = mockConnection(provider)
            val config = jooqConfiguration(connection)
            return dslContext(config)
        }
    }
}
