package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.order.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.order.application.port.repo.CartRepo
import com.codejsha.bookstore.order.application.port.repo.CartResult
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CART
import com.codejsha.platform.shared.data.ActorContext
import org.jooq.DSLContext
import org.jooq.Record
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class CartRepoImpl(
    private val dslContext: DSLContext
) : CartRepo {
    private val c = CART.`as`("c")

    override fun findByUser(userUid: UUID, context: ActorContext): CartResult? {
        return dslContext
            .select(c.asterisk())
            .from(c)
            .where(c.USER_UID.eq(uuidToBytes(userUid)))
            .fetchOne { toCartResult(it) }
    }

    override fun createForUser(userUid: UUID, context: ActorContext): CartResult {
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(c)
            .set(c.UID, uuidToBytes(uid))
            .set(c.USER_UID, uuidToBytes(userUid))
            .set(c.CREATED_AT, LocalDateTime.now())
            .execute()

        return findByUser(userUid, context)
            ?: throw NoSuchElementException("Cart for user $userUid not found after creation")
    }

    override fun delete(userUid: UUID, context: ActorContext) {
        dslContext
            .deleteFrom(c)
            .where(c.USER_UID.eq(uuidToBytes(userUid)))
            .execute()
    }

    private fun toCartResult(record: Record): CartResult {
        return CartResult(
            id = record.get(c.ID)!!.toLong(),
            uid = bytesToUuid(record.get(c.UID)!!),
            userUid = bytesToUuid(record.get(c.USER_UID)!!),
            createdAt = record.get(c.CREATED_AT)!!,
            updatedAt = record.get(c.UPDATED_AT),
        )
    }
}
