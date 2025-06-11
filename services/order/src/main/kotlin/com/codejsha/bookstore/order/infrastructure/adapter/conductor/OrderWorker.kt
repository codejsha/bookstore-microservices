package com.codejsha.bookstore.order.infrastructure.adapter.conductor

import com.codejsha.bookstore.order.application.port.saga.ChangeOrderStatusPayload
import com.codejsha.bookstore.order.application.port.saga.PlaceOrderPayload
import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.handler.OrderHandler
import com.codejsha.bookstore.order.domain.model.OrderMapper
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.netflix.conductor.common.metadata.tasks.TaskResult
import com.netflix.conductor.sdk.workflow.task.WorkerTask
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.stereotype.Component

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "command")
class OrderWorker(
    val orderUseCase: OrderUseCase,
    val orderHandler: OrderHandler
) {

    @WorkerTask(value = "place_order")
    @WithSpan
    fun placeOrderTask(payload: PlaceOrderPayload): TaskResult {
        val dto = OrderMapper.toOrderDto(payload)
        val order = orderUseCase.createOrder(dto).block()
            ?: throw IllegalArgumentException("Order creation failed")
        val output: Map<String, Any?> = mapOf(
            "order_id" to order.id,
            "user_id" to order.userId,
            "payment_type" to payload.paymentType,
            "card_number" to payload.cardNumber,
            "amount" to order.totalPrice,
        )

        val result = TaskResult()
        result.outputData = output
        result.status = TaskResult.Status.COMPLETED
        return result
    }

    @WorkerTask(value = "change_order_status_paid")
    @WithSpan
    fun changeOrderStatusPaidTask(payload: ChangeOrderStatusPayload): TaskResult {
        val command = OrderCommand.ChangeOrderStatusCommand(payload.orderId, OrderStatus.PAID)
        orderHandler.handle(command).block()
            ?: throw IllegalArgumentException("Order update failed")
        val order = orderUseCase.findOrder(payload.orderId).block()
            ?: throw IllegalArgumentException("Order not found")

        val result = TaskResult()
        result.outputData = mapOf(
            "order_id" to order.id,
            "order_items" to order.orderItems.map {
                mapOf(
                    "book_id" to it.bookId,
                    "quantity" to it.quantity,
                )
            },
        )
        result.status = TaskResult.Status.COMPLETED
        return result
    }
}
