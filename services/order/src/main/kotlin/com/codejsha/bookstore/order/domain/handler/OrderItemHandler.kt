package com.codejsha.bookstore.order.domain.handler

import com.codejsha.bookstore.order.application.port.repo.OrderItemExtendedRepo
import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.usecase.OrderItemCommand
import com.codejsha.bookstore.order.application.usecase.OrderRead
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.model.OrderItemDto

import io.opentelemetry.instrumentation.annotations.WithSpan
import org.reactivestreams.Publisher
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional
import reactor.core.publisher.Flux

interface OrderItemHandler {
    fun handle(command: OrderItemCommand): Flux<OrderItemEntity>

    fun handle(read: OrderRead): Publisher<OrderItemEntity>
}

@Component
class OrderItemHandlerImpl(
    private val orderItemRepo: OrderItemRepo,
    private val orderItemExtendedRepo: OrderItemExtendedRepo
) : OrderItemHandler {
    @Transactional
    @WithSpan
    override fun handle(command: OrderItemCommand): Flux<OrderItemEntity> =
        when (command) {
            is OrderItemCommand.CreateOrderItemsCommand -> {
                val entities: List<OrderItemEntity> = command.newEntities()
                orderItemRepo.saveAll(entities)
            }

            is OrderItemCommand.UpdateOrderItemsCommand -> {
                orderItemExtendedRepo
                    .findAllByOrderId(command.orderId)
                    .collectList()
                    .flatMapMany { existingEntities ->
                        val synced: List<OrderItemEntity> =
                            syncOrderItems(existingEntities, command.orderItems, command.orderId)
                        val toDelete: List<OrderItemEntity?> =
                            existingEntities.filterNot { old ->
                                synced.any { it.bookId == old.bookId }
                            }

                        val saveFlux = orderItemRepo.saveAll(synced)
                        val deleteMono = orderItemRepo.deleteAll(toDelete).then()
                        deleteMono.thenMany(saveFlux)
                    }
            }

            is OrderItemCommand.DeleteOrderItemsCommand -> {
                orderItemRepo
                    .deleteById(command.id)
                    .thenMany(Flux.empty())
            }
        }

    fun syncOrderItems(
        existingEntities: List<OrderItemEntity>,
        incomingDtos: List<OrderItemDto>,
        orderId: Long
    ): List<OrderItemEntity> {
        val existingMap = existingEntities.associateBy { it.bookId }.toMutableMap()
        val updatedEntities = mutableListOf<OrderItemEntity>()

        for (dto in incomingDtos) {
            val existing = existingMap.remove(dto.bookId)
            val entity =
                existing?.apply {
                    if (quantity != dto.quantity) {
                        quantity = dto.quantity
                    }
                } ?: OrderItemEntity(
                    id = null,
                    orderId = orderId,
                    bookId = dto.bookId,
                    quantity = dto.quantity
                )
            updatedEntities += entity
        }

        return updatedEntities
    }

    @Transactional(readOnly = true)
    @WithSpan
    override fun handle(read: OrderRead): Publisher<OrderItemEntity> =
        when (read) {
            is OrderRead.FindAllOrdersRead -> {
                orderItemRepo.findAll()
            }

            is OrderRead.FindOrderRead -> {
                orderItemRepo.findById(read.id)
            }
        }
}
