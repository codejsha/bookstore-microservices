package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.Dashboard
import com.codejsha.bookstore.admin.domain.model.Metric
import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.external.OrderLine
import com.codejsha.bookstore.admin.domain.model.external.OrderShipping
import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementDetailLine
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.StockAtWarehouse
import com.codejsha.bookstore.admin.domain.model.external.Subject
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.external.Warehouse
import com.codejsha.bookstore.admin.domain.model.external.Work
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminAuthorItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminDashboardResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminIdentityResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminMetric
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderLineItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderShipping
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminPaymentResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRefundResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementDetailItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementRunResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminStockResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminStockWarehouseItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSubjectItem
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminUserResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWarehouseResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkResponse

internal fun toAdminIdentityResponse(principal: Principal) = AdminIdentityResponse(
    uid = principal.sub,
    email = principal.email,
    name = principal.name,
    roles = principal.roles,
)

internal fun toAdminDashboardResponse(dashboard: Dashboard) = AdminDashboardResponse(
    works = toAdminMetric(dashboard.works),
    orders = toAdminMetric(dashboard.orders),
    warehouses = toAdminMetric(dashboard.warehouses),
)

internal fun toAdminMetric(metric: Metric) = AdminMetric(
    count = metric.count,
    available = metric.available,
)

internal fun toAdminAuthorItem(author: Author) = AdminAuthorItem(
    uid = author.uid,
    name = author.name,
)

internal fun toAdminSubjectItem(subject: Subject) = AdminSubjectItem(
    uid = subject.uid,
    name = subject.name,
)

internal fun toAdminWorkResponse(work: Work) = AdminWorkResponse(
    uid = work.uid,
    title = work.title,
    description = work.description,
    firstPublishDate = work.firstPublishDate,
    olKey = work.olKey,
    authors = work.authors.map { toAdminAuthorItem(it) },
    subjects = work.subjects.map { toAdminSubjectItem(it) },
    createdAt = work.createdAt,
    updatedAt = work.updatedAt,
)

internal fun toAdminUserResponse(user: User) = AdminUserResponse(
    uid = user.uid,
    email = user.email,
    firstName = user.firstName,
    lastName = user.lastName,
    phone = user.phone,
    status = user.status,
    roles = user.roles,
    lastLoginAt = user.lastLoginAt,
    createdAt = user.createdAt,
    updatedAt = user.updatedAt,
)

internal fun toAdminWarehouseResponse(warehouse: Warehouse) = AdminWarehouseResponse(
    uid = warehouse.uid,
    name = warehouse.name,
    address = warehouse.address,
    capacity = warehouse.capacity,
    createdAt = warehouse.createdAt,
    updatedAt = warehouse.updatedAt,
)

internal fun toAdminStockResponse(stock: Stock) = AdminStockResponse(
    uid = stock.uid,
    editionUid = stock.editionUid,
    totalQuantity = stock.totalQuantity,
    warehouses = stock.warehouses.map { toAdminStockWarehouseItem(it) },
)

internal fun toAdminStockWarehouseItem(stock: StockAtWarehouse) = AdminStockWarehouseItem(
    warehouseUid = stock.warehouseUid,
    warehouseName = stock.warehouseName,
    quantity = stock.quantity,
)

internal fun toAdminOrderItem(order: Order) = AdminOrderItem(
    uid = order.uid,
    orderNumber = order.orderNumber,
    userUid = order.userUid,
    status = order.status,
    currency = order.currency,
    itemsAmount = order.itemsAmount,
    discountAmount = order.discountAmount,
    shippingAmount = order.shippingAmount,
    taxAmount = order.taxAmount,
    totalAmount = order.totalAmount,
    createdAt = order.createdAt,
    updatedAt = order.updatedAt,
)

internal fun toAdminOrderResponse(order: Order) = AdminOrderResponse(
    uid = order.uid,
    orderNumber = order.orderNumber,
    userUid = order.userUid,
    status = order.status,
    currency = order.currency,
    itemsAmount = order.itemsAmount,
    discountAmount = order.discountAmount,
    shippingAmount = order.shippingAmount,
    taxAmount = order.taxAmount,
    totalAmount = order.totalAmount,
    items = order.lines.takeIf { it.isNotEmpty() }?.map { toAdminOrderLineItem(it) },
    shipping = order.shipping?.let { toAdminOrderShipping(it) },
    createdAt = order.createdAt,
    updatedAt = order.updatedAt,
)

