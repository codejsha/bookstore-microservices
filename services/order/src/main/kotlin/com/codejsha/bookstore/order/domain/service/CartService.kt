package com.codejsha.bookstore.order.domain.service

import com.codejsha.bookstore.order.application.port.repo.*
import com.codejsha.bookstore.order.application.port.support.DistributedLock
import com.codejsha.bookstore.order.application.port.support.IdempotencyService
import com.codejsha.bookstore.order.application.port.support.TransactionRunner
import com.codejsha.bookstore.order.application.usecase.CartUseCase
import com.codejsha.bookstore.order.domain.aggregate.*
import com.codejsha.bookstore.order.domain.model.CartStateConflictException
import com.codejsha.bookstore.order.domain.model.OrderOwnershipException
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.domain.model.command.CartCheckoutCommand
import com.codejsha.bookstore.order.domain.model.command.InvalidCommandException
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.stereotype.Service
import java.math.BigDecimal
import java.util.*

@Service
class CartService(
    private val cartRepo: CartRepo,
    private val cartItemRepo: CartItemRepo,
    private val orderRepo: OrderRepo,
    private val orderItemRepo: OrderItemRepo,
    private val orderShippingRepo: OrderShippingRepo,
    private val distributedLock: DistributedLock,
    private val idempotencyService: IdempotencyService,
    private val txRunner: TransactionRunner,
) : CartUseCase {

    // ─── Query ──────────────────────────────────────────────────────────────

    @WithSpan
    override suspend fun getCart(userUid: UUID, context: ActorContext): CartAggregate {
        val cart = getOrCreateCart(userUid, context)
        return txRunner.tx {
            val items = cartItemRepo.findAllByCart(cart.id, context).map { it.toEntity() }
            cart.toAggregate().copy(items = items)
        }
    }

    // ─── Item management ────────────────────────────────────────────────────

    @WithSpan
    override suspend fun addItem(
        userUid: UUID, command: CartAddItemCommand, context: ActorContext
    ): CartItemEntity {
        val cart = getOrCreateCart(userUid, context)
        return txRunner.tx {
            val existing = cartItemRepo.findByProduct(cart.id, command.productId, context)
            if (existing != null) {
                cartItemRepo.updateQuantity(cart.id, existing.uid, existing.quantity + command.quantity, context)
                    .toEntity()
            } else {
                cartItemRepo.create(cart.id, command, context).toEntity()
            }
        }
    }

    @WithSpan
    override suspend fun updateItemQuantity(
        userUid: UUID, itemUid: UUID, quantity: Int, context: ActorContext
    ): CartItemEntity {
        if (quantity < 0) {
            throw InvalidCommandException("cart item quantity must be zero or positive, was $quantity")
        }
        val cart = getOrCreateCart(userUid, context)
        return txRunner.tx {
            if (quantity == 0) {
                val item = cartItemRepo.findOne(cart.id, itemUid, context)
                cartItemRepo.delete(cart.id, itemUid, context)
                item.toEntity().copy(quantity = 0)
            } else {
                cartItemRepo.updateQuantity(cart.id, itemUid, quantity, context).toEntity()
            }
        }
    }

    @WithSpan
    override suspend fun removeItem(userUid: UUID, itemUid: UUID, context: ActorContext) {
        val cart = getOrCreateCart(userUid, context)
        txRunner.tx {
            cartItemRepo.delete(cart.id, itemUid, context)
        }
    }

    @WithSpan
    override suspend fun clearCart(userUid: UUID, context: ActorContext) {
        txRunner.tx {
            val cart = cartRepo.findByUser(userUid, context) ?: return@tx
            cartItemRepo.deleteAllByCart(cart.id, context)
        }
    }

    // ─── Checkout ────────────────────────────────────────────────────────────

    @WithSpan
    override suspend fun checkout(
        userUid: UUID, command: CartCheckoutCommand, context: ActorContext
    ): OrderAggregate {
        val idempotencyKey = "$userUid:${command.idempotencyKey}"

        idempotencyService.getResult(idempotencyKey)?.let { existingOrderUid ->
            return txRunner.tx {
                loadOwnedOrder(UUID.fromString(existingOrderUid), userUid, context)
            }
        }

        val lockKey = "checkout:$idempotencyKey"
        val lockToken = distributedLock.tryLock(lockKey)
            ?: throw CartStateConflictException("Another checkout is already in progress for this cart")
        try {
            idempotencyService.getResult(idempotencyKey)?.let { existingOrderUid ->
                return txRunner.tx {
                    loadOwnedOrder(UUID.fromString(existingOrderUid), userUid, context)
                }
            }

            val created = txRunner.tx {
                val cart = cartRepo.findByUser(userUid, context)
                    ?: throw CartStateConflictException("Cart is empty")
                val items = cartItemRepo.findAllByCart(cart.id, context)
                if (items.isEmpty()) throw CartStateConflictException("Cart is empty")

                val itemsAmount = items.fold(BigDecimal.ZERO) { acc, item ->
                    acc + (item.price * BigDecimal(item.quantity))
                }

                val orderCreateCommand = OrderCreateCommand(
                    userUid = userUid,
                    currency = command.currency,
                    itemsAmount = itemsAmount,
                    discountAmount = BigDecimal.ZERO,
                    shippingAmount = BigDecimal.ZERO,
                    taxAmount = BigDecimal.ZERO,
                    totalAmount = itemsAmount,
                    idempotencyKey = idempotencyKey,
                )

                val orderResult = orderRepo.create(orderCreateCommand, context)

                val orderItemEntities = items.map { cartItem ->
                    val orderItemCmd = OrderItemCreateCommand(
                        productId = cartItem.productId,
                        sku = null,
                        productName = cartItem.productName,
                        options = null,
                        quantity = cartItem.quantity,
                        currency = cartItem.currency,
                        price = cartItem.price,
                        taxRate = BigDecimal.ZERO,
                    )
                    orderItemRepo.create(orderResult.uid, orderItemCmd, context).toEntity()
                }

                val shippingEntity = command.shipping?.let {
                    orderShippingRepo.create(orderResult.uid, it, context).toEntity()
                }

                cartItemRepo.deleteAllByCart(cart.id, context)
                cartRepo.delete(userUid, context)

                orderResult.toAggregate().copy(
                    items = orderItemEntities,
                    shipping = shippingEntity,
                )
            }

            val won = idempotencyService.tryMarkProcessed(idempotencyKey, created.uid.toString())
            if (!won) {
                val canonicalUid = idempotencyService.getResult(idempotencyKey)
                if (canonicalUid != null && canonicalUid != created.uid.toString()) {
                    return txRunner.tx {
                        loadOwnedOrder(UUID.fromString(canonicalUid), userUid, context)
                    }
                }
            }
            return created
        } finally {
            distributedLock.unlock(lockKey, lockToken)
        }
    }

    private fun loadOwnedOrder(orderUid: UUID, userUid: UUID, context: ActorContext): OrderAggregate {
        val order = orderRepo.findOne(orderUid, context)
        if (order.userUid != userUid) {
            throw OrderOwnershipException("order does not belong to the requesting user")
        }
        return order.toAggregate()
    }

    private suspend fun getOrCreateCart(userUid: UUID, context: ActorContext): CartResult {
        txRunner.tx { cartRepo.findByUser(userUid, context) }?.let { return it }
        return try {
            txRunner.tx { cartRepo.createForUser(userUid, context) }
        } catch (e: DataIntegrityViolationException) {
            txRunner.tx { cartRepo.findByUser(userUid, context) }
                ?: throw e
        }
    }
}
