package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.order.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.order.application.port.repo.OrderShippingRepo
import com.codejsha.bookstore.order.application.port.repo.OrderShippingResult
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingUpdateCommand
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_SHIPPING
import com.codejsha.platform.shared.data.ActorContext
import org.jooq.DSLContext
import org.jooq.Record
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class OrderShippingRepoImpl(
    private val dslContext: DSLContext
) : OrderShippingRepo {
    private val os = ORDER_SHIPPING.`as`("os")
    private val o = ORDERS.`as`("o")

    private fun resolveOrderId(orderUid: UUID): Long {
        return dslContext
            .select(o.ID)
            .from(o)
            .where(o.UID.eq(uuidToBytes(orderUid)).and(o.DELETED_AT.isNull))
            .fetchOneInto(Long::class.java)
            ?: throw NoSuchElementException("Order with uid $orderUid not found")
    }

    override fun findByOrder(orderUid: UUID, context: ActorContext): OrderShippingResult? {
        val orderId = resolveOrderId(orderUid)

        return dslContext
            .select(os.asterisk())
            .from(os)
            .where(os.ORDER_ID.eq(orderId).and(os.DELETED_AT.isNull))
            .fetchOne { toOrderShippingResult(it) }
    }

    override fun create(
        orderUid: UUID, command: OrderShippingCreateCommand, context: ActorContext
    ): OrderShippingResult {
        val orderId = resolveOrderId(orderUid)
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(os)
            .set(os.UID, uuidToBytes(uid))
            .set(os.ORDER_ID, orderId)
            .set(os.RECIPIENT_NAME, command.recipientName)
            .set(os.RECIPIENT_PHONE, command.recipientPhone)
            .set(os.ADDRESS_LINE1, command.addressLine1)
            .set(os.ADDRESS_LINE2, command.addressLine2)
            .set(os.CITY, command.city)
            .set(os.STATE, command.state)
            .set(os.POSTAL_CODE, command.postalCode)
            .set(os.COUNTRY, command.country)
            .set(os.SHIPPING_METHOD, command.shippingMethod)
            .set(os.CREATED_AT, LocalDateTime.now())
            .set(os.ACTOR_ID, context.actorId)
            .execute()

        return findByOrder(orderUid, context)
            ?: throw IllegalStateException("OrderShipping not found after creation")
    }

    override fun update(
        orderUid: UUID, command: OrderShippingUpdateCommand, context: ActorContext
    ): OrderShippingResult {
        val orderId = resolveOrderId(orderUid)

        val updated = dslContext
            .update(os)
            .apply {
                command.recipientName?.let { set(os.RECIPIENT_NAME, it) }
                command.recipientPhone?.let { set(os.RECIPIENT_PHONE, it) }
                command.addressLine1?.let { set(os.ADDRESS_LINE1, it) }
                command.addressLine2?.let { set(os.ADDRESS_LINE2, it) }
                command.city?.let { set(os.CITY, it) }
                command.state?.let { set(os.STATE, it) }
                command.postalCode?.let { set(os.POSTAL_CODE, it) }
                command.country?.let { set(os.COUNTRY, it) }
                command.shippingMethod?.let { set(os.SHIPPING_METHOD, it) }
            }
            .set(os.UPDATED_AT, LocalDateTime.now())
            .set(os.ACTOR_ID, context.actorId)
            .where(os.ORDER_ID.eq(orderId).and(os.DELETED_AT.isNull))
            .execute()

        check(updated > 0) { "OrderShipping for order uid $orderUid not found" }
        return findByOrder(orderUid, context)
            ?: throw IllegalStateException("OrderShipping not found after update")
    }

    override fun delete(orderUid: UUID, context: ActorContext) {
        val orderId = resolveOrderId(orderUid)

        dslContext
            .update(os)
            .set(os.DELETED_AT, LocalDateTime.now())
            .set(os.ACTOR_ID, context.actorId)
            .where(os.ORDER_ID.eq(orderId).and(os.DELETED_AT.isNull))
            .execute()
    }

    private fun toOrderShippingResult(record: Record): OrderShippingResult {
        return OrderShippingResult(
            id = record.get(os.ID)!!.toLong(),
            uid = bytesToUuid(record.get(os.UID)!!),
            orderId = record.get(os.ORDER_ID)!!,
            recipientName = record.get(os.RECIPIENT_NAME)!!,
            recipientPhone = record.get(os.RECIPIENT_PHONE)!!,
            addressLine1 = record.get(os.ADDRESS_LINE1)!!,
            addressLine2 = record.get(os.ADDRESS_LINE2),
            city = record.get(os.CITY)!!,
            state = record.get(os.STATE)!!,
            postalCode = record.get(os.POSTAL_CODE)!!,
            country = record.get(os.COUNTRY)!!,
            shippingMethod = record.get(os.SHIPPING_METHOD)!!,
            createdAt = record.get(os.CREATED_AT)!!,
            updatedAt = record.get(os.UPDATED_AT),
        )
    }
}
