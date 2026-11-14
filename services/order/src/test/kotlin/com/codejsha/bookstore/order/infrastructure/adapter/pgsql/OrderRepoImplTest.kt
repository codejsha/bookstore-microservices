package com.codejsha.bookstore.order.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import com.codejsha.bookstore.order.domain.model.command.OrderUpdateCommand
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.order.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.order.support.OrderTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.Test
import org.springframework.dao.DuplicateKeyException
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

class OrderRepoImplTest : AbstractPostgresTestSupport() {

    private val otherUser: UUID = UUID.fromString("88888888-8888-8888-8888-888888888888")
    private val otherActor = ActorContext(actorId = 9L, actorType = ActorType.USER)

    @Test
    fun `create_commandGiven_persistsPendingOrderWithVersionOneAndActor`() {
        val created = orderRepo.create(orderCommand(idempotencyKey = "idem-create"), ctx)

        val found = orderRepo.findOne(created.uid, ctx)
        assertTrue(found.id > 0L)
        assertEquals(OrderTestFixtures.USER_UID, found.userUid)
        assertEquals("PENDING", found.status)
        assertEquals("KRW", found.currency)
        assertEquals(0, BigDecimal("10000.00").compareTo(found.totalAmount))
        assertEquals("idem-create", found.idempotencyKey)
        assertTrue(found.orderNumber.matches(Regex("ord_[0-9a-f]{24}")))
        assertNull(found.paymentUid)
        assertNull(found.updatedAt)
        assertEquals(1L, versionOf(created.uid))
        assertEquals(ctx.actorId, actorOf(created.uid))
    }

    @Test
    fun `create_duplicateIdempotencyKey_rejectsWithUniqueViolation`() {
        orderRepo.create(orderCommand(idempotencyKey = "idem-dup"), ctx)

        val ex = assertFailsWith<DuplicateKeyException> {
            orderRepo.create(orderCommand(idempotencyKey = "idem-dup"), ctx)
        }

        assertEquals(UNIQUE_VIOLATION, sqlStateOf(ex))
        assertEquals(1, dsl.fetchCount(ORDERS))
    }

    @Test
    fun `findByIdempotencyKey_keyStored_findsThatOrder`() {
        val created = orderRepo.create(orderCommand(idempotencyKey = "idem-hit"), ctx)
        orderRepo.create(orderCommand(idempotencyKey = "idem-other"), ctx)

        assertEquals(created.uid, orderRepo.findByIdempotencyKey("idem-hit", ctx)?.uid)
    }

    @Test
    fun `findByIdempotencyKey_keyUnknown_yieldsNull`() {
        orderRepo.create(orderCommand(idempotencyKey = "idem-hit"), ctx)

        assertNull(orderRepo.findByIdempotencyKey("idem-miss", ctx))
    }

    @Test
    fun `findByIdempotencyKey_orderSoftDeleted_yieldsNull`() {
        val created = orderRepo.create(orderCommand(idempotencyKey = "idem-gone"), ctx)
        orderRepo.delete(created.uid, ctx)

        assertNull(orderRepo.findByIdempotencyKey("idem-gone", ctx))
    }

    @Test
    fun `update_partialCommand_changesOnlyGivenFieldsAndStampsUpdate`() {
        val created = orderRepo.create(orderCommand(), ctx)

        val updated = orderRepo.update(
            created.uid,
            updateCommand(status = "PAID", totalAmount = BigDecimal("9000.00")),
            otherActor,
        )

        assertEquals("PAID", updated.status)
        assertEquals(0, BigDecimal("9000.00").compareTo(updated.totalAmount))
        assertEquals(0, BigDecimal("10000.00").compareTo(updated.itemsAmount))
        assertEquals("KRW", updated.currency)
        assertEquals(created.orderNumber, updated.orderNumber)
        assertNotNull(updated.updatedAt)
        assertEquals(otherActor.actorId, actorOf(created.uid))
    }

