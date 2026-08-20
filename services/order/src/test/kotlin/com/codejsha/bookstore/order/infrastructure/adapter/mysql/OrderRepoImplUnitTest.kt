package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderUpdateCommand
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
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
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class OrderRepoImplUnitTest {

    private val ctx = OrderTestFixtures.DEFAULT_CONTEXT

    @Test
    fun `findOne empty result throws NoSuchElementException`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> { repo.findOne(UUID.randomUUID(), ctx) }
        val sql = capture.queries.first().lowercase()
        assertContains(sql, "from orders")
        assertContains(sql, "deleted_at is null")
    }

    @Test
    fun `findAll without filters emits SELECT with deleted_at is null only`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderRepoImpl(dsl)

        val result = repo.findAll(OrderQueryOption(), PageRequest.of(0, 10), ctx)

        assertEquals(0L, result.totalElements)
        assertEquals(2, capture.queries.size)
        val data = capture.queries[0].lowercase()
        assertContains(data, "deleted_at is null")
        assertTrue(
            !data.contains("user_uid ="),
            "WHERE clause should not include user_uid: $data",
        )
    }

    @Test
    fun `findAll with userUid and status binds both`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = OrderRepoImpl(dsl)

        val option = OrderQueryOption(userUid = OrderTestFixtures.USER_UID, status = "PENDING")
        repo.findAll(option, PageRequest.of(0, 10), ctx)

        val data = capture.queries[0].lowercase()
        assertContains(data, "user_uid")
        assertContains(data, "status")
        val bindings = capture.bindings[0].toList()
        assertContains(bindings, "PENDING")
    }

    @Test
    fun `create emits INSERT INTO orders with default status PENDING and actor id`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = OrderRepoImpl(dsl)

        val command = OrderCreateCommand(
            userUid = OrderTestFixtures.USER_UID,
            currency = "KRW",
            itemsAmount = BigDecimal("10000.00"),
            discountAmount = BigDecimal.ZERO,
            shippingAmount = BigDecimal.ZERO,
            taxAmount = BigDecimal.ZERO,
            totalAmount = BigDecimal("10000.00"),
            idempotencyKey = "idem_001",
        )
        runCatching { repo.create(command, ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val insert = capture.queries.first().lowercase()
        assertContains(insert, "insert into orders")
        assertContains(insert, "user_uid")

        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "PENDING")
        assertContains(bindings, "idem_001")
        assertContains(bindings, ctx.actorId)
    }

    @Test
    fun `update only sets non-null fields and updates updated_at`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = OrderRepoImpl(dsl)

        val uid = UUID.randomUUID()
        val command = OrderUpdateCommand(
            userUid = null,
            status = "PAID",
            currency = null,
            itemsAmount = null,
            discountAmount = null,
            shippingAmount = null,
            taxAmount = null,
            totalAmount = null,
        )
        runCatching { repo.update(uid, command, ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val update = capture.queries.first().lowercase()
        assertContains(update, "update orders")
        assertContains(update, "status")
        assertContains(update, "updated_at")
        assertTrue(
            !Regex("set\\b.*currency\\s*=").containsMatchIn(update),
            "currency should not be in SET clause: $update",
        )
        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "PAID")
    }

    @Test
    fun `delete soft-deletes via UPDATE with deleted_at`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = OrderRepoImpl(dsl)

        val uid = UUID.randomUUID()
        repo.delete(uid, ctx)

        assertTrue(capture.queries.isNotEmpty())
        val update = capture.queries.first().lowercase()
        assertContains(update, "update orders")
        assertContains(update, "deleted_at")
        assertTrue(
            !update.startsWith("delete"),
            "delete should be soft (UPDATE), not DELETE: $update",
        )
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
