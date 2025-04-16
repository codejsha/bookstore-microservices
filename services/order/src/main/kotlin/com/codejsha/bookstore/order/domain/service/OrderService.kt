package com.codejsha.bookstore.order.domain.service

import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.application.usecase.OrderItemCommand
import com.codejsha.bookstore.order.application.usecase.OrderQuery
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.handler.OrderHandler
import com.codejsha.bookstore.order.domain.handler.OrderItemHandler
import com.codejsha.bookstore.order.domain.model.FilterCondition
import com.codejsha.bookstore.order.domain.model.OrderDto
import com.codejsha.bookstore.service.application.port.pb.bookpb.BookServiceGrpc
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentServiceGrpc
import com.codejsha.bookstore.service.application.port.pb.userpb.UserServiceGrpc
import org.springframework.stereotype.Service
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.core.scheduler.Schedulers

@Service
class OrderService(
    private val orderHandler: OrderHandler,
    private val orderItemHandler: OrderItemHandler,
    private val bookStub: BookServiceGrpc.BookServiceStub,
    private val paymentStub: PaymentServiceGrpc.PaymentServiceStub,
    private val userStub: UserServiceGrpc.UserServiceStub
) : OrderUseCase {

    override fun findAllOrders(cond: FilterCondition): Mono<List<OrderAggregate>> {
        val orderFlux = orderHandler.handle(OrderQuery.FindAllOrders(cond)) as Flux<OrderEntity>
        val orderAgg = orderFlux.publishOn(Schedulers.boundedElastic()).mapNotNull {
            val orderItemFlux =
                orderItemHandler.handle(OrderQuery.FindOrder(it.id!!)) as Flux<OrderItemEntity>
            OrderAggregate(it.id!!, it, orderItemFlux.collectList().block().orEmpty())
        }
        return orderAgg.collectList()
    }

    override fun findOrder(id: Long): Mono<OrderAggregate> {
        val orderMono = orderHandler.handle(OrderQuery.FindOrder(id)) as Mono<OrderEntity>
        val orderAggMono = orderMono.publishOn(Schedulers.boundedElastic()).map {
            val orderItemFlux =
                orderItemHandler.handle(OrderQuery.FindOrder(it.id!!)) as Flux<OrderItemEntity>
            OrderAggregate(it.id!!, it, orderItemFlux.collectList().block().orEmpty())
        }
        return orderAggMono
    }

    override fun createOrder(dto: OrderDto): Mono<OrderAggregate> {
        val orderEntity = dto.toEntity()
        val orderItemEntities = dto.orderItems?.map { it.toEntity() }.orEmpty()
        val orderMono = orderHandler.handle(OrderCommand.CreateOrder(orderEntity))
        val orderItemFlux = orderItemHandler.handle(OrderItemCommand.CreateOrderItem(orderItemEntities))
        val aggMono = orderMono.flatMap { order ->
            orderItemFlux.collectList().map {
                OrderAggregate(order.id!!, order, it)
            }
        }
        return aggMono
    }

    override fun updateOrder(dto: OrderDto): Mono<OrderAggregate> {
        val orderEntity = dto.toEntity()
        val orderItemEntities = dto.orderItems?.map { it.toEntity() }.orEmpty()
        val orderMono = orderHandler.handle(OrderCommand.UpdateOrder(orderEntity))
        val orderItemFlux = orderItemHandler.handle(OrderItemCommand.UpdateOrderItem(orderItemEntities))
        val aggMono = orderMono.flatMap { order ->
            orderItemFlux.collectList().map {
                OrderAggregate(order.id!!, order, it)
            }
        }
        return aggMono
    }

    override fun deleteOrder(id: Long): Mono<Void> {
        val voidMono = orderHandler.handle(OrderCommand.DeleteOrder(id))
            .then<Void>(Mono.empty())
        return voidMono
    }
}