    @Test
    fun `update_everyAmountGiven_overwritesAllAmounts`() {
        val created = orderRepo.create(orderCommand(), ctx)

        val updated = orderRepo.update(
            created.uid,
            OrderUpdateCommand(
                userUid = null,
                status = null,
                currency = "USD",
                itemsAmount = BigDecimal("50.00"),
                discountAmount = BigDecimal("5.00"),
                shippingAmount = BigDecimal("3.00"),
                taxAmount = BigDecimal("4.80"),
                totalAmount = BigDecimal("52.80"),
            ),
            ctx,
        )

        assertEquals("PENDING", updated.status)
        assertEquals("USD", updated.currency)
        assertEquals(0, BigDecimal("50.00").compareTo(updated.itemsAmount))
        assertEquals(0, BigDecimal("5.00").compareTo(updated.discountAmount))
        assertEquals(0, BigDecimal("3.00").compareTo(updated.shippingAmount))
        assertEquals(0, BigDecimal("4.80").compareTo(updated.taxAmount))
        assertEquals(0, BigDecimal("52.80").compareTo(updated.totalAmount))
    }

    @Test
    fun `update_unknownUid_throwsIllegalStateException`() {
        assertFailsWith<IllegalStateException> {
            orderRepo.update(UUID.randomUUID(), updateCommand(status = "PAID"), ctx)
        }
    }

    @Test
    fun `update_orderSoftDeleted_throwsIllegalStateException`() {
        val created = orderRepo.create(orderCommand(), ctx)
        orderRepo.delete(created.uid, ctx)

        assertFailsWith<IllegalStateException> {
            orderRepo.update(created.uid, updateCommand(status = "PAID"), ctx)
        }
    }

    @Test
    fun `transitionStatus_currentStatusAllowed_movesStatusAndBumpsVersion`() {
        val created = orderRepo.create(orderCommand(), ctx)

        val moved = orderRepo.transitionStatus(created.uid, setOf("PENDING"), "PAID", otherActor)

        assertTrue(moved)
        assertEquals("PAID", orderRepo.findOne(created.uid, ctx).status)
        assertEquals(2L, versionOf(created.uid))
        assertEquals(otherActor.actorId, actorOf(created.uid))
    }

    @Test
    fun `transitionStatus_currentStatusNotAllowed_leavesOrderUntouched`() {
        val created = orderRepo.create(orderCommand(), ctx)

        val moved = orderRepo.transitionStatus(created.uid, setOf("PAID"), "SHIPPED", otherActor)

        assertFalse(moved)
        val found = orderRepo.findOne(created.uid, ctx)
        assertEquals("PENDING", found.status)
        assertNull(found.updatedAt)
        assertEquals(1L, versionOf(created.uid))
    }

    @Test
    fun `transitionStatus_alreadyAtTarget_reportsSuccess`() {
        val created = orderRepo.create(orderCommand(), ctx)
        orderRepo.transitionStatus(created.uid, setOf("PENDING"), "PAID", ctx)

        val repeated = orderRepo.transitionStatus(created.uid, setOf("PENDING"), "PAID", ctx)

        assertTrue(repeated)
        assertEquals("PAID", orderRepo.findOne(created.uid, ctx).status)
    }

    @Test
    fun `transitionStatus_orderSoftDeleted_reportsFailure`() {
        val created = orderRepo.create(orderCommand(), ctx)
        orderRepo.delete(created.uid, ctx)

        assertFalse(orderRepo.transitionStatus(created.uid, setOf("PENDING"), "PAID", ctx))
    }

    @Test
    fun `setPaymentUid_orderExists_storesPaymentUid`() {
        val created = orderRepo.create(orderCommand(), ctx)
        val paymentUid = UUID.randomUUID()

        orderRepo.setPaymentUid(created.uid, paymentUid, ctx)

        val found = orderRepo.findOne(created.uid, ctx)
        assertEquals(paymentUid, found.paymentUid)
        assertNotNull(found.updatedAt)
    }

    @Test
    fun `setPaymentUid_unknownUid_throwsIllegalStateException`() {
        assertFailsWith<IllegalStateException> {
            orderRepo.setPaymentUid(UUID.randomUUID(), UUID.randomUUID(), ctx)
        }
    }

    @Test
    fun `delete_orderExists_keepsRowAndHidesItFromReads`() {
        val created = orderRepo.create(orderCommand(), ctx)

        orderRepo.delete(created.uid, otherActor)

        assertFailsWith<NoSuchElementException> { orderRepo.findOne(created.uid, ctx) }
        val row = dsl.select(ORDERS.DELETED_AT, ORDERS.ACTOR_ID).from(ORDERS).where(ORDERS.UID.eq(created.uid)).fetchOne()
        assertNotNull(row)
        assertNotNull(row.get(ORDERS.DELETED_AT))
        assertEquals(otherActor.actorId, row.get(ORDERS.ACTOR_ID))
    }

