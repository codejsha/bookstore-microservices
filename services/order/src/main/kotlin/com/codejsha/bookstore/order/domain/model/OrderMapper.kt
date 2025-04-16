package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.service.application.port.openapi.model.OrderCreateWebReq
import com.codejsha.bookstore.service.application.port.openapi.model.OrderFindAllWebParam
import com.codejsha.bookstore.service.application.port.openapi.model.OrderUpdateWebReq
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq

sealed class OrderMapper {
    companion object {
        fun toOrderEntity(req: OrderCreateWebReq): OrderEntity = req.toOrderEntity()
        fun toOrderEntity(id: Long, req: OrderUpdateWebReq): OrderEntity = req.toOrderEntity(id)
        fun toFilterCondition(param: OrderFindAllWebParam): FilterCondition = param.toFilterCondition()
        fun toFilterCondition(req: OrderFindAllProtoReq): FilterCondition = req.toFilterCondition()
    }
}

private fun OrderCreateWebReq.toOrderEntity(): OrderEntity {
    return OrderEntity(
        id = null,
        userId = this.userId!!,
        totalPrice = this.totalPrice!!,
        status = this.status!!,
        createdAt = null,
        updatedAt = null
    )
}

private fun OrderUpdateWebReq.toOrderEntity(id: Long): OrderEntity {
    return OrderEntity(
        id = id,
        userId = this.userId!!,
        totalPrice = this.totalPrice!!,
        status = this.status!!,
        createdAt = null,
        updatedAt = null
    )
}

private fun OrderFindAllWebParam.toFilterCondition(): FilterCondition {
    return FilterCondition(
        userId = this.userId,
        filter = FilterCondition.QueryFilter(
            sort = this.sort,
            order = this.order,
            limit = this.limit,
            offset = this.offset
        )
    )
}

private fun OrderFindAllProtoReq.toFilterCondition(): FilterCondition {
    return FilterCondition(
        userId = this.userId,
        filter = FilterCondition.QueryFilter(
            sort = this.filter.sort,
            order = this.filter.order,
            limit = this.filter.limit,
            offset = this.filter.offset
        )
    )
}
