package com.codejsha.bookstore.order.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CART
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CART_ITEM
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.order.support.OrderTestFixtures
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.dao.DuplicateKeyException
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

class CartRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var cartRepo: CartRepoImpl

    @Autowired
    private lateinit var cartItemRepo: CartItemRepoImpl

    private val userUid: UUID = OrderTestFixtures.USER_UID
    private val otherUser: UUID = UUID.fromString("88888888-8888-8888-8888-888888888888")

    @Test
    fun `createForUser_noCartYet_persistsCartForThatUser`() {
        val created = cartRepo.createForUser(userUid, ctx)

        val found = assertNotNull(cartRepo.findByUser(userUid, ctx))
        assertTrue(found.id > 0L)
        assertEquals(created.uid, found.uid)
        assertEquals(userUid, found.userUid)
        assertNull(found.updatedAt)
    }

    @Test
    fun `createForUser_cartAlreadyExists_rejectsWithUniqueViolation`() {
        cartRepo.createForUser(userUid, ctx)

        val ex = assertFailsWith<DuplicateKeyException> { cartRepo.createForUser(userUid, ctx) }

        assertEquals(UNIQUE_VIOLATION, sqlStateOf(ex))
        assertEquals(1, dsl.fetchCount(CART))
    }

    @Test
    fun `findByUser_userWithoutCart_yieldsNull`() {
        cartRepo.createForUser(otherUser, ctx)

        assertNull(cartRepo.findByUser(userUid, ctx))
    }

    @Test
    fun `delete_cartExists_removesOnlyThatUsersCart`() {
        cartRepo.createForUser(userUid, ctx)
        val kept = cartRepo.createForUser(otherUser, ctx)

        cartRepo.delete(userUid, ctx)

        assertNull(cartRepo.findByUser(userUid, ctx))
        assertEquals(kept.uid, cartRepo.findByUser(otherUser, ctx)?.uid)
        assertEquals(1, dsl.fetchCount(CART))
    }

    @Test
    fun `itemCreate_commandGiven_persistsItemInThatCart`() {
        val cart = cartRepo.createForUser(userUid, ctx)

        val created = cartItemRepo.create(cart.id, itemCommand(quantity = 3), ctx)

        val found = cartItemRepo.findOne(cart.id, created.uid, ctx)
        assertEquals(cart.id, found.cartId)
        assertEquals(200L, found.productId)
        assertEquals("Book A", found.productName)
        assertEquals(3, found.quantity)
        assertEquals("KRW", found.currency)
        assertEquals(0, BigDecimal("5000.00").compareTo(found.price))
        assertNull(found.updatedAt)
    }

    @Test
    fun `itemCreate_productNameNull_persistsNull`() {
        val cart = cartRepo.createForUser(userUid, ctx)

        val created = cartItemRepo.create(cart.id, itemCommand(productName = null), ctx)

        assertNull(created.productName)
    }

    @Test
    fun `itemCreate_sameProductTwiceInOneCart_rejectsWithUniqueViolation`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        cartItemRepo.create(cart.id, itemCommand(productId = 200L), ctx)

        val ex = assertFailsWith<DuplicateKeyException> {
            cartItemRepo.create(cart.id, itemCommand(productId = 200L), ctx)
        }

        assertEquals(UNIQUE_VIOLATION, sqlStateOf(ex))
        assertEquals(1, dsl.fetchCount(CART_ITEM))
    }

    @Test
    fun `itemCreate_sameProductInAnotherCart_persistsBoth`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        cartItemRepo.create(cart.id, itemCommand(productId = 200L), ctx)

        cartItemRepo.create(otherCart.id, itemCommand(productId = 200L), ctx)

        assertEquals(2, dsl.fetchCount(CART_ITEM))
    }

    @Test
    fun `itemFindOne_itemBelongsToAnotherCart_throwsNoSuchElementException`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        val item = cartItemRepo.create(cart.id, itemCommand(), ctx)

        assertFailsWith<NoSuchElementException> { cartItemRepo.findOne(otherCart.id, item.uid, ctx) }
    }

    @Test
    fun `itemFindByProduct_productInCart_findsThatItem`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        cartItemRepo.create(cart.id, itemCommand(productId = 100L), ctx)
        val wanted = cartItemRepo.create(cart.id, itemCommand(productId = 200L), ctx)

        assertEquals(wanted.uid, cartItemRepo.findByProduct(cart.id, 200L, ctx)?.uid)
    }

    @Test
    fun `itemFindByProduct_productOnlyInAnotherCart_yieldsNull`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        cartItemRepo.create(otherCart.id, itemCommand(productId = 200L), ctx)

        assertNull(cartItemRepo.findByProduct(cart.id, 200L, ctx))
    }

    @Test
    fun `itemFindAllByCart_severalItems_listsThatCartOldestFirst`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        val uids = listOf(3, 1, 2).map { day ->
            val item = cartItemRepo.create(cart.id, itemCommand(productId = day.toLong()), ctx)
            dsl.update(CART_ITEM)
                .set(CART_ITEM.CREATED_AT, LocalDateTime.of(2026, 1, day, 0, 0))
                .where(CART_ITEM.UID.eq(item.uid))
                .execute()
            day to item.uid
        }.sortedBy { it.first }.map { it.second }
        cartItemRepo.create(otherCart.id, itemCommand(productId = 9L), ctx)

        val items = cartItemRepo.findAllByCart(cart.id, ctx)

        assertEquals(uids, items.map { it.uid })
    }

    @Test
    fun `itemUpdateQuantity_itemExists_changesQuantityAndStampsUpdate`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val item = cartItemRepo.create(cart.id, itemCommand(quantity = 1), ctx)

        val updated = cartItemRepo.updateQuantity(cart.id, item.uid, 7, ctx)

        assertEquals(7, updated.quantity)
        assertEquals(0, item.price.compareTo(updated.price))
        assertNotNull(updated.updatedAt)
    }

    @Test
    fun `itemUpdateQuantity_itemBelongsToAnotherCart_throwsAndKeepsQuantity`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        val item = cartItemRepo.create(cart.id, itemCommand(quantity = 1), ctx)

        assertFailsWith<IllegalStateException> { cartItemRepo.updateQuantity(otherCart.id, item.uid, 7, ctx) }
        assertEquals(1, cartItemRepo.findOne(cart.id, item.uid, ctx).quantity)
    }

    @Test
    fun `itemDelete_itemExists_removesOnlyThatItem`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val removed = cartItemRepo.create(cart.id, itemCommand(productId = 100L), ctx)
        val kept = cartItemRepo.create(cart.id, itemCommand(productId = 200L), ctx)

        cartItemRepo.delete(cart.id, removed.uid, ctx)

        assertEquals(listOf(kept.uid), cartItemRepo.findAllByCart(cart.id, ctx).map { it.uid })
        assertEquals(1, dsl.fetchCount(CART_ITEM))
    }

    @Test
    fun `itemDelete_itemBelongsToAnotherCart_leavesItemInPlace`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        val item = cartItemRepo.create(cart.id, itemCommand(), ctx)

        cartItemRepo.delete(otherCart.id, item.uid, ctx)

        assertEquals(item.uid, cartItemRepo.findOne(cart.id, item.uid, ctx).uid)
    }

    @Test
    fun `itemDeleteAllByCart_twoCartsWithItems_emptiesOnlyThatCart`() {
        val cart = cartRepo.createForUser(userUid, ctx)
        val otherCart = cartRepo.createForUser(otherUser, ctx)
        cartItemRepo.create(cart.id, itemCommand(productId = 100L), ctx)
        cartItemRepo.create(cart.id, itemCommand(productId = 200L), ctx)
        val kept = cartItemRepo.create(otherCart.id, itemCommand(productId = 100L), ctx)

        cartItemRepo.deleteAllByCart(cart.id, ctx)

        assertTrue(cartItemRepo.findAllByCart(cart.id, ctx).isEmpty())
        assertEquals(listOf(kept.uid), cartItemRepo.findAllByCart(otherCart.id, ctx).map { it.uid })
    }

    private fun itemCommand(
        productId: Long = 200L,
        productName: String? = "Book A",
        quantity: Int = 1,
    ) = CartAddItemCommand(
        productId = productId,
        productName = productName,
        quantity = quantity,
        currency = "KRW",
        price = BigDecimal("5000.00"),
    )
}
