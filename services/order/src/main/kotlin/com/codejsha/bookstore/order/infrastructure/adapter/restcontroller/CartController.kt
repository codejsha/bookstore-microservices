package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.CartApi
import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.order.application.usecase.CartUseCase
import com.codejsha.bookstore.order.domain.aggregate.CartAggregate
import com.codejsha.bookstore.order.domain.aggregate.CartItemEntity
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.domain.model.command.CartCheckoutCommand
import com.codejsha.bookstore.order.infrastructure.support.auth.BadRequestException
import com.codejsha.bookstore.order.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.order.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.order.infrastructure.support.auth.Principal
import com.codejsha.bookstore.order.infrastructure.support.auth.subjectUserUid
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.math.BigDecimal
import java.util.*

@RestController
class CartController(
    private val cartUseCase: CartUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : CartApi {

    override fun cartGet(): ResponseEntity<CartFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val cart = cartUseCase.getCart(principal.cartOwner(), buildContext())
        ResponseEntity.ok(toCartFindResponse(cart))
    }

    override fun cartAddItem(requestBody: CartAddItemRequest): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val command = CartAddItemCommand(
            productId = requestBody.productId,
            productName = requestBody.productName,
            quantity = requestBody.quantity,
            currency = requestBody.currency,
            price = BigDecimal.valueOf(requestBody.price),
        )
        cartUseCase.addItem(principal.cartOwner(), command, buildContext())
        ResponseEntity.noContent().build()
    }

    override fun cartUpdateItem(
        uid: String,
        requestBody: CartUpdateItemRequest,
    ): ResponseEntity<CartItemFindResponse> = runBlocking {
        val principal = principalResolver.require()
        if (requestBody.quantity < 0) {
            throw BadRequestException("cart item quantity must be zero or positive")
        }
        val item = cartUseCase.updateItemQuantity(
            principal.cartOwner(),
            UUID.fromString(uid),
            requestBody.quantity,
            buildContext(),
        )
        ResponseEntity.ok(toCartItemFindResponse(item))
    }

    override fun cartRemoveItem(uid: String): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        cartUseCase.removeItem(principal.cartOwner(), UUID.fromString(uid), buildContext())
        ResponseEntity.noContent().build()
    }

    override fun cartClear(): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        cartUseCase.clearCart(principal.cartOwner(), buildContext())
        ResponseEntity.noContent().build()
    }

    override fun cartCheckout(requestBody: CartCheckoutRequest): ResponseEntity<OrderFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val command = CartCheckoutCommand(
            currency = requestBody.currency,
            idempotencyKey = requestBody.idempotencyKey,
            shipping = requestBody.shipping?.toCommand(),
        )
        val order = cartUseCase.checkout(principal.cartOwner(), command, buildContext())
        ResponseEntity.ok(toOrderFindResponse(order))
    }

    private fun Principal.cartOwner(): UUID =
        subjectUserUid() ?: throw ForbiddenException("subject is not a platform user uid")

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toCartFindResponse(agg: CartAggregate) = CartFindResponse(
        uid = agg.uid.toString(),
        userUid = agg.userUid.toString(),
        items = agg.items.map { toCartItemFindResponse(it) },
        totalItems = agg.items.sumOf { it.quantity },
        totalAmount = agg.items
            .fold(BigDecimal.ZERO) { acc, item -> acc + item.price.multiply(BigDecimal(item.quantity)) }
            .toDouble(),
    )

    private fun toCartItemFindResponse(entity: CartItemEntity) = CartItemFindResponse(
        uid = entity.uid.toString(),
        productId = entity.productId,
        productName = entity.productName,
        quantity = entity.quantity,
        currency = entity.currency,
        price = entity.price.toDouble(),
        subtotal = entity.price.multiply(BigDecimal(entity.quantity)).toDouble(),
    )

    private fun buildContext() = ActorContext(actorId = 0L, ActorType.USER)
}
