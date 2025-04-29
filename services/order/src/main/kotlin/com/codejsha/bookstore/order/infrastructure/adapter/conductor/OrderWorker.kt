package com.codejsha.bookstore.order.infrastructure.adapter.conductor

import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.handler.OrderHandler
import com.codejsha.bookstore.order.domain.model.OrderDto
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.netflix.conductor.common.metadata.tasks.TaskResult
import com.netflix.conductor.sdk.workflow.task.WorkerTask
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.stereotype.Component

@Component
class OrderWorker(
    val orderHandler: OrderHandler,
    val orderUseCase: OrderUseCase
) {

    @WorkerTask(value = "place_order")
    @WithSpan
    fun placeOrder(orderDto: OrderDto): TaskResult {
        val order = orderUseCase.createOrder(orderDto)
        val output = order.block()!!.toMap()
        val result = TaskResult()
        result.outputData = output
        result.status = TaskResult.Status.COMPLETED
        return result
    }

    @WorkerTask(value = "change_order_status")
    @WithSpan
    fun changeStatusToPaid(orderId: Long): TaskResult {
        val orderAgg = orderUseCase.findOrder(orderId).block()
            ?: throw IllegalArgumentException("Order not found")
        val updated = orderAgg.order.changeStatus(OrderStatus.PAID)

        val order = orderHandler.handle(OrderCommand.UpdateOrder(updated))
        val output = order.block()!!.toMap()
        val result = TaskResult()
        result.outputData = output
        result.status = TaskResult.Status.COMPLETED
        return result
    }
}
