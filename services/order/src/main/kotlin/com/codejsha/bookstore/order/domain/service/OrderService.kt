package com.codejsha.bookstore.order.domain.service

import com.codejsha.bookstore.order.application.port.saga.PlaceOrderPayload
import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.application.usecase.OrderItemCommand
import com.codejsha.bookstore.order.application.usecase.OrderRead
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.constant.TaskProperty
import com.codejsha.bookstore.order.domain.constant.Workflow
import com.codejsha.bookstore.order.domain.handler.OrderHandler
import com.codejsha.bookstore.order.domain.handler.OrderItemHandler
import com.codejsha.bookstore.order.domain.model.FilterCondition
import com.codejsha.bookstore.order.domain.model.OrderDto
import com.codejsha.bookstore.service.application.port.pb.bookpb.BookServiceGrpc
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentServiceGrpc
import com.codejsha.bookstore.service.application.port.pb.userpb.UserServiceGrpc
import com.netflix.conductor.client.http.WorkflowClient
import com.netflix.conductor.common.metadata.workflow.StartWorkflowRequest
import org.springframework.stereotype.Service
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono

@Service
class OrderService(
    private val orderHandler: OrderHandler,
    private val orderItemHandler: OrderItemHandler,
    private val bookStub: BookServiceGrpc.BookServiceStub,
    private val paymentStub: PaymentServiceGrpc.PaymentServiceStub,
    private val userStub: UserServiceGrpc.UserServiceStub,
    private val workflowClient: WorkflowClient
) : OrderUseCase {

    override fun startPlaceOrderWorkflow(dto: OrderDto): Mono<Unit> {
        val payload = PlaceOrderPayload.create(dto)
        val request = StartWorkflowRequest()
        request.name = Workflow.BOOKSTORE_ORDER_BOOK_WORKFLOW.wfName
        request.version = 1
        request.correlationId = payload.requestId
        request.input = payload.toMap()
        val taskToDomain: Map<String, String> = mapOf(
            TaskProperty.ALL_WORKERS.propertyName to TaskProperty.ALL_WORKERS.value,
        )
        request.taskToDomain = taskToDomain

        var workflowId: String?
        try {
            workflowId = workflowClient.startWorkflow(request)
        } catch (e: Exception) {
            return Mono.error(e)
        }

        return if (workflowId.isNullOrBlank()) {
            Mono.error(IllegalStateException("Failed to start workflow"))
        } else {
            Mono.just(Unit)
        }
    }

    override fun findAllOrders(cond: FilterCondition): Mono<List<OrderAggregate>> {
        val orderFlux = orderHandler.handle(OrderRead.FindAllOrdersRead(cond)) as Flux<OrderEntity>

        val orderAggFlux = orderFlux.flatMap {
            val orderItemFlux = orderItemHandler.handle(OrderRead.FindOrderRead(requireNotNull(it.id)))
                    as Flux<OrderItemEntity>
            orderItemFlux.collectList().map { orderItems ->
                OrderAggregate.create(
                    id = requireNotNull(it.id), order = it, orderItems = orderItems
                )
            }
        }

        return orderAggFlux.collectList()
    }

    override fun findOrder(id: Long): Mono<OrderAggregate> {
        val orderMono = orderHandler.handle(OrderRead.FindOrderRead(id)) as Mono<OrderEntity>

        val orderAggMono = orderMono.flatMap { order ->
            val orderItemFlux = orderItemHandler.handle(OrderRead.FindOrderRead(requireNotNull(order.id)))
                    as Flux<OrderItemEntity>
            orderItemFlux.collectList().map { items ->
                OrderAggregate.create(
                    id = requireNotNull(order.id), order = order, orderItems = items
                )
            }
        }

        return orderAggMono
    }

    override fun createOrder(dto: OrderDto): Mono<OrderAggregate> {
        val command = OrderCommand.CreateOrderCommand(dto)
        val orderMono = orderHandler.handle(command)

        val aggMono = orderMono.flatMap { order ->
            val orderId = requireNotNull(order.id)
            val orderItems = requireNotNull(dto.orderItems)

            val orderItemCommand = OrderItemCommand.CreateOrderItemsCommand(orderId, orderItems)
            orderItemHandler.handle(orderItemCommand)
                .collectList()
                .map { items -> OrderAggregate.create(orderId, order, items) }
        }

        return aggMono
    }

    override fun updateOrder(dto: OrderDto): Mono<OrderAggregate> {
        val command = OrderCommand.UpdateOrderCommand(dto)
        val orderMono = orderHandler.handle(command)

        val aggMono = orderMono.flatMap { order ->
            val orderId = requireNotNull(order.id)
            val orderItems = requireNotNull(dto.orderItems)

            val orderItemCommand = OrderItemCommand.UpdateOrderItemsCommand(orderId, orderItems)
            orderItemHandler.handle(orderItemCommand)
                .collectList()
                .map { items -> OrderAggregate.create(orderId, order, items) }
        }

        return aggMono
    }

    override fun deleteOrder(id: Long): Mono<Void> {
        val voidMono = orderHandler.handle(OrderCommand.DeleteOrderCommand(id))
            .then<Void>(Mono.empty())
        return voidMono
    }
}
