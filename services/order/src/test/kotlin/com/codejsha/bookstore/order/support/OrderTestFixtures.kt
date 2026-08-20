package com.codejsha.bookstore.order.support

import com.codejsha.bookstore.order.application.port.repo.CartItemResult
import com.codejsha.bookstore.order.application.port.repo.CartResult
import com.codejsha.bookstore.order.application.port.repo.OrderAdjustmentResult
import com.codejsha.bookstore.order.application.port.repo.OrderItemResult
import com.codejsha.bookstore.order.application.port.repo.OrderResult
import com.codejsha.bookstore.order.application.port.repo.OrderShippingResult
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID

object OrderTestFixtures {

    val FIXED_TIME: LocalDateTime = LocalDateTime.of(2026, 1, 1, 0, 0)
    val DEFAULT_CONTEXT: ActorContext = ActorContext(actorId = 1L, actorType = ActorType.USER)

    val ORDER_UID: UUID = UUID.fromString("11111111-1111-1111-1111-111111111111")
    val ORDER_ITEM_UID: UUID = UUID.fromString("22222222-2222-2222-2222-222222222222")
    val ORDER_ADJUSTMENT_UID: UUID = UUID.fromString("33333333-3333-3333-3333-333333333333")
    val ORDER_SHIPPING_UID: UUID = UUID.fromString("44444444-4444-4444-4444-444444444444")
    val CART_UID: UUID = UUID.fromString("55555555-5555-5555-5555-555555555555")
    val CART_ITEM_UID: UUID = UUID.fromString("66666666-6666-6666-6666-666666666666")
    val USER_UID: UUID = UUID.fromString("77777777-7777-7777-7777-777777777777")

    fun orderResult(
        id: Long = 1L,
        uid: UUID = ORDER_UID,
        userUid: UUID = USER_UID,
        orderNumber: String = "ORD-20260101-0001",
        status: String = OrderStatus.PENDING.value,
        currency: String = "KRW",
        itemsAmount: BigDecimal = BigDecimal("10000.00"),
        totalAmount: BigDecimal = BigDecimal("10000.00"),
        idempotencyKey: String = "idem_001",
        paymentUid: UUID? = null,
    ): OrderResult = OrderResult(
        id = id,
        uid = uid,
        userUid = userUid,
        orderNumber = orderNumber,
        status = status,
        currency = currency,
        itemsAmount = itemsAmount,
        discountAmount = BigDecimal.ZERO,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = totalAmount,
        idempotencyKey = idempotencyKey,
        paymentUid = paymentUid,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )

    fun orderItemResult(
        id: Long = 1L,
        uid: UUID = ORDER_ITEM_UID,
        orderId: Long = 1L,
        productId: Long = 200L,
        productName: String? = "Book A",
        quantity: Int = 2,
        currency: String = "KRW",
        price: BigDecimal = BigDecimal("5000.00"),
        taxRate: BigDecimal = BigDecimal("0.00"),
    ): OrderItemResult = OrderItemResult(
        id = id,
        uid = uid,
        orderId = orderId,
        productId = productId,
        sku = "SKU-001",
        productName = productName,
        options = null,
        quantity = quantity,
        currency = currency,
        price = price,
        taxRate = taxRate,
        subtotal = price * BigDecimal(quantity),
        createdAt = FIXED_TIME,
        updatedAt = null,
    )

    fun orderAdjustmentResult(
        id: Long = 1L,
        uid: UUID = ORDER_ADJUSTMENT_UID,
        orderId: Long = 1L,
        type: String = "COUPON",
        label: String? = "WELCOME",
        amount: BigDecimal = BigDecimal("-1000.00"),
    ): OrderAdjustmentResult = OrderAdjustmentResult(
        id = id,
        uid = uid,
        orderId = orderId,
        type = type,
        label = label,
        amount = amount,
        meta = null,
        createdAt = FIXED_TIME,
        updatedAt = null,
    )

    fun orderShippingResult(
        id: Long = 1L,
        uid: UUID = ORDER_SHIPPING_UID,
        orderId: Long = 1L,
    ): OrderShippingResult = OrderShippingResult(
        id = id,
        uid = uid,
        orderId = orderId,
        recipientName = "Alice",
        recipientPhone = "010-1234-5678",
        addressLine1 = "1 Foo St",
        addressLine2 = null,
        city = "Seoul",
        state = "Seoul",
        postalCode = "12345",
        country = "KR",
        shippingMethod = "STANDARD",
        createdAt = FIXED_TIME,
        updatedAt = null,
    )

    fun cartResult(
        id: Long = 1L,
        uid: UUID = CART_UID,
        userUid: UUID = USER_UID,
    ): CartResult = CartResult(
        id = id,
        uid = uid,
        userUid = userUid,
        createdAt = FIXED_TIME,
        updatedAt = null,
    )

    fun cartItemResult(
        id: Long = 1L,
        uid: UUID = CART_ITEM_UID,
        cartId: Long = 1L,
        productId: Long = 200L,
        productName: String? = "Book A",
        quantity: Int = 1,
        currency: String = "KRW",
        price: BigDecimal = BigDecimal("5000.00"),
    ): CartItemResult = CartItemResult(
        id = id,
        uid = uid,
        cartId = cartId,
        productId = productId,
        productName = productName,
        quantity = quantity,
        currency = currency,
        price = price,
        createdAt = FIXED_TIME,
        updatedAt = null,
    )
}
