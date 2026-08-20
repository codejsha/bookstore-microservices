package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface OrderUseCase {
    suspend fun findAllOrders(option: OrderQueryOption, pageable: Pageable, context: ActorContext): Page<Order>

    suspend fun findOrder(uid: String, context: ActorContext): Order

    suspend fun cancelOrder(uid: String, context: ActorContext): Order
}
