package com.codejsha.bookstore.order.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_SHIPPING
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingUpdateCommand
import com.codejsha.bookstore.order.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class OrderShippingRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var shippingRepo: OrderShippingRepoImpl

    private val otherActor = ActorContext(actorId = 9L, actorType = ActorType.USER)

    @Test
    fun `create_orderExists_persistsEveryAddressField`() {
        val order = orderRepo.create(orderCommand(), ctx)

        val created = shippingRepo.create(order.uid, shippingCommand(addressLine2 = "Apt 2"), ctx)

        val found = assertNotNull(shippingRepo.findByOrder(order.uid, ctx))
        assertEquals(created.uid, found.uid)
        assertEquals(order.id, found.orderId)
        assertEquals("Alice", found.recipientName)
        assertEquals("010-1234-5678", found.recipientPhone)
        assertEquals("1 Foo St", found.addressLine1)
        assertEquals("Apt 2", found.addressLine2)
        assertEquals("Seoul", found.city)
        assertEquals("Gangnam", found.state)
        assertEquals("12345", found.postalCode)
        assertEquals("KR", found.country)
        assertEquals("STANDARD", found.shippingMethod)
        assertNull(found.updatedAt)
        val row = rowOf(created.uid)
        assertEquals(1L, row.get(ORDER_SHIPPING.VERSION))
        assertEquals(ctx.actorId, row.get(ORDER_SHIPPING.ACTOR_ID))
    }

    @Test
    fun `create_addressLine2Null_persistsNull`() {
        val order = orderRepo.create(orderCommand(), ctx)

        val created = shippingRepo.create(order.uid, shippingCommand(addressLine2 = null), ctx)

        assertNull(created.addressLine2)
    }

    @Test
    fun `findByOrder_orderWithoutShipping_yieldsNull`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        shippingRepo.create(other.uid, shippingCommand(), ctx)

        assertNull(shippingRepo.findByOrder(order.uid, ctx))
    }

    @Test
    fun `update_partialCommand_changesOnlyGivenFieldsAndStampsUpdate`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val created = shippingRepo.create(order.uid, shippingCommand(), ctx)

        val updated = shippingRepo.update(order.uid, updateCommand(city = "Busan"), otherActor)

        assertEquals(created.uid, updated.uid)
        assertEquals("Busan", updated.city)
        assertEquals("Alice", updated.recipientName)
        assertEquals("1 Foo St", updated.addressLine1)
        assertEquals("STANDARD", updated.shippingMethod)
        assertNotNull(updated.updatedAt)
        assertEquals(otherActor.actorId, rowOf(created.uid).get(ORDER_SHIPPING.ACTOR_ID))
    }

    @Test
    fun `update_everyFieldGiven_overwritesAllFields`() {
        val order = orderRepo.create(orderCommand(), ctx)
        shippingRepo.create(order.uid, shippingCommand(), ctx)

        val updated = shippingRepo.update(
            order.uid,
            OrderShippingUpdateCommand(
                recipientName = "Bob",
                recipientPhone = "010-9999-0000",
                addressLine1 = "9 Bar Ave",
                addressLine2 = "Unit 7",
                city = "Incheon",
                state = "Yeonsu",
                postalCode = "54321",
                country = "US",
                shippingMethod = "EXPRESS",
            ),
            ctx,
        )

        assertEquals("Bob", updated.recipientName)
        assertEquals("010-9999-0000", updated.recipientPhone)
        assertEquals("9 Bar Ave", updated.addressLine1)
        assertEquals("Unit 7", updated.addressLine2)
        assertEquals("Incheon", updated.city)
        assertEquals("Yeonsu", updated.state)
        assertEquals("54321", updated.postalCode)
        assertEquals("US", updated.country)
        assertEquals("EXPRESS", updated.shippingMethod)
    }

    @Test
    fun `update_orderWithoutShipping_throwsIllegalStateException`() {
        val order = orderRepo.create(orderCommand(), ctx)

        assertFailsWith<IllegalStateException> { shippingRepo.update(order.uid, updateCommand(city = "Busan"), ctx) }
    }

    @Test
    fun `update_unknownOrder_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> {
            shippingRepo.update(UUID.randomUUID(), updateCommand(city = "Busan"), ctx)
        }
    }

    @Test
    fun `delete_shippingExists_keepsRowAndHidesItFromReads`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val created = shippingRepo.create(order.uid, shippingCommand(), ctx)

        shippingRepo.delete(order.uid, otherActor)

        assertNull(shippingRepo.findByOrder(order.uid, ctx))
        val row = rowOf(created.uid)
        assertNotNull(row.get(ORDER_SHIPPING.DELETED_AT))
        assertEquals(otherActor.actorId, row.get(ORDER_SHIPPING.ACTOR_ID))
    }

    @Test
    fun `delete_shippingOfAnotherOrder_leavesItUntouched`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val other = orderRepo.create(orderCommand(), ctx)
        val created = shippingRepo.create(order.uid, shippingCommand(), ctx)

        shippingRepo.delete(other.uid, ctx)

        assertEquals(created.uid, shippingRepo.findByOrder(order.uid, ctx)?.uid)
    }

    @Test
    fun `create_previousShippingSoftDeleted_readsTheReplacement`() {
        val order = orderRepo.create(orderCommand(), ctx)
        val first = shippingRepo.create(order.uid, shippingCommand(), ctx)
        shippingRepo.delete(order.uid, ctx)

        val second = shippingRepo.create(order.uid, shippingCommand(), ctx)

        assertEquals(second.uid, shippingRepo.findByOrder(order.uid, ctx)?.uid)
        assertEquals(2, dsl.fetchCount(ORDER_SHIPPING))
        assertNotNull(rowOf(first.uid).get(ORDER_SHIPPING.DELETED_AT))
    }

    @Test
    fun `update_shippingSoftDeleted_throwsIllegalStateException`() {
        val order = orderRepo.create(orderCommand(), ctx)
        shippingRepo.create(order.uid, shippingCommand(), ctx)
        shippingRepo.delete(order.uid, ctx)

        assertFailsWith<IllegalStateException> { shippingRepo.update(order.uid, updateCommand(city = "Busan"), ctx) }
    }

    private fun shippingCommand(addressLine2: String? = null) = OrderShippingCreateCommand(
        recipientName = "Alice",
        recipientPhone = "010-1234-5678",
        addressLine1 = "1 Foo St",
        addressLine2 = addressLine2,
        city = "Seoul",
        state = "Gangnam",
        postalCode = "12345",
        country = "KR",
        shippingMethod = "STANDARD",
    )

    private fun updateCommand(city: String? = null) = OrderShippingUpdateCommand(
        recipientName = null, recipientPhone = null, addressLine1 = null, addressLine2 = null,
        city = city, state = null, postalCode = null, country = null, shippingMethod = null,
    )

    private fun rowOf(uid: UUID) =
        dsl.selectFrom(ORDER_SHIPPING).where(ORDER_SHIPPING.UID.eq(uid)).fetchSingle()
}
