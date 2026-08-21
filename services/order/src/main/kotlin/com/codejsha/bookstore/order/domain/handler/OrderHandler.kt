package com.codejsha.bookstore.order.domain.handler

import com.codejsha.bookstore.order.application.port.repo.OrderExtendedRepo
import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.application.usecase.OrderRead
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity

import io.opentelemetry.instrumentation.annotations.WithSpan
import org.reactivestreams.Publisher
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional
import reactor.core.publisher.Mono

interface OrderHandler {
    fun handle(command: OrderCommand): Mono<OrderEntity>

    fun handle(read: OrderRead): Publisher<OrderEntity>
}

@Component
class OrderHandlerImpl(
    private val orderRepo: OrderRepo,
    private val orderExtendedRepo: OrderExtendedRepo
) : OrderHandler {
    @Transactional
    @WithSpan
    override fun handle(command: OrderCommand): Mono<OrderEntity> =
        when (command) {
            is OrderCommand.CreateOrderCommand -> {
                val entity = command.newEntity()
                orderRepo.save(entity)
            }

            is OrderCommand.UpdateOrderCommand -> {
                orderRepo
                    .findById(command.id)
                    .switchIfEmpty(Mono.error(IllegalArgumentException("Order with id ${command.id} not found")))
                    .map { it.update(command) }
                    .flatMap { updatedEntity -> orderRepo.save(updatedEntity) }
            }

            is OrderCommand.DeleteOrderCommand -> {
                orderRepo
                    .deleteById(command.id)
                    .then(Mono.empty())
            }

            is OrderCommand.ChangeOrderStatusCommand -> {
                orderRepo
                    .findById(command.id)
                    .switchIfEmpty(Mono.error(IllegalArgumentException("Order with id ${command.id} not found")))
                    .map { it.update(command) }
                    .flatMap { updatedEntity -> orderRepo.save(updatedEntity) }
            }
        }

    @Transactional(readOnly = true)
    @WithSpan
    override fun handle(read: OrderRead): Publisher<OrderEntity> =
        when (read) {
            is OrderRead.FindAllOrdersRead -> {
                orderExtendedRepo.findAllPaged(read.cond)
            }

            is OrderRead.FindOrderRead -> {
                orderRepo.findById(read.id)
            }
        }
}
