package com.codejsha.bookstore.order.domain.handler

import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.application.usecase.OrderQuery
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import org.reactivestreams.Publisher
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional
import reactor.core.publisher.Mono

interface OrderHandler {
    fun handle(command: OrderCommand): Mono<OrderEntity>
    fun handle(query: OrderQuery): Publisher<OrderEntity>
}

@Component
class OrderHandlerImpl(
    private val orderRepo: OrderRepo
) : OrderHandler {

    @Transactional(readOnly = true)
    override fun handle(query: OrderQuery): Publisher<OrderEntity> {
        return when (query) {
            is OrderQuery.FindAllOrders -> orderRepo.findAllPaged(query.cond)
            is OrderQuery.FindOrder -> orderRepo.findById(query.id)
        }
    }

    @Transactional
    override fun handle(command: OrderCommand): Mono<OrderEntity> {
        return when (command) {
            is OrderCommand.CreateOrder -> orderRepo.save(command.ent)
            is OrderCommand.UpdateOrder -> orderRepo.save(command.ent)
            is OrderCommand.DeleteOrder -> orderRepo.deleteById(command.id).then(Mono.empty())
        }
    }
}