internal fun toAdminOrderLineItem(line: OrderLine) = AdminOrderLineItem(
    uid = line.uid,
    productId = line.productId,
    productName = line.productName,
    sku = line.sku,
    quantity = line.quantity,
    price = line.price,
    subtotal = line.subtotal,
    taxRate = line.taxRate,
    currency = line.currency,
    options = line.options,
)

internal fun toAdminOrderShipping(shipping: OrderShipping) = AdminOrderShipping(
    recipientName = shipping.recipientName,
    recipientPhone = shipping.recipientPhone,
    addressLine1 = shipping.addressLine1,
    addressLine2 = shipping.addressLine2,
    city = shipping.city,
    state = shipping.state,
    postalCode = shipping.postalCode,
    country = shipping.country,
    shippingMethod = shipping.shippingMethod,
)

internal fun toAdminPaymentResponse(payment: Payment) = AdminPaymentResponse(
    uid = payment.uid,
    paymentId = payment.paymentId,
    customerId = payment.customerId,
    status = payment.status,
    amount = payment.amount,
    amountCaptured = payment.amountCaptured,
    amountCapturable = payment.amountCapturable,
    currency = payment.currency,
    paymentMethod = payment.paymentMethod,
    connector = payment.connector,
    errorCode = payment.errorCode,
    errorMessage = payment.errorMessage,
    confirmedAt = payment.confirmedAt,
    capturedAt = payment.capturedAt,
    cancelledAt = payment.cancelledAt,
    createdAt = payment.createdAt,
    updatedAt = payment.updatedAt,
)

internal fun toAdminRefundResponse(refund: Refund) = AdminRefundResponse(
    uid = refund.uid,
    refundId = refund.refundId,
    paymentId = refund.paymentId,
    status = refund.status,
    refundType = refund.refundType,
    amount = refund.amount,
    currency = refund.currency,
    reason = refund.reason,
    connector = refund.connector,
    errorCode = refund.errorCode,
    errorMessage = refund.errorMessage,
    createdAt = refund.createdAt,
    updatedAt = refund.updatedAt,
)

internal fun toAdminSettlementItem(bucket: SettlementBucket) = AdminSettlementItem(
    uid = bucket.uid,
    settlementDate = bucket.settlementDate,
    currency = bucket.currency,
    paymentMethod = bucket.paymentMethod,
    grossAmount = bucket.grossAmount,
    refundAmount = bucket.refundAmount,
    feeAmount = bucket.feeAmount,
    netAmount = bucket.netAmount,
    paymentCount = bucket.paymentCount,
    refundCount = bucket.refundCount,
    status = bucket.status,
    createdAt = bucket.createdAt,
    updatedAt = bucket.updatedAt,
)

internal fun toAdminSettlementResponse(bucket: SettlementBucket) = AdminSettlementResponse(
    uid = bucket.uid,
    settlementDate = bucket.settlementDate,
    currency = bucket.currency,
    paymentMethod = bucket.paymentMethod,
    grossAmount = bucket.grossAmount,
    refundAmount = bucket.refundAmount,
    feeAmount = bucket.feeAmount,
    netAmount = bucket.netAmount,
    paymentCount = bucket.paymentCount,
    refundCount = bucket.refundCount,
    status = bucket.status,
    details = bucket.details.takeIf { it.isNotEmpty() }?.map { toAdminSettlementDetailItem(it) },
    createdAt = bucket.createdAt,
    updatedAt = bucket.updatedAt,
)

internal fun toAdminSettlementDetailItem(line: SettlementDetailLine) = AdminSettlementDetailItem(
    uid = line.uid,
    sourceType = line.sourceType,
    sourceId = line.sourceId,
    paymentId = line.paymentId,
    amount = line.amount,
    currency = line.currency,
    paymentMethod = line.paymentMethod,
    occurredAt = line.occurredAt,
)

internal fun toAdminSettlementRunResponse(ack: SettlementRunAck) = AdminSettlementRunResponse(
    uid = ack.uid,
    targetDate = ack.targetDate,
    status = ack.status,
    startedAt = ack.startedAt,
)
