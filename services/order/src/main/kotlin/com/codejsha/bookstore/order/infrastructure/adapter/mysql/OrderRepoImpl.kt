package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.order.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.port.repo.OrderResult
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderUpdateCommand
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
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
class OrderRepoImpl(
    private val dslContext: DSLContext
) : OrderRepo {
    private val o = ORDERS.`as`("o")

    private val paymentUidField = o.PAYMENT_UID

    override fun findAll(
        option: OrderQueryOption, pageable: Pageable, context: ActorContext
    ): Page<OrderResult> {
        var condition: Condition = o.DELETED_AT.isNull
        option.userUid?.let { condition = condition.and(o.USER_UID.eq(uuidToBytes(it))) }
        option.status?.let { condition = condition.and(o.STATUS.eq(it)) }

        val content = dslContext
            .select(o.asterisk())
            .from(o)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toOrderResult(it) }

        val total = dslContext
            .selectCount()
            .from(o)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): OrderResult {
        return dslContext
            .select(o.asterisk())
            .from(o)
            .where(o.UID.eq(uuidToBytes(uid)).and(o.DELETED_AT.isNull))
            .fetchOne { toOrderResult(it) }
            ?: throw NoSuchElementException("Order with uid $uid not found")
    }

    override fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): OrderResult? {
        return dslContext
            .select(o.asterisk())
            .from(o)
            .where(o.IDEMPOTENCY_KEY.eq(idempotencyKey).and(o.DELETED_AT.isNull))
            .fetchOne { toOrderResult(it) }
    }

    override fun create(command: OrderCreateCommand, context: ActorContext): OrderResult {
        val uid = Uuid.generateV7().toJavaUuid()
        val orderNumber = "ord_${UUID.randomUUID().toString().replace("-", "").take(24)}"

        dslContext
            .insertInto(o)
            .set(o.UID, uuidToBytes(uid))
            .set(o.USER_UID, uuidToBytes(command.userUid))
            .set(o.ORDER_NUMBER, orderNumber)
            .set(o.STATUS, "PENDING")
            .set(o.CURRENCY, command.currency)
            .set(o.ITEMS_AMOUNT, command.itemsAmount)
            .set(o.DISCOUNT_AMOUNT, command.discountAmount)
            .set(o.SHIPPING_AMOUNT, command.shippingAmount)
            .set(o.TAX_AMOUNT, command.taxAmount)
            .set(o.TOTAL_AMOUNT, command.totalAmount)
            .set(o.IDEMPOTENCY_KEY, command.idempotencyKey)
            .set(o.CREATED_AT, LocalDateTime.now())
            .set(o.ACTOR_ID, context.actorId)
            .set(o.VERSION, 1L)
            .execute()

        return findOne(uid, context)
    }

    override fun update(uid: UUID, command: OrderUpdateCommand, context: ActorContext): OrderResult {
        val updated = dslContext
            .update(o)
            .apply {
                command.status?.let { set(o.STATUS, it) }
                command.currency?.let { set(o.CURRENCY, it) }
                command.itemsAmount?.let { set(o.ITEMS_AMOUNT, it) }
                command.discountAmount?.let { set(o.DISCOUNT_AMOUNT, it) }
                command.shippingAmount?.let { set(o.SHIPPING_AMOUNT, it) }
                command.taxAmount?.let { set(o.TAX_AMOUNT, it) }
                command.totalAmount?.let { set(o.TOTAL_AMOUNT, it) }
            }
            .set(o.UPDATED_AT, LocalDateTime.now())
            .set(o.ACTOR_ID, context.actorId)
            .where(o.UID.eq(uuidToBytes(uid)).and(o.DELETED_AT.isNull))
            .execute()

        check(updated > 0) { "Order with uid $uid not found" }
        return findOne(uid, context)
    }

    override fun transitionStatus(uid: UUID, from: Set<String>, to: String, context: ActorContext): Boolean {
        val updated = dslContext
            .update(o)
            .set(o.STATUS, to)
            .set(o.VERSION, o.VERSION.plus(1))
            .set(o.UPDATED_AT, LocalDateTime.now())
            .set(o.ACTOR_ID, context.actorId)
            .where(
                o.UID.eq(uuidToBytes(uid))
                    .and(o.DELETED_AT.isNull)
                    .and(o.STATUS.`in`(from + to)),
            )
            .execute()

        return updated > 0
    }

    override fun setPaymentUid(uid: UUID, paymentUid: UUID, context: ActorContext) {
        val updated = dslContext
            .update(o)
            .set(paymentUidField, uuidToBytes(paymentUid))
            .set(o.UPDATED_AT, LocalDateTime.now())
            .set(o.ACTOR_ID, context.actorId)
            .where(o.UID.eq(uuidToBytes(uid)).and(o.DELETED_AT.isNull))
            .execute()

        check(updated > 0) { "Order with uid $uid not found" }
    }

    override fun delete(uid: UUID, context: ActorContext) {
        dslContext
            .update(o)
            .set(o.DELETED_AT, LocalDateTime.now())
            .set(o.ACTOR_ID, context.actorId)
            .where(o.UID.eq(uuidToBytes(uid)).and(o.DELETED_AT.isNull))
            .execute()
    }

    private fun toOrderResult(record: Record): OrderResult {
        return OrderResult(
            id = record.get(o.ID)!!.toLong(),
            uid = bytesToUuid(record.get(o.UID)!!),
            userUid = bytesToUuid(record.get(o.USER_UID)!!),
            orderNumber = record.get(o.ORDER_NUMBER)!!,
            status = record.get(o.STATUS)!!,
            currency = record.get(o.CURRENCY)!!,
            itemsAmount = record.get(o.ITEMS_AMOUNT)!!,
            discountAmount = record.get(o.DISCOUNT_AMOUNT)!!,
            shippingAmount = record.get(o.SHIPPING_AMOUNT)!!,
            taxAmount = record.get(o.TAX_AMOUNT)!!,
            totalAmount = record.get(o.TOTAL_AMOUNT)!!,
            idempotencyKey = record.get(o.IDEMPOTENCY_KEY)!!,
            paymentUid = record.get(paymentUidField)?.let { bytesToUuid(it) },
            createdAt = record.get(o.CREATED_AT)!!,
            updatedAt = record.get(o.UPDATED_AT),
        )
    }
}
