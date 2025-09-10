package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.order.infrastructure.support.utils.parseJson
import com.codejsha.bookstore.order.infrastructure.support.utils.toJson
import com.codejsha.bookstore.order.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.order.application.port.repo.OrderAdjustmentRepo
import com.codejsha.bookstore.order.application.port.repo.OrderAdjustmentResult
import com.codejsha.bookstore.order.domain.model.command.OrderAdjustmentCreateCommand
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_ADJUSTMENT
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
class OrderAdjustmentRepoImpl(
    private val dslContext: DSLContext
) : OrderAdjustmentRepo {
    private val oa = ORDER_ADJUSTMENT.`as`("oa")
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
    ): Page<OrderAdjustmentResult> {
        val orderId = resolveOrderId(orderUid)

        val content = dslContext
            .select(oa.asterisk())
            .from(oa)
            .where(oa.ORDER_ID.eq(orderId).and(oa.DELETED_AT.isNull))
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toOrderAdjustmentResult(it) }

        val total = dslContext
            .selectCount()
            .from(oa)
            .where(oa.ORDER_ID.eq(orderId).and(oa.DELETED_AT.isNull))
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(orderUid: UUID, uid: UUID, context: ActorContext): OrderAdjustmentResult {
        val orderId = resolveOrderId(orderUid)

        return dslContext
            .select(oa.asterisk())
            .from(oa)
            .where(
                oa.ORDER_ID.eq(orderId)
                    .and(oa.UID.eq(uuidToBytes(uid)))
                    .and(oa.DELETED_AT.isNull)
            )
            .fetchOne { toOrderAdjustmentResult(it) }
            ?: throw NoSuchElementException("OrderAdjustment with uid $uid not found")
    }

    override fun create(
        orderUid: UUID, command: OrderAdjustmentCreateCommand, context: ActorContext
    ): OrderAdjustmentResult {
        val orderId = resolveOrderId(orderUid)
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(oa)
            .set(oa.UID, uuidToBytes(uid))
            .set(oa.ORDER_ID, orderId)
            .set(oa.TYPE, command.type)
            .set(oa.LABEL, command.label)
            .set(oa.AMOUNT, command.amount)
            .set(oa.META, toJson(command.meta))
            .set(oa.CREATED_AT, LocalDateTime.now())
            .set(oa.ACTOR_ID, context.actorId)
            .execute()

        return findOne(orderUid, uid, context)
    }

    override fun delete(orderUid: UUID, uid: UUID, context: ActorContext) {
        val orderId = resolveOrderId(orderUid)

        dslContext
            .update(oa)
            .set(oa.DELETED_AT, LocalDateTime.now())
            .set(oa.ACTOR_ID, context.actorId)
            .where(
                oa.ORDER_ID.eq(orderId)
                    .and(oa.UID.eq(uuidToBytes(uid)))
                    .and(oa.DELETED_AT.isNull)
            )
            .execute()
    }

    private fun toOrderAdjustmentResult(record: Record): OrderAdjustmentResult {
        return OrderAdjustmentResult(
            id = record.get(oa.ID)!!.toLong(),
            uid = bytesToUuid(record.get(oa.UID)!!),
            orderId = record.get(oa.ORDER_ID)!!,
            type = record.get(oa.TYPE)!!,
            label = record.get(oa.LABEL),
            amount = record.get(oa.AMOUNT)!!,
            meta = record.get(oa.META)?.let { parseJson(it) },
            createdAt = record.get(oa.CREATED_AT)!!,
            updatedAt = record.get(oa.UPDATED_AT),
        )
    }
}
