package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.external.OrderLine
import com.codejsha.bookstore.admin.domain.model.external.OrderShipping
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.order.api.OrderApi
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderItemFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderShippingFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderStatus
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

@Component
class OrderRestClient(
    private val orderApi: OrderApi,
) : OrderClient {
    override fun countOrders(): Long =
        orderApi.ordersGetAll(null, null, COUNT_ONLY_PAGE_SIZE, null, null).total

    override fun findAllOrders(option: OrderQueryOption, pageable: Pageable): Page<Order> {
        val response = orderApi.ordersGetAll(
            option.userUid,
            option.status?.toOrderStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toOrder() }, pageable, response.total)
    }

    override fun findOrder(uid: String): Order = orderApi.ordersRead(uid).toOrder()

    override fun cancelOrder(uid: String): Order = orderApi.ordersCancel(uid).toOrder()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun String.toOrderStatus(): OrderStatus =
        parseEnumValue("order status", this, OrderStatus.Companion::fromValue)

    private fun OrderFindResponse.toOrder() = Order(
        uid = uid,
        orderNumber = orderNumber,
        userUid = userUid,
        status = status.value,
        currency = currency,
        itemsAmount = itemsAmount,
        discountAmount = discountAmount,
        shippingAmount = shippingAmount,
        taxAmount = taxAmount,
        totalAmount = totalAmount,
        lines = items.orEmpty().map { it.toLine() },
        shipping = shipping?.toShipping(),
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun OrderItemFindResponse.toLine() = OrderLine(
        uid = uid,
        productId = productId,
        productName = productName,
        sku = sku,
        quantity = quantity,
        price = price,
        subtotal = subtotal,
        taxRate = taxRate,
        currency = currency,
        options = options,
    )

    private fun OrderShippingFindResponse.toShipping() = OrderShipping(
        recipientName = recipientName,
        recipientPhone = recipientPhone,
        addressLine1 = addressLine1,
        addressLine2 = addressLine2,
        city = city,
        state = state,
        postalCode = postalCode,
        country = country,
        shippingMethod = shippingMethod,
    )
}
