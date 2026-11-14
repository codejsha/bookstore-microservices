package com.codejsha.bookstore.order.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_ADJUSTMENT
import com.codejsha.bookstore.order.domain.model.command.OrderAdjustmentCreateCommand
import com.codejsha.bookstore.order.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.jooq.JSONB
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import java.math.BigDecimal
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class OrderAdjustmentRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var adjustmentRepo: OrderAdjustmentRepoImpl

    private val otherActor = ActorContext(actorId = 9L, actorType = ActorType.USER)

    @Test
    fun `create_metaGiven_persistsAdjustmentAndReadsMetaBack`() {
        val order = orderRepo.create(orderCommand(), ctx)

        val created = adjustmentRepo.create(
            order.uid,
            adjustmentCommand(meta = mapOf("code" to "WELCOME", "percent" to 10, "tags" to listOf("new", "web"))),
            ctx,
        )

        val found = adjustmentRepo.findOne(order.uid, created.uid, ctx)
        assertEquals(order.id, found.orderId)
        assertEquals("COUPON", found.type)
        assertEquals("WELCOME", found.label)
        assertEquals(0, BigDecimal("-1000.00").compareTo(found.amount))
        assertEquals("WELCOME", found.meta?.get("code"))
        assertEquals(10, found.meta?.get("percent"))
        assertEquals(listOf("new", "web"), found.meta?.get("tags"))
        assertNull(found.updatedAt)
        val row = rowOf(created.uid)
        assertEquals(1L, row.get(ORDER_ADJUSTMENT.VERSION))
        assertEquals(ctx.actorId, row.get(ORDER_ADJUSTMENT.ACTOR_ID))
    }

    @Test
    fun `create_metaAndLabelNull_persistsNulls`() {
        val order = orderRepo.create(orderCommand(), ctx)

        val created = adjustmentRepo.create(
            order.uid,
            OrderAdjustmentCreateCommand(type = "SHIPPING", label = null, amount = BigDecimal("2500.00"), meta = null),
            ctx,
        )

        assertNull(created.meta)
        assertNull(created.label)
        assertNull(rowOf(created.uid).get(ORDER_ADJUSTMENT.META))
    }

    @Test
    fun `create_orderSoftDeleted_throwsNoSuchElementException`() {
        val order = orderRepo.create(orderCommand(), ctx)
        orderRepo.delete(order.uid, ctx)

        assertFailsWith<NoSuchElementException> { adjustmentRepo.create(order.uid, adjustmentCommand(), ctx) }
        assertEquals(0, dsl.fetchCount(ORDER_ADJUSTMENT))
    }

    @Test
    fun `findOne_storedMetaNotAnObject_readsMetaAsNull`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val created = adjustmentRepo.create(order.uid, adjustmentCommand(), ctx)
        dsl.update(ORDER_ADJUSTMENT)
            .set(ORDER_ADJUSTMENT.META, JSONB.jsonb("[1, 2]"))
            .where(ORDER_ADJUSTMENT.UID.eq(created.uid))
            .execute()

        assertNull(adjustmentRepo.findOne(order.uid, created.uid, ctx).meta)
    }

    @Test
    fun `findOne_adjustmentBelongsToAnotherOrder_throwsNoSuchElementException`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val adjustment = adjustmentRepo.create(order.uid, adjustmentCommand(), ctx)

        assertFailsWith<NoSuchElementException> { adjustmentRepo.findOne(other.uid, adjustment.uid, ctx) }
    }

    @Test
    fun `findOne_unknownAdjustmentUid_throwsNoSuchElementException`() {
        val order = orderRepo.create(orderCommand(), ctx)

        assertFailsWith<NoSuchElementException> { adjustmentRepo.findOne(order.uid, UUID.randomUUID(), ctx) }
    }

    @Test
    fun `findAllByOrder_adjustmentsAcrossOrders_listsOnlyLiveAdjustmentsOfThatOrder`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val kept = adjustmentRepo.create(order.uid, adjustmentCommand(), ctx)
        val removed = adjustmentRepo.create(order.uid, adjustmentCommand(), ctx)
        adjustmentRepo.create(other.uid, adjustmentCommand(), ctx)
        adjustmentRepo.delete(order.uid, removed.uid, ctx)

        val page = adjustmentRepo.findAllByOrder(order.uid, PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(kept.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAllByOrder_sortedByAmountAscending_pagesInAmountOrder`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val high = adjustmentRepo.create(order.uid, adjustmentCommand(amount = BigDecimal("300.00")), ctx)
        val low = adjustmentRepo.create(order.uid, adjustmentCommand(amount = BigDecimal("-500.00")), ctx)
        val mid = adjustmentRepo.create(order.uid, adjustmentCommand(amount = BigDecimal("100.00")), ctx)
        val sort = Sort.by(Sort.Direction.ASC, "amount")

        val first = adjustmentRepo.findAllByOrder(order.uid, PageRequest.of(0, 2, sort), ctx)
        val second = adjustmentRepo.findAllByOrder(order.uid, PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(listOf(low.uid, mid.uid), first.content.map { it.uid })
        assertEquals(listOf(high.uid), second.content.map { it.uid })
    }

    @Test
    fun `delete_adjustmentExists_keepsRowAndHidesItFromReads`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val adjustment = adjustmentRepo.create(order.uid, adjustmentCommand(), ctx)

        adjustmentRepo.delete(order.uid, adjustment.uid, otherActor)

        assertFailsWith<NoSuchElementException> { adjustmentRepo.findOne(order.uid, adjustment.uid, ctx) }
        val row = rowOf(adjustment.uid)
        assertNotNull(row.get(ORDER_ADJUSTMENT.DELETED_AT))
        assertEquals(otherActor.actorId, row.get(ORDER_ADJUSTMENT.ACTOR_ID))
    }

    @Test
    fun `delete_adjustmentBelongsToAnotherOrder_leavesItUntouched`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val adjustment = adjustmentRepo.create(order.uid, adjustmentCommand(), ctx)

        adjustmentRepo.delete(other.uid, adjustment.uid, ctx)

        assertEquals(adjustment.uid, adjustmentRepo.findOne(order.uid, adjustment.uid, ctx).uid)
    }

    @Test
    fun `delete_unknownOrder_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> { adjustmentRepo.delete(UUID.randomUUID(), UUID.randomUUID(), ctx) }
    }

    private fun adjustmentCommand(
        amount: BigDecimal = BigDecimal("-1000.00"),
        meta: Map<String, Any>? = null,
    ) = OrderAdjustmentCreateCommand(type = "COUPON", label = "WELCOME", amount = amount, meta = meta)

    private fun rowOf(uid: UUID) =
        dsl.selectFrom(ORDER_ADJUSTMENT).where(ORDER_ADJUSTMENT.UID.eq(uid)).fetchSingle()
}
