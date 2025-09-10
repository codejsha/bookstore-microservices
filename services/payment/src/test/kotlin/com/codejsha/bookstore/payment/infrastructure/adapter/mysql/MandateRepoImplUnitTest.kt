package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.TestPersistenceUtils
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import org.jooq.SQLDialect
import org.jooq.impl.DSL
import org.jooq.tools.jdbc.MockDataProvider
import org.jooq.tools.jdbc.MockResult
import org.junit.jupiter.api.Test
import org.springframework.data.domain.PageRequest
import java.util.UUID
import kotlin.test.assertContains
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class MandateRepoImplUnitTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @Test
    fun `findOne empty throws NoSuchElementException`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = MandateRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> { repo.findOne(UUID.randomUUID(), ctx) }
    }

    @Test
    fun `findAll with mandateStatus filter binds value`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = MandateRepoImpl(dsl)

        repo.findAll(MandateQueryOption(mandateStatus = "active"), PageRequest.of(0, 10), ctx)

        val data = capture.queries[0].lowercase()
        assertContains(data, "mandate_status")
        assertContains(capture.bindings[0].toList(), "active")
    }

    @Test
    fun `revoke updates only when status is active`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = MandateRepoImpl(dsl)

        runCatching { repo.revoke(UUID.randomUUID(), ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val update = capture.queries.first().lowercase()
        assertContains(update, "update")
        assertContains(update, "mandates")
        assertContains(update, "mandate_status")
        assertContains(update, "where")
        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "revoked")
        assertContains(bindings, "active")
    }

    @Test
    fun `revoke throws IllegalStateException when no active row exists`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = MandateRepoImpl(dsl)

        assertFailsWith<IllegalStateException> { repo.revoke(UUID.randomUUID(), ctx) }
    }

    private class QueryCapture {
        val queries: MutableList<String> = mutableListOf()
        val bindings: MutableList<Array<Any?>> = mutableListOf()

        fun provider(handler: () -> Array<MockResult>): MockDataProvider =
            MockDataProvider { ctx ->
                queries.add(ctx.sql())
                bindings.add(ctx.bindings())
                handler()
            }
    }

    companion object {
        private val EMPTY_RESULT: org.jooq.Result<*> =
            DSL.using(SQLDialect.MYSQL).newResult(DSL.field("x", Int::class.javaObjectType))
    }
}
