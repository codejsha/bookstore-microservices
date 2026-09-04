package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.order.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.port.repo.OrderItemResult
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemUpdateCommand
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_ITEM
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.DSLContext
import org.jooq.Record
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class OrderItemRepoImpl(
    private val dslContext: DSLContext
) : OrderItemRepo {
    private val oi = ORDER_ITEM.`as`("oi")
    private val o = ORDERS.`as`("o")

    private fun resolveOrderId(orderUid: UUID): Long {
        return dslContext
            .select(o.ID)
            .from(o)
            .where(o.UID.eq(uuidToBytes(orderUid)).and(o.DELETED_AT.isNull))
            .fetchOneInto(Long::class.java)
            ?: throw NoSuchElementException("Order with uid $orderUid not found")
    }

    override fun findAllByOrder(
        orderUid: UUID, pageable: Pageable, context: ActorContext
    ): Page<OrderItemResult> {
        val orderId = resolveOrderId(orderUid)

        val content = dslContext
            .select(oi.asterisk())
            .from(oi)
            .where(oi.ORDER_ID.eq(orderId).and(oi.DELETED_AT.isNull))
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toOrderItemResult(it) }

        val total = dslContext
            .selectCount()
            .from(oi)
            .where(oi.ORDER_ID.eq(orderId).and(oi.DELETED_AT.isNull))
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(orderUid: UUID, uid: UUID, context: ActorContext): OrderItemResult {
        val orderId = resolveOrderId(orderUid)

        return dslContext
            .select(oi.asterisk())
            .from(oi)
            .where(
                oi.ORDER_ID.eq(orderId)
                    .and(oi.UID.eq(uuidToBytes(uid)))
                    .and(oi.DELETED_AT.isNull)
            )
            .fetchOne { toOrderItemResult(it) }
            ?: throw NoSuchElementException("OrderItem with uid $uid not found")
    }

    override fun create(
        orderUid: UUID, command: OrderItemCreateCommand, context: ActorContext
    ): OrderItemResult {
        val orderId = resolveOrderId(orderUid)
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(oi)
            .set(oi.UID, uuidToBytes(uid))
            .set(oi.ORDER_ID, orderId)
            .set(oi.PRODUCT_ID, command.productId)
            .set(oi.SKU, command.sku)
            .set(oi.PRODUCT_NAME, command.productName)
            .set(oi.OPTIONS, command.options)
            .set(oi.QUANTITY, command.quantity)
            .set(oi.CURRENCY, command.currency)
            .set(oi.PRICE, command.price)
            .set(oi.TAX_RATE, command.taxRate)
            .set(oi.CREATED_AT, LocalDateTime.now())
            .set(oi.ACTOR_ID, context.actorId)
            .execute()

        return findOne(orderUid, uid, context)
    }

    override fun update(
        orderUid: UUID, uid: UUID, command: OrderItemUpdateCommand, context: ActorContext
    ): OrderItemResult {
        val orderId = resolveOrderId(orderUid)

        val updated = dslContext
            .update(oi)
            .apply {
                command.productId?.let { set(oi.PRODUCT_ID, it) }
                command.sku?.let { set(oi.SKU, it) }
                command.productName?.let { set(oi.PRODUCT_NAME, it) }
                command.options?.let { set(oi.OPTIONS, it) }
                command.quantity?.let { set(oi.QUANTITY, it) }
                command.currency?.let { set(oi.CURRENCY, it) }
                command.price?.let { set(oi.PRICE, it) }
                command.taxRate?.let { set(oi.TAX_RATE, it) }
            }
            .set(oi.UPDATED_AT, LocalDateTime.now())
            .set(oi.ACTOR_ID, context.actorId)
            .where(
                oi.ORDER_ID.eq(orderId)
                    .and(oi.UID.eq(uuidToBytes(uid)))
                    .and(oi.DELETED_AT.isNull)
            )
            .execute()

        check(updated > 0) { "OrderItem with uid $uid not found" }
        return findOne(orderUid, uid, context)
    }

    override fun delete(orderUid: UUID, uid: UUID, context: ActorContext) {
        val orderId = resolveOrderId(orderUid)

        dslContext
            .update(oi)
            .set(oi.DELETED_AT, LocalDateTime.now())
            .set(oi.ACTOR_ID, context.actorId)
            .where(
                oi.ORDER_ID.eq(orderId)
                    .and(oi.UID.eq(uuidToBytes(uid)))
                    .and(oi.DELETED_AT.isNull)
            )
            .execute()
    }

    private fun toOrderItemResult(record: Record): OrderItemResult {
        return OrderItemResult(
            id = record.get(oi.ID)!!.toLong(),
            uid = bytesToUuid(record.get(oi.UID)!!),
            orderId = record.get(oi.ORDER_ID)!!,
            productId = record.get(oi.PRODUCT_ID)!!,
            sku = record.get(oi.SKU),
            productName = record.get(oi.PRODUCT_NAME),
            options = record.get(oi.OPTIONS),
            quantity = record.get(oi.QUANTITY)!!,
            currency = record.get(oi.CURRENCY)!!,
            price = record.get(oi.PRICE)!!,
            taxRate = record.get(oi.TAX_RATE)!!,
            subtotal = record.get(oi.SUBTOTAL)!!,
            createdAt = record.get(oi.CREATED_AT)!!,
            updatedAt = record.get(oi.UPDATED_AT),
        )
    }
}
