package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
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

class RefundRepoImplUnitTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @Test
    fun `findOne_whenResultEmpty_throwsNoSuchElementExceptionNamingTheUid`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = RefundRepoImpl(dsl)

        val uid = UUID.randomUUID()
        val ex = assertFailsWith<NoSuchElementException> { repo.findOne(uid, ctx) }
        assertContains(ex.message ?: "", uid.toString())

        val sql = capture.queries.first().lowercase()
        assertContains(sql, "from refunds")
        assertContains(sql, "where")
        assertContains(sql, "deleted_at is null")
    }

    @Test
    fun `findAll_whenPaymentIdFilterGiven_bindsValue`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = RefundRepoImpl(dsl)

        val option = RefundQueryOption(paymentId = "pay_001")
        repo.findAll(option, PageRequest.of(0, 10), ctx)

        val data = capture.queries[0].lowercase()
        assertContains(data, "from refunds")
        assertContains(data, "payment_id")
        assertContains(capture.bindings[0].toList(), "pay_001")
    }

    @Test
    fun `create_whenRefundTypeNull_emitsInsertWithInstantDefault`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = RefundRepoImpl(dsl)

        val command = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 5_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = null,
            metadata = null,
        )
        runCatching { repo.create(command, ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val insert = capture.queries.first().lowercase()
        assertContains(insert, "insert into refunds")

        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "pay_001")
        assertContains(bindings, 5_000L)
        assertContains(bindings, "KRW")
        assertContains(bindings, "pending")
        assertContains(bindings, "instant")
    }

    @Test
    fun `create_whenRefundTypeGiven_emitsInsertWithThatValue`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = RefundRepoImpl(dsl)

        val command = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 5_000L,
            currency = "KRW",
            reason = null,
            refundType = "scheduled",
            metadata = null,
        )
        runCatching { repo.create(command, ctx) }

        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "scheduled")
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
