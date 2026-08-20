package com.codejsha.bookstore.order.domain.workflow

import io.temporal.activity.ActivityInterface
import io.temporal.activity.ActivityMethod

@ActivityInterface
interface OrderActivities {
    @ActivityMethod
    fun createOrder(request: CreateOrderRequest): CreateOrderResult

    @ActivityMethod
    fun confirmOrder(orderUid: String)

    @ActivityMethod
    fun cancelOrder(orderUid: String)

    @ActivityMethod
    fun recordPayment(orderUid: String, paymentUid: String)

    @ActivityMethod
    fun loadCancellationState(orderUid: String): OrderCancellationState

    @ActivityMethod
    fun loadOrderItems(orderUid: String): List<StockReservationItem>

    @ActivityMethod
    fun markOrderRefunded(orderUid: String)

    @ActivityMethod
    fun markOrderShipped(orderUid: String)

    @ActivityMethod
    fun restoreOrderPaid(orderUid: String)
}
