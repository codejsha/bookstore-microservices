package com.codejsha.bookstore.order.domain.handler

import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.usecase.OrderItemCommand
import com.codejsha.bookstore.order.application.usecase.OrderQuery
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import org.reactivestreams.Publisher
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional
import reactor.core.publisher.Flux

interface OrderItemHandler {
    fun handle(command: OrderItemCommand): Flux<OrderItemEntity>
    fun handle(query: OrderQuery): Publisher<OrderItemEntity>
}

@Component
class OrderItemHandlerImpl(
    private val orderItemRepo: OrderItemRepo
) : OrderItemHandler {

    @Transactional
    override fun handle(command: OrderItemCommand): Flux<OrderItemEntity> {
        return when (command) {
            is OrderItemCommand.CreateOrderItem -> orderItemRepo.saveAll(command.ents)
            is OrderItemCommand.UpdateOrderItem -> orderItemRepo.saveAll(command.ents)
            is OrderItemCommand.DeleteOrderItem -> {
                orderItemRepo.deleteById(command.id).thenMany(Flux.empty())
            }
        }
    }


    @Transactional(readOnly = true)
    override fun handle(query: OrderQuery): Publisher<OrderItemEntity> {
        return when (query) {
            is OrderQuery.FindAllOrders -> orderItemRepo.findAll()
            is OrderQuery.FindOrder -> orderItemRepo.findById(query.id)
        }
    }
}