    @Test
    fun `delete_alreadyDeleted_keepsOriginalDeletionTime`() {
        val created = orderRepo.create(orderCommand(), ctx)
        orderRepo.delete(created.uid, ctx)
        val firstDeletedAt = deletedAtOf(created.uid)

        orderRepo.delete(created.uid, otherActor)

        assertEquals(firstDeletedAt, deletedAtOf(created.uid))
        assertEquals(ctx.actorId, actorOf(created.uid))
    }

    @Test
    fun `findAll_userUidFilter_listsOnlyThatUsersOrders`() {
        orderRepo.create(orderCommand(), ctx)
        val theirs = orderRepo.create(orderCommand(userUid = otherUser), ctx)

        val page = orderRepo.findAll(OrderQueryOption(userUid = otherUser), PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(theirs.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAll_statusFilter_listsOnlyMatchingStatus`() {
        val paid = orderRepo.create(orderCommand(), ctx)
        orderRepo.create(orderCommand(), ctx)
        orderRepo.transitionStatus(paid.uid, setOf("PENDING"), "PAID", ctx)

        val page = orderRepo.findAll(OrderQueryOption(status = "PAID"), PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(paid.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAll_softDeletedOrderPresent_excludesItFromContentAndTotal`() {
        val kept = orderRepo.create(orderCommand(), ctx)
        val removed = orderRepo.create(orderCommand(), ctx)
        orderRepo.delete(removed.uid, ctx)

        val page = orderRepo.findAll(OrderQueryOption(), PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(kept.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAll_sortedByCreatedAtDescending_pagesNewestFirst`() {
        val uids = (1..3).map { day ->
            val created = orderRepo.create(orderCommand(), ctx)
            dsl.update(ORDERS)
                .set(ORDERS.CREATED_AT, LocalDateTime.of(2026, 1, day, 0, 0))
                .where(ORDERS.UID.eq(created.uid))
                .execute()
            created.uid
        }
        val sort = Sort.by(Sort.Direction.DESC, "created_at")

        val first = orderRepo.findAll(OrderQueryOption(), PageRequest.of(0, 2, sort), ctx)
        val second = orderRepo.findAll(OrderQueryOption(), PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(2, first.totalPages)
        assertEquals(listOf(uids[2], uids[1]), first.content.map { it.uid })
        assertEquals(listOf(uids[0]), second.content.map { it.uid })
    }

    @Test
    fun `findAll_sortedByTotalAmountAscending_ordersByAmount`() {
        val large = orderRepo.create(orderCommand(totalAmount = BigDecimal("300.00")), ctx)
        val small = orderRepo.create(orderCommand(totalAmount = BigDecimal("100.00")), ctx)
        val medium = orderRepo.create(orderCommand(totalAmount = BigDecimal("200.00")), ctx)

        val page = orderRepo.findAll(
            OrderQueryOption(),
            PageRequest.of(0, 10, Sort.by(Sort.Direction.ASC, "total_amount")),
            ctx,
        )

        assertEquals(listOf(small.uid, medium.uid, large.uid), page.content.map { it.uid })
    }

    private fun updateCommand(status: String? = null, totalAmount: BigDecimal? = null) = OrderUpdateCommand(
        userUid = null,
        status = status,
        currency = null,
        itemsAmount = null,
        discountAmount = null,
        shippingAmount = null,
        taxAmount = null,
        totalAmount = totalAmount,
    )

    private fun versionOf(uid: UUID): Long? =
        dsl.select(ORDERS.VERSION).from(ORDERS).where(ORDERS.UID.eq(uid)).fetchOne(ORDERS.VERSION)

    private fun actorOf(uid: UUID): Long? =
        dsl.select(ORDERS.ACTOR_ID).from(ORDERS).where(ORDERS.UID.eq(uid)).fetchOne(ORDERS.ACTOR_ID)

    private fun deletedAtOf(uid: UUID): LocalDateTime? =
        dsl.select(ORDERS.DELETED_AT).from(ORDERS).where(ORDERS.UID.eq(uid)).fetchOne(ORDERS.DELETED_AT)
}
