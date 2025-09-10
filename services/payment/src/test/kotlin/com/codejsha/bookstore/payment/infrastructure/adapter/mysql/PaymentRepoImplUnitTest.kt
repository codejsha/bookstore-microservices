package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
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
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class PaymentRepoImplUnitTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @Test
    fun `findOne throws NoSuchElementException when no record matches`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = PaymentRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> { repo.findOne(UUID.randomUUID(), ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val select = capture.queries.first().lowercase()
        assertContains(select, "select")
        assertContains(select, "from payments")
        assertContains(select, "where")
        assertContains(select, "uid")
        assertContains(select, "deleted_at is null")
    }

    @Test
    fun `findAll without filters emits SELECT with deleted_at IS NULL only`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = PaymentRepoImpl(dsl)

        val pageable = PageRequest.of(0, 10)
        val result = repo.findAll(PaymentQueryOption(), pageable, ctx)

        assertEquals(0L, result.totalElements)
        assertEquals(2, capture.queries.size)
        val dataQuery = capture.queries[0].lowercase()
        assertContains(dataQuery, "deleted_at is null")
        assertTrue(
            !dataQuery.contains("customer_id ="),
            "WHERE clause should not include customer_id: $dataQuery",
        )
    }

    @Test
    fun `findAll with all filters adds WHERE conditions and bindings`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = PaymentRepoImpl(dsl)

        val option = PaymentQueryOption(
            customerId = "cus_001",
            status = "succeeded",
            connector = "stripe",
        )
        repo.findAll(option, PageRequest.of(0, 10), ctx)

        val dataQuery = capture.queries[0].lowercase()
        assertContains(dataQuery, "customer_id")
        assertContains(dataQuery, "status")
        assertContains(dataQuery, "connector")
        assertContains(dataQuery, "deleted_at is null")

        val dataBindings = capture.bindings[0].toList()
        assertContains(dataBindings, "cus_001")
        assertContains(dataBindings, "succeeded")
        assertContains(dataBindings, "stripe")
    }

    @Test
    fun `create emits INSERT INTO payments with command bindings and default status`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = PaymentRepoImpl(dsl)

        val command = PaymentCreateCommand(
            amount = 9_900L,
            currency = "KRW",
            customerId = "cus_1",
            paymentMethod = "card",
            paymentMethodType = "credit",
            authenticationType = "no_three_ds",
            setupFutureUsage = null,
            description = "buy book",
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
        )

        runCatching { repo.create(command, ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val insert = capture.queries.first().lowercase()
        assertContains(insert, "insert into payments")

        val insertBindings = capture.bindings.first().toList()
        assertContains(insertBindings, 9_900L)
        assertContains(insertBindings, "KRW")
        assertContains(insertBindings, "cus_1")
        assertContains(insertBindings, "requires_payment_method")
        assertContains(insertBindings, "automatic")
        assertContains(insertBindings, "buy book")
    }

    @Test
    fun `create with a captured terminal status stamps captured_at`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = PaymentRepoImpl(dsl)

        val command = PaymentCreateCommand(
            amount = 9_900L,
            currency = "KRW",
            customerId = "cus_1",
            paymentMethod = "card",
            paymentMethodType = "credit",
            authenticationType = "no_three_ds",
            setupFutureUsage = null,
            description = null,
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
            status = "succeeded",
            amountCaptured = 9_900L,
        )

        runCatching { repo.create(command, ctx) }

        val insert = capture.queries.first().lowercase()
        assertContains(insert, "insert into payments")
        assertContains(insert, "captured_at")
    }

    @Test
    fun `create with a non-captured status leaves captured_at unset`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = PaymentRepoImpl(dsl)

        val command = PaymentCreateCommand(
            amount = 9_900L,
            currency = "KRW",
            customerId = "cus_1",
            paymentMethod = "card",
            paymentMethodType = "credit",
            authenticationType = "no_three_ds",
            setupFutureUsage = null,
            description = null,
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
            status = "failed",
        )

        runCatching { repo.create(command, ctx) }

        val insert = capture.queries.first().lowercase()
        assertContains(insert, "insert into payments")
        assertTrue(
            !insert.contains("captured_at"),
            "captured_at should not be set for a non-captured status: $insert",
        )
    }

    @Test
    fun `update only sets non-null fields and always sets updated_at`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = PaymentRepoImpl(dsl)

        val uid = UUID.randomUUID()
        val command = PaymentUpdateCommand(
            amount = null,
            currency = null,
            customerId = null,
            paymentMethod = null,
            paymentMethodType = null,
            authenticationType = null,
            setupFutureUsage = null,
            description = "updated desc",
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
        )

        runCatching { repo.update(uid, command, ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val update = capture.queries.first().lowercase()
        assertContains(update, "update payments")
        assertContains(update, "set")
        assertContains(update, "description")
        assertContains(update, "updated_at")
        assertTrue(
            !Regex("set\\b.*amount\\s*=").containsMatchIn(update),
            "amount should not be in SET clause: $update",
        )

        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "updated desc")
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
