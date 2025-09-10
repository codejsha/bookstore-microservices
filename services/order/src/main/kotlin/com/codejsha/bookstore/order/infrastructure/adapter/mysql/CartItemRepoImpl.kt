package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.order.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.order.application.port.repo.CartItemRepo
import com.codejsha.bookstore.order.application.port.repo.CartItemResult
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CART_ITEM
import com.codejsha.platform.shared.data.ActorContext
import org.jooq.DSLContext
import org.jooq.Record
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class CartItemRepoImpl(
    private val dslContext: DSLContext
) : CartItemRepo {
    private val ci = CART_ITEM.`as`("ci")

    override fun findAllByCart(cartId: Long, context: ActorContext): List<CartItemResult> {
        return dslContext
            .select(ci.asterisk())
            .from(ci)
            .where(ci.CART_ID.eq(cartId))
            .orderBy(ci.CREATED_AT.asc())
            .fetch { toCartItemResult(it) }
    }

    override fun findOne(cartId: Long, uid: UUID, context: ActorContext): CartItemResult {
        return dslContext
            .select(ci.asterisk())
            .from(ci)
            .where(ci.CART_ID.eq(cartId).and(ci.UID.eq(uuidToBytes(uid))))
            .fetchOne { toCartItemResult(it) }
            ?: throw NoSuchElementException("CartItem with uid $uid not found")
    }

    override fun findByProduct(cartId: Long, productId: Long, context: ActorContext): CartItemResult? {
        return dslContext
            .select(ci.asterisk())
            .from(ci)
            .where(ci.CART_ID.eq(cartId).and(ci.PRODUCT_ID.eq(productId)))
            .fetchOne { toCartItemResult(it) }
    }

    override fun create(cartId: Long, command: CartAddItemCommand, context: ActorContext): CartItemResult {
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(ci)
            .set(ci.UID, uuidToBytes(uid))
            .set(ci.CART_ID, cartId)
            .set(ci.PRODUCT_ID, command.productId)
            .set(ci.PRODUCT_NAME, command.productName)
            .set(ci.QUANTITY, command.quantity)
            .set(ci.CURRENCY, command.currency)
            .set(ci.PRICE, command.price)
            .set(ci.CREATED_AT, LocalDateTime.now())
            .execute()

        return findOne(cartId, uid, context)
    }

    override fun updateQuantity(cartId: Long, uid: UUID, quantity: Int, context: ActorContext): CartItemResult {
        val updated = dslContext
            .update(ci)
            .set(ci.QUANTITY, quantity)
            .set(ci.UPDATED_AT, LocalDateTime.now())
            .where(ci.CART_ID.eq(cartId).and(ci.UID.eq(uuidToBytes(uid))))
            .execute()

        check(updated > 0) { "CartItem with uid $uid not found" }
        return findOne(cartId, uid, context)
    }

    override fun delete(cartId: Long, uid: UUID, context: ActorContext) {
        dslContext
            .deleteFrom(ci)
            .where(ci.CART_ID.eq(cartId).and(ci.UID.eq(uuidToBytes(uid))))
            .execute()
    }

    override fun deleteAllByCart(cartId: Long, context: ActorContext) {
        dslContext
            .deleteFrom(ci)
            .where(ci.CART_ID.eq(cartId))
            .execute()
    }

    private fun toCartItemResult(record: Record): CartItemResult {
        return CartItemResult(
            id = record.get(ci.ID)!!.toLong(),
            uid = bytesToUuid(record.get(ci.UID)!!),
            cartId = record.get(ci.CART_ID)!!,
            productId = record.get(ci.PRODUCT_ID)!!,
            productName = record.get(ci.PRODUCT_NAME),
            quantity = record.get(ci.QUANTITY)!!,
            currency = record.get(ci.CURRENCY)!!,
            price = record.get(ci.PRICE)!!,
            createdAt = record.get(ci.CREATED_AT)!!,
            updatedAt = record.get(ci.UPDATED_AT),
        )
    }
}
