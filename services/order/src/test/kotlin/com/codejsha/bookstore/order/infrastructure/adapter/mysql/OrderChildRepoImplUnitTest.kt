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

class OrderChildRepoImplUnitTest {

    private val ctx = OrderTestFixtures.DEFAULT_CONTEXT
    private val orderUid: UUID = UUID.randomUUID()

    @Test
    fun `OrderItemRepoImpl findAllByOrder throws when order not found`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderItemRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> {
            repo.findAllByOrder(orderUid, PageRequest.of(0, 10), ctx)
        }
        assertContains(capture.queries.first().lowercase(), "from orders")
    }

    @Test
    fun `OrderItemRepoImpl create throws when order not found`() {
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
    fun `OrderShippingRepoImpl findByOrder throws when order not found`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderShippingRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> { repo.findByOrder(orderUid, ctx) }
        assertContains(capture.queries.first().lowercase(), "from orders")
    }

    @Test
    fun `OrderShippingRepoImpl create throws when order not found`() {
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
    fun `OrderAdjustmentRepoImpl findAllByOrder throws when order not found`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderAdjustmentRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> {
            repo.findAllByOrder(orderUid, PageRequest.of(0, 10), ctx)
        }
        assertContains(capture.queries.first().lowercase(), "from orders")
    }

    @Test
    fun `OrderAdjustmentRepoImpl create throws when order not found`() {
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
    }
}
