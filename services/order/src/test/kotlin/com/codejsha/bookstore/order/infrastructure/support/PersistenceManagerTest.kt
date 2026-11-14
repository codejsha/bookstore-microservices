package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import org.jooq.SQLDialect
import org.jooq.conf.RenderKeywordCase
import org.jooq.conf.RenderQuotedNames
import org.junit.jupiter.api.Test
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.boot.jooq.autoconfigure.ExceptionTranslatorExecuteListener
import org.springframework.jdbc.datasource.TransactionAwareDataSourceProxy
import javax.sql.DataSource
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertIs
import kotlin.test.assertSame
import kotlin.test.assertTrue

class PersistenceManagerTest {

    private val manager = PersistenceManager()
    private val dataSource = mock(DataSource::class.java)

    @Test
    fun `connectionProvider_dataSourceGiven_wrapsItInTransactionAwareProxy`() {
        val provider = manager.connectionProvider(dataSource)

        val proxy = assertIs<TransactionAwareDataSourceProxy>(provider.dataSource())
        assertSame(dataSource, proxy.targetDataSource)
        verifyNoInteractions(dataSource)
    }

    @Test
    fun `jooqConfiguration_providerGiven_usesPostgresDialectWithOptimisticLocking`() {
        val provider = manager.connectionProvider(dataSource)

        val configuration = manager.jooqConfiguration(provider)

        assertEquals(SQLDialect.POSTGRES, configuration.dialect())
        assertSame(provider, configuration.connectionProvider())
        val settings = configuration.settings()
        assertTrue(settings.isExecuteWithOptimisticLocking)
        assertFalse(settings.isExecuteLogging)
        assertFalse(settings.isRenderSchema)
        assertEquals(RenderKeywordCase.UPPER, settings.renderKeywordCase)
        assertEquals(RenderQuotedNames.EXPLICIT_DEFAULT_UNQUOTED, settings.renderQuotedNames)
    }

    @Test
    fun `jooqConfiguration_providerGiven_registersSpringExceptionTranslator`() {
        val configuration = manager.jooqConfiguration(manager.connectionProvider(dataSource))

        val listeners = configuration.executeListenerProviders().map { it.provide() }

        assertEquals(listOf(ExceptionTranslatorExecuteListener.DEFAULT), listeners)
    }

    @Test
    fun `dslContext_configurationGiven_rendersUnquotedSchemaLessSql`() {
        val configuration = manager.jooqConfiguration(manager.connectionProvider(dataSource))

        val dsl = manager.dslContext(configuration)

        assertSame(configuration, dsl.configuration())
        val sql = dsl.select(ORDERS.UID).from(ORDERS).where(ORDERS.DELETED_AT.isNull).sql
        assertTrue(sql.contains("SELECT"), sql)
        assertTrue(sql.contains("orders.uid"), sql)
        assertFalse(sql.contains("\""), sql)
        assertFalse(sql.contains("public."), sql)
    }
}
