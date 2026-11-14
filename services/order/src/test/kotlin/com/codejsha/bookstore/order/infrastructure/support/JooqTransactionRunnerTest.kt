package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull

class JooqTransactionRunnerTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var txRunner: JooqTransactionRunner

    @Test
    fun `tx_blockCompletes_commitsWritesAndYieldsBlockValue`() {
        val uid = txRunner.tx { orderRepo.create(orderCommand(idempotencyKey = "tx-commit"), ctx).uid }

        assertEquals(uid, orderRepo.findByIdempotencyKey("tx-commit", ctx)?.uid)
    }

    @Test
    fun `tx_blockYieldsNull_passesNullThrough`() {
        val result: String? = txRunner.tx { null }

        assertNull(result)
    }

    @Test
    fun `tx_blockThrowsAfterWrites_rollsBackEveryWrite`() {
        val ex = assertFailsWith<IllegalStateException> {
            txRunner.tx {
                orderRepo.create(orderCommand(idempotencyKey = "tx-rollback-1"), ctx)
                orderRepo.create(orderCommand(idempotencyKey = "tx-rollback-2"), ctx)
                error("boom")
            }
        }

        assertEquals("boom", ex.message)
        assertEquals(0, dsl.fetchCount(ORDERS))
    }

    @Test
    fun `tx_uniqueViolationInsideBlock_rollsBackEarlierWrites`() {
        val ex = assertFailsWith<RuntimeException> {
            txRunner.tx {
                orderRepo.create(orderCommand(idempotencyKey = "tx-dup"), ctx)
                orderRepo.create(orderCommand(idempotencyKey = "tx-dup"), ctx)
            }
        }

        assertEquals(UNIQUE_VIOLATION, sqlStateOf(ex))
        assertEquals(0, dsl.fetchCount(ORDERS))
    }
}
