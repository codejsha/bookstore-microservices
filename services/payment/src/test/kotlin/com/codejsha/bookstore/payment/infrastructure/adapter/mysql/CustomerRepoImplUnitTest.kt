package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
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

class CustomerRepoImplUnitTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @Test
    fun `findOne empty throws NoSuchElementException`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = CustomerRepoImpl(dsl)

        assertFailsWith<NoSuchElementException> { repo.findOne(UUID.randomUUID(), ctx) }
        val sql = capture.queries.first().lowercase()
        assertContains(sql, "from customers")
    }

    @Test
    fun `findAll with email filter adds WHERE email and binding`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(0, EMPTY_RESULT)) })
        val repo = CustomerRepoImpl(dsl)

        val option = CustomerQueryOption(email = "alice@example.com")
        repo.findAll(option, PageRequest.of(0, 10), ctx)

        val data = capture.queries[0].lowercase()
        assertContains(data, "email")
        assertContains(capture.bindings[0].toList(), "alice@example.com")
    }

    @Test
    fun `create emits INSERT INTO customers with command bindings`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = CustomerRepoImpl(dsl)

        val command = CustomerCreateCommand(
            customerId = "cus_001",
            name = "Alice",
            email = "alice@example.com",
            phone = "1234",
            phoneCountryCode = "+82",
            description = "test",
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        runCatching { repo.create(command, ctx) }

        assertTrue(capture.queries.isNotEmpty())
        val insert = capture.queries.first().lowercase()
        assertContains(insert, "insert into customers")

        val bindings = capture.bindings.first().toList()
        assertContains(bindings, "cus_001")
        assertContains(bindings, "Alice")
        assertContains(bindings, "alice@example.com")
    }

    @Test
    fun `delete soft-deletes by setting deleted_at when row matches`() {
        val capture = QueryCapture()
        val dsl = TestPersistenceUtils.dslWithProvider(capture.provider { arrayOf(MockResult(1, null)) })
        val repo = CustomerRepoImpl(dsl)

        val uid = UUID.randomUUID()
        repo.delete(uid, ctx)

        assertTrue(capture.queries.isNotEmpty())
        val update = capture.queries.first().lowercase()
        assertContains(update, "update customers")
        assertContains(update, "deleted_at")
        assertContains(update, "where")
        assertContains(update, "uid")
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
