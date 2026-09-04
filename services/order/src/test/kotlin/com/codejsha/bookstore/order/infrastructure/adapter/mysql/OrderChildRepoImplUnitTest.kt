package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.domain.model.command.OrderAdjustmentCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.infrastructure.support.TestPersistenceUtils
import com.codejsha.bookstore.order.support.OrderTestFixtures
import org.jooq.SQLDialect
import org.jooq.impl.DSL
import org.jooq.tools.jdbc.MockDataProvider
import org.jooq.tools.jdbc.MockResult
import org.junit.jupiter.api.Test
import org.springframework.data.domain.PageRequest
import java.math.BigDecimal
import java.util.UUID
import kotlin.test.assertContains
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse

class OrderChildRepoImplUnitTest {

    private val ctx = OrderTestFixtures.DEFAULT_CONTEXT
    private val orderUid: UUID = UUID.randomUUID()

    @Test
    fun `orderItemFindAllByOrder_whenOrderMissing_throws`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderItemRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> {
            repo.findAllByOrder(orderUid, PageRequest.of(0, 10), ctx)
        }
        assertContains(capture.queries.first().lowercase(), "from orders")
    }

    @Test
    fun `orderItemCreate_whenOrderMissing_throws`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderItemRepoImpl(dsl)

        val cmd = OrderItemCreateCommand(
            productId = 1L, sku = null, productName = null, options = null,
            quantity = 1, currency = "KRW", price = BigDecimal.ONE, taxRate = BigDecimal.ZERO,
        )
        assertFailsWith<NoSuchElementException> { repo.create(orderUid, cmd, ctx) }
    }

    @Test
    fun `orderItemCreate_whenInserting_omitsTheGeneratedSubtotalColumn`() {
        val capture = QueryCapture()
        var call = 0
        val dsl = TestPersistenceUtils.dslWithProvider(
            capture.provider {
                call += 1
                if (call == 1) arrayOf(MockResult(1, orderIdResult())) else arrayOf(MockResult(1, EMPTY_RESULT))
            }
        )
        val repo = OrderItemRepoImpl(dsl)

        val cmd = OrderItemCreateCommand(
            productId = 1L, sku = null, productName = null, options = null,
            quantity = 2, currency = "KRW", price = BigDecimal("100.00"), taxRate = BigDecimal.ZERO,
        )
        assertFailsWith<NoSuchElementException> { repo.create(orderUid, cmd, ctx) }

        val insert = capture.queries.first { it.lowercase().startsWith("insert") }
        assertFalse(insert.lowercase().contains("subtotal"))
    }

    @Test
    fun `orderShippingFindByOrder_whenOrderMissing_throws`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderShippingRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> { repo.findByOrder(orderUid, ctx) }
        assertContains(capture.queries.first().lowercase(), "from orders")
    }

    @Test
    fun `orderShippingCreate_whenOrderMissing_throws`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderShippingRepoImpl(dsl)

        val cmd = OrderShippingCreateCommand(
            recipientName = "Alice", recipientPhone = "010-1",
            addressLine1 = "1", addressLine2 = null,
            city = "Seoul", state = "Seoul", postalCode = "12345",
            country = "KR", shippingMethod = "STANDARD",
        )
        assertFailsWith<NoSuchElementException> { repo.create(orderUid, cmd, ctx) }
    }

    @Test
    fun `orderAdjustmentFindAllByOrder_whenOrderMissing_throws`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderAdjustmentRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> {
            repo.findAllByOrder(orderUid, PageRequest.of(0, 10), ctx)
        }
        assertContains(capture.queries.first().lowercase(), "from orders")
    }

    @Test
    fun `orderAdjustmentCreate_whenOrderMissing_throws`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderAdjustmentRepoImpl(dsl)

        val cmd = OrderAdjustmentCreateCommand(
            type = "COUPON", label = "WELCOME",
            amount = BigDecimal("-1000.00"), meta = null,
        )
        assertFailsWith<NoSuchElementException> { repo.create(orderUid, cmd, ctx) }
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

        private fun orderIdResult(): org.jooq.Result<*> {
            val idField = DSL.field("id", Long::class.javaObjectType)
            val dsl = DSL.using(SQLDialect.MYSQL)
            val result = dsl.newResult(idField)
            result.add(dsl.newRecord(idField).also { it.setValue(idField, 42L) })
            return result
        }
    }
}
