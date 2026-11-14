package com.codejsha.bookstore.order.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_ITEM
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemUpdateCommand
import com.codejsha.bookstore.order.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
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

class OrderItemRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var itemRepo: OrderItemRepoImpl

    private val otherActor = ActorContext(actorId = 9L, actorType = ActorType.USER)

    @Test
    fun `create_orderExists_persistsItemWithGeneratedSubtotal`() {
        val order = orderRepo.create(orderCommand(), ctx)

        val created = itemRepo.create(order.uid, itemCommand(quantity = 3, price = BigDecimal("1500.00")), ctx)

        val found = itemRepo.findOne(order.uid, created.uid, ctx)
        assertEquals(order.id, found.orderId)
        assertEquals(200L, found.productId)
        assertEquals("SKU-001", found.sku)
        assertEquals("Book A", found.productName)
        assertEquals("hardcover", found.options)
        assertEquals(3, found.quantity)
        assertEquals("KRW", found.currency)
        assertEquals(0, BigDecimal("1500.00").compareTo(found.price))
        assertEquals(0, BigDecimal("0.1000").compareTo(found.taxRate))
        assertEquals(0, BigDecimal("4500.00").compareTo(found.subtotal))
        assertNull(found.updatedAt)
        val row = rowOf(created.uid)
        assertEquals(1L, row.get(ORDER_ITEM.VERSION))
        assertEquals(ctx.actorId, row.get(ORDER_ITEM.ACTOR_ID))
    }

    @Test
    fun `create_optionalFieldsNull_persistsNulls`() {
        val order = orderRepo.create(orderCommand(), ctx)

        val created = itemRepo.create(
            order.uid,
            OrderItemCreateCommand(
                productId = 7L, sku = null, productName = null, options = null,
                quantity = 1, currency = "KRW", price = BigDecimal("10.00"), taxRate = BigDecimal.ZERO,
            ),
            ctx,
        )

        assertNull(created.sku)
        assertNull(created.productName)
        assertNull(created.options)
    }

    @Test
    fun `create_orderSoftDeleted_throwsNoSuchElementException`() {
        val order = orderRepo.create(orderCommand(), ctx)
        orderRepo.delete(order.uid, ctx)

        assertFailsWith<NoSuchElementException> { itemRepo.create(order.uid, itemCommand(), ctx) }
        assertEquals(0, dsl.fetchCount(ORDER_ITEM))
    }

    @Test
    fun `findOne_itemBelongsToAnotherOrder_throwsNoSuchElementException`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val item = itemRepo.create(order.uid, itemCommand(), ctx)

        assertFailsWith<NoSuchElementException> { itemRepo.findOne(other.uid, item.uid, ctx) }
    }

    @Test
    fun `findOne_unknownItemUid_throwsNoSuchElementException`() {
        val order = orderRepo.create(orderCommand(), ctx)

        assertFailsWith<NoSuchElementException> { itemRepo.findOne(order.uid, UUID.randomUUID(), ctx) }
    }

    @Test
    fun `findAllByOrder_itemsAcrossOrders_listsOnlyLiveItemsOfThatOrder`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val kept = itemRepo.create(order.uid, itemCommand(productId = 1L), ctx)
        val removed = itemRepo.create(order.uid, itemCommand(productId = 2L), ctx)
        itemRepo.create(other.uid, itemCommand(productId = 3L), ctx)
        itemRepo.delete(order.uid, removed.uid, ctx)

        val page = itemRepo.findAllByOrder(order.uid, PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(kept.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAllByOrder_sortedByPriceDescending_pagesInPriceOrder`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val cheap = itemRepo.create(order.uid, itemCommand(productId = 1L, price = BigDecimal("100.00")), ctx)
        val pricey = itemRepo.create(order.uid, itemCommand(productId = 2L, price = BigDecimal("300.00")), ctx)
        val mid = itemRepo.create(order.uid, itemCommand(productId = 3L, price = BigDecimal("200.00")), ctx)
        val sort = Sort.by(Sort.Direction.DESC, "price")

        val first = itemRepo.findAllByOrder(order.uid, PageRequest.of(0, 2, sort), ctx)
        val second = itemRepo.findAllByOrder(order.uid, PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(listOf(pricey.uid, mid.uid), first.content.map { it.uid })
        assertEquals(listOf(cheap.uid), second.content.map { it.uid })
    }

    @Test
    fun `update_partialCommand_changesOnlyGivenFieldsAndRecomputesSubtotal`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val item = itemRepo.create(order.uid, itemCommand(quantity = 2, price = BigDecimal("1000.00")), ctx)

        val updated = itemRepo.update(order.uid, item.uid, updateCommand(quantity = 5), otherActor)

        assertEquals(5, updated.quantity)
        assertEquals(0, BigDecimal("5000.00").compareTo(updated.subtotal))
        assertEquals(0, BigDecimal("1000.00").compareTo(updated.price))
        assertEquals("SKU-001", updated.sku)
        assertEquals("Book A", updated.productName)
        assertNotNull(updated.updatedAt)
        assertEquals(otherActor.actorId, rowOf(item.uid).get(ORDER_ITEM.ACTOR_ID))
    }

    @Test
    fun `update_everyFieldGiven_overwritesAllFields`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val item = itemRepo.create(order.uid, itemCommand(), ctx)

        val updated = itemRepo.update(
            order.uid,
            item.uid,
            OrderItemUpdateCommand(
                productId = 900L, sku = "SKU-900", productName = "Book Z", options = "paperback",
                quantity = 4, currency = "USD", price = BigDecimal("12.50"), taxRate = BigDecimal("0.0500"),
            ),
            ctx,
        )

        assertEquals(900L, updated.productId)
        assertEquals("SKU-900", updated.sku)
        assertEquals("Book Z", updated.productName)
        assertEquals("paperback", updated.options)
        assertEquals(4, updated.quantity)
        assertEquals("USD", updated.currency)
        assertEquals(0, BigDecimal("12.50").compareTo(updated.price))
        assertEquals(0, BigDecimal("0.05").compareTo(updated.taxRate))
        assertEquals(0, BigDecimal("50.00").compareTo(updated.subtotal))
    }

    @Test
    fun `update_unknownItemUid_throwsIllegalStateException`() {
        val order = orderRepo.create(orderCommand(), ctx)

        assertFailsWith<IllegalStateException> {
            itemRepo.update(order.uid, UUID.randomUUID(), updateCommand(quantity = 2), ctx)
        }
    }

    @Test
    fun `update_itemSoftDeleted_throwsIllegalStateException`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val item = itemRepo.create(order.uid, itemCommand(), ctx)
        itemRepo.delete(order.uid, item.uid, ctx)

        assertFailsWith<IllegalStateException> {
            itemRepo.update(order.uid, item.uid, updateCommand(quantity = 2), ctx)
        }
    }

    @Test
    fun `delete_itemExists_keepsRowAndHidesItFromReads`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val item = itemRepo.create(order.uid, itemCommand(), ctx)

        itemRepo.delete(order.uid, item.uid, otherActor)

        assertFailsWith<NoSuchElementException> { itemRepo.findOne(order.uid, item.uid, ctx) }
        val row = rowOf(item.uid)
        assertNotNull(row.get(ORDER_ITEM.DELETED_AT))
        assertEquals(otherActor.actorId, row.get(ORDER_ITEM.ACTOR_ID))
    }

    @Test
    fun `delete_itemBelongsToAnotherOrder_leavesItemUntouched`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val item = itemRepo.create(order.uid, itemCommand(), ctx)

        itemRepo.delete(other.uid, item.uid, ctx)

        assertEquals(item.uid, itemRepo.findOne(order.uid, item.uid, ctx).uid)
    }

    @Test
    fun `delete_unknownOrder_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> { itemRepo.delete(UUID.randomUUID(), UUID.randomUUID(), ctx) }
    }

    private fun itemCommand(
        productId: Long = 200L,
        quantity: Int = 2,
        price: BigDecimal = BigDecimal("5000.00"),
    ) = OrderItemCreateCommand(
        productId = productId,
        sku = "SKU-001",
        productName = "Book A",
        options = "hardcover",
        quantity = quantity,
        currency = "KRW",
        price = price,
        taxRate = BigDecimal("0.1000"),
    )

    private fun updateCommand(quantity: Int? = null) = OrderItemUpdateCommand(
        productId = null, sku = null, productName = null, options = null,
        quantity = quantity, currency = null, price = null, taxRate = null,
    )

    private fun rowOf(uid: UUID) =
        dsl.selectFrom(ORDER_ITEM).where(ORDER_ITEM.UID.eq(uid)).fetchSingle()
}
