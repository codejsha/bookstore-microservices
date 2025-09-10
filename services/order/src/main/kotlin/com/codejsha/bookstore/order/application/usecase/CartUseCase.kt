package com.codejsha.bookstore.order.application.usecase

import com.codejsha.bookstore.order.domain.aggregate.CartAggregate
import com.codejsha.bookstore.order.domain.aggregate.CartItemEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.domain.model.command.CartCheckoutCommand
import com.codejsha.platform.shared.data.ActorContext
import java.util.UUID

interface CartUseCase {
    suspend fun getCart(userUid: UUID, context: ActorContext): CartAggregate
    suspend fun addItem(userUid: UUID, command: CartAddItemCommand, context: ActorContext): CartItemEntity
    suspend fun updateItemQuantity(userUid: UUID, itemUid: UUID, quantity: Int, context: ActorContext): CartItemEntity
    suspend fun removeItem(userUid: UUID, itemUid: UUID, context: ActorContext)
    suspend fun clearCart(userUid: UUID, context: ActorContext)
    suspend fun checkout(userUid: UUID, command: CartCheckoutCommand, context: ActorContext): OrderAggregate
}
