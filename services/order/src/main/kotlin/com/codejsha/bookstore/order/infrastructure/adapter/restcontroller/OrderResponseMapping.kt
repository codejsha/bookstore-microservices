package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.order.domain.aggregate.OrderAdjustmentEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.OrderItemEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderShippingEntity
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import java.time.ZoneOffset

internal fun toOrderFindResponse(agg: OrderAggregate) = OrderFindResponse(
    uid = agg.uid.toString(),
    userUid = agg.userUid.toString(),
    orderNumber = agg.orderNumber,
    status = OrderStatus.fromValue(agg.status.value),
    currency = agg.currency,
    itemsAmount = agg.itemsAmount.toDouble(),
    discountAmount = agg.discountAmount.toDouble(),
    shippingAmount = agg.shippingAmount.toDouble(),
    taxAmount = agg.taxAmount.toDouble(),
    totalAmount = agg.totalAmount.toDouble(),
    idempotencyKey = agg.idempotencyKey,
    items = agg.items.takeIf { it.isNotEmpty() }?.map { toOrderItemFindResponse(it) },
    shipping = agg.shipping?.let { toOrderShippingFindResponse(it) },
    adjustments = agg.adjustments.takeIf { it.isNotEmpty() }?.map { toOrderAdjustmentFindResponse(it) },
    createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
    updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
)

internal fun toOrderItemFindResponse(entity: OrderItemEntity) = OrderItemFindResponse(
    uid = entity.uid.toString(),
    productId = entity.productId,
    sku = entity.sku,
    productName = entity.productName,
    options = entity.options,
    quantity = entity.quantity,
    currency = entity.currency,
    price = entity.price.toDouble(),
    taxRate = entity.taxRate.toDouble(),
    subtotal = entity.subtotal.toDouble(),
    createdAt = entity.createdAt.atOffset(ZoneOffset.UTC),
    updatedAt = entity.updatedAt?.atOffset(ZoneOffset.UTC),
)

internal fun toOrderShippingFindResponse(entity: OrderShippingEntity) = OrderShippingFindResponse(
    uid = entity.uid.toString(),
    recipientName = entity.recipientName,
    recipientPhone = entity.recipientPhone,
    addressLine1 = entity.addressLine1,
    addressLine2 = entity.addressLine2,
    city = entity.city,
    state = entity.state,
    postalCode = entity.postalCode,
    country = entity.country,
    shippingMethod = entity.shippingMethod,
    createdAt = entity.createdAt.atOffset(ZoneOffset.UTC),
    updatedAt = entity.updatedAt?.atOffset(ZoneOffset.UTC),
)

internal fun toOrderAdjustmentFindResponse(entity: OrderAdjustmentEntity) = OrderAdjustmentFindResponse(
    uid = entity.uid.toString(),
    type = OrderAdjustmentType.fromValue(entity.type.value),
    label = entity.label,
    amount = entity.amount.toDouble(),
    createdAt = entity.createdAt.atOffset(ZoneOffset.UTC),
    updatedAt = entity.updatedAt?.atOffset(ZoneOffset.UTC),
)

internal fun OrderShippingCreateRequest.toCommand() = OrderShippingCreateCommand(
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
