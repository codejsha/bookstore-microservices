package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.OrderAdjustmentApi
import com.codejsha.bookstore.generated.application.port.openapi.api.OrderApi
import com.codejsha.bookstore.generated.application.port.openapi.api.OrderItemApi
import com.codejsha.bookstore.generated.application.port.openapi.api.OrderShippingApi
import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAdjustmentEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.OrderItemEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderShippingEntity
import com.codejsha.bookstore.order.domain.model.command.*
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.order.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.order.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.order.infrastructure.support.auth.Principal
import com.codejsha.bookstore.order.infrastructure.support.auth.ROLE_ADMIN
import com.codejsha.bookstore.order.infrastructure.support.auth.assertOrderOwner
import com.codejsha.bookstore.order.infrastructure.support.auth.subjectUserUid
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.math.BigDecimal
import java.net.URI
import java.time.ZoneOffset
import java.util.*

@RestController
class OrderController(
    private val orderUseCase: OrderUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : OrderApi, OrderItemApi, OrderShippingApi, OrderAdjustmentApi {

    // ─── OrderApi ───────────────────────────────────────────────────────────

    override fun ordersGetAll(
        userUid: String?, status: OrderStatus?, pageable: Pageable?
    ): ResponseEntity<OrderFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val ownerFilter: UUID? =
            if (principal.hasRole(ROLE_ADMIN)) userUid?.let { UUID.fromString(it) }
            else principal.subjectUserUid()
        val option = OrderQueryOption(userUid = ownerFilter, status = status?.value)
        val context = buildContext(principal)
        val result = orderUseCase.findAllOrders(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            OrderFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toOrderFindResponse(it) }
            )
        )
    }

    override fun ordersPlace(requestBody: PlaceOrderRequest): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        if (!principal.hasRole(ROLE_ADMIN)) {
            throw ForbiddenException(
                "direct order placement is admin-only; use cart checkout for server-side pricing",
            )
        }
        val ownerUserUid: UUID = UUID.fromString(requestBody.userUid)

        val command = OrderCreateCommand(
            userUid = ownerUserUid,
            currency = requestBody.currency,
            itemsAmount = requestBody.itemsAmount.toBigDecimal(),
            discountAmount = requestBody.discountAmount?.toBigDecimal() ?: BigDecimal.ZERO,
            shippingAmount = requestBody.shippingAmount?.toBigDecimal() ?: BigDecimal.ZERO,
            taxAmount = requestBody.taxAmount?.toBigDecimal() ?: BigDecimal.ZERO,
            totalAmount = requestBody.totalAmount.toBigDecimal(),
            idempotencyKey = requestBody.idempotencyKey,
        )

        val items = requestBody.items.map {
            OrderItemCreateCommand(
                productId = it.productId,
                sku = it.sku,
                productName = it.productName,
                options = it.options,
                quantity = it.quantity,
                currency = it.currency,
                price = it.price.toBigDecimal(),
                taxRate = it.taxRate?.toBigDecimal() ?: BigDecimal.ZERO,
            )
        }

        val shipping = requestBody.shipping?.let {
            OrderShippingCreateCommand(
                recipientName = it.recipientName,
                recipientPhone = it.recipientPhone,
                addressLine1 = it.addressLine1,
                addressLine2 = it.addressLine2,
                city = it.city,
                state = it.state,
                postalCode = it.postalCode,
                country = it.country,
                shippingMethod = it.shippingMethod,
            )
        }

        val order = orderUseCase.placeOrder(command, items, shipping, context)
        ResponseEntity.created(URI.create("/api/v1/orders/${order.uid}")).build()
    }

    override fun ordersRead(uid: String): ResponseEntity<OrderFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        val order = orderUseCase.findOrder(UUID.fromString(uid), context)
        principal.assertOrderOwner(order.userUid)
        ResponseEntity.ok(toOrderFindResponse(order))
    }

    override fun ordersCancel(uid: String): ResponseEntity<OrderFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        principal.assertOrderOwner(orderUseCase.findOrder(UUID.fromString(uid), context).userUid)
        val order = orderUseCase.cancelOrder(UUID.fromString(uid), context)
        ResponseEntity.ok(toOrderFindResponse(order))
    }

    // ─── OrderItemApi ───────────────────────────────────────────────────────

    override fun orderItemsAdd(orderUid: String, requestBody: OrderItemCreateRequest): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        val command = OrderItemCreateCommand(
            productId = requestBody.productId,
            sku = requestBody.sku,
            productName = requestBody.productName,
            options = requestBody.options,
            quantity = requestBody.quantity,
            currency = requestBody.currency,
            price = requestBody.price.toBigDecimal(),
            taxRate = requestBody.taxRate?.toBigDecimal() ?: BigDecimal.ZERO,
        )
        val item = orderUseCase.addItem(UUID.fromString(orderUid), command, context)
        ResponseEntity.created(URI.create("/api/v1/orders/$orderUid/items/${item.uid}")).build()
    }

    override fun orderItemsUpdate(
        orderUid: String,
        uid: String,
        requestBody: OrderItemUpdateRequest
    ): ResponseEntity<OrderItemFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        val command = OrderItemUpdateCommand(
            productId = requestBody.productId,
            sku = requestBody.sku,
            productName = requestBody.productName,
            options = requestBody.options,
            quantity = requestBody.quantity,
            currency = requestBody.currency,
            price = requestBody.price?.toBigDecimal(),
            taxRate = requestBody.taxRate?.toBigDecimal(),
        )
        val item = orderUseCase.updateItem(UUID.fromString(orderUid), UUID.fromString(uid), command, context)
        ResponseEntity.ok(toOrderItemFindResponse(item))
    }

    override fun orderItemsRemove(orderUid: String, uid: String): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        orderUseCase.removeItem(UUID.fromString(orderUid), UUID.fromString(uid), context)
        ResponseEntity.noContent().build()
    }

    // ─── OrderShippingApi ───────────────────────────────────────────────────

    override fun orderShippingSet(
        orderUid: String,
        requestBody: OrderShippingCreateRequest
    ): ResponseEntity<OrderShippingFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        val command = OrderShippingCreateCommand(
            recipientName = requestBody.recipientName,
            recipientPhone = requestBody.recipientPhone,
            addressLine1 = requestBody.addressLine1,
            addressLine2 = requestBody.addressLine2,
            city = requestBody.city,
            state = requestBody.state,
            postalCode = requestBody.postalCode,
            country = requestBody.country,
            shippingMethod = requestBody.shippingMethod,
        )
        val shipping = orderUseCase.setShipping(UUID.fromString(orderUid), command, context)
        ResponseEntity.ok(toOrderShippingFindResponse(shipping))
    }

    override fun orderShippingUpdate(
        orderUid: String,
        requestBody: OrderShippingUpdateRequest
    ): ResponseEntity<OrderShippingFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        val command = OrderShippingUpdateCommand(
            recipientName = requestBody.recipientName,
            recipientPhone = requestBody.recipientPhone,
            addressLine1 = requestBody.addressLine1,
            addressLine2 = requestBody.addressLine2,
            city = requestBody.city,
            state = requestBody.state,
            postalCode = requestBody.postalCode,
            country = requestBody.country,
            shippingMethod = requestBody.shippingMethod,
        )
        val shipping = orderUseCase.updateShipping(UUID.fromString(orderUid), command, context)
        ResponseEntity.ok(toOrderShippingFindResponse(shipping))
    }

    // ─── OrderAdjustmentApi ─────────────────────────────────────────────────

    override fun orderAdjustmentsApply(
        orderUid: String,
        requestBody: OrderAdjustmentCreateRequest
    ): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        val command = OrderAdjustmentCreateCommand(
            type = requestBody.type.value,
            label = requestBody.label,
            amount = requestBody.amount.toBigDecimal(),
            meta = null,
        )
        val adj = orderUseCase.applyAdjustment(UUID.fromString(orderUid), command, context)
        ResponseEntity.created(URI.create("/api/v1/orders/$orderUid/adjustments/${adj.uid}")).build()
    }

    override fun orderAdjustmentsRemove(orderUid: String, uid: String): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        authorizeOrder(orderUid, principal, context)
        orderUseCase.removeAdjustment(UUID.fromString(orderUid), UUID.fromString(uid), context)
        ResponseEntity.noContent().build()
    }

    private suspend fun authorizeOrder(orderUid: String, principal: Principal, context: ActorContext) {
        val order = orderUseCase.findOrder(UUID.fromString(orderUid), context)
        principal.assertOrderOwner(order.userUid)
    }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
