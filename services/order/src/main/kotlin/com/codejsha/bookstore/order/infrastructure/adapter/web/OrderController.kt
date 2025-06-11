package com.codejsha.bookstore.order.infrastructure.adapter.web

import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.model.FilterCondition
import com.codejsha.bookstore.order.domain.model.OrderMapper
import com.codejsha.bookstore.service.application.port.openapi.api.OrderApi
import com.codejsha.bookstore.service.application.port.openapi.model.OrderCreateWebReq
import com.codejsha.bookstore.service.application.port.openapi.model.OrderFindAllWebResp
import com.codejsha.bookstore.service.application.port.openapi.model.OrderFindWebResp
import org.springframework.web.bind.annotation.RestController
import reactor.core.publisher.Mono

@RestController
class OrderController(
    private val orderUseCase: OrderUseCase
) : OrderApi {

    override fun apiV1OrdersGet(
        userId: String?,
        sort: String,
        order: String,
        limit: Int,
        offset: Int
    ): Mono<OrderFindAllWebResp> {
        val cond = FilterCondition(
            userId = userId,
            filter = FilterCondition.QueryFilter(
                sort = sort,
                order = order,
                limit = limit,
                offset = offset
            )
        )
        val orders = orderUseCase.findAllOrders(cond)
            .map { it -> it.map { it.toOrderFindWebResp() } }
        val responseMono = orders.map {
            OrderFindAllWebResp(total = it.size.toLong(), items = it)
        }
        return responseMono
    }

    override fun apiV1OrdersIdGet(id: Long): Mono<OrderFindWebResp> {
        val responseMono = orderUseCase.findOrder(id)
            .map { it.toOrderFindWebResp() }
        return responseMono
    }

    override fun apiV1OrdersPost(orderCreateWebReq: OrderCreateWebReq): Mono<Unit> {
        val dto = OrderMapper.toOrderDto(orderCreateWebReq)
        return orderUseCase.startPlaceOrderWorkflow(dto)
    }
}
