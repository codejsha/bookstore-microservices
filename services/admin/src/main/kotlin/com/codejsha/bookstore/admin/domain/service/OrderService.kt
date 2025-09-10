package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.application.usecase.OrderUseCase
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class OrderService(
    private val orderClient: OrderClient,
) : OrderUseCase {

    override suspend fun findAllOrders(option: OrderQueryOption, pageable: Pageable, context: ActorContext) =
        orderClient.findAllOrders(option, pageable)

    override suspend fun findOrder(uid: String, context: ActorContext) = orderClient.findOrder(uid)

    override suspend fun cancelOrder(uid: String, context: ActorContext) = orderClient.cancelOrder(uid)
}
