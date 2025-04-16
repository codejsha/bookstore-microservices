package com.codejsha.bookstore.payment.domain.model

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentCreateWebReq
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindAllWebParam
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentUpdateWebReq
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq

sealed class PaymentMapper {
    companion object {
        fun toPaymentDto(req: PaymentCreateWebReq): PaymentDto = req.toPaymentDto()
        fun toPaymentDto(id: Long, req: PaymentUpdateWebReq): PaymentDto = req.toPaymentDto(id)

        fun toFilterCondition(param: PaymentFindAllWebParam): FilterCondition = param.toFilterCondition()
        fun toFilterCondition(request: PaymentFindAllProtoReq): FilterCondition = request.toFilterCondition()
    }
}

private fun PaymentCreateWebReq.toPaymentDto(): PaymentDto {
    return PaymentDto(
        id = null,
        orderId = this.orderId,
        userId = this.userId,
        paymentType = this.paymentType,
        cardNumber = this.cardNumber,
        amount = this.amount,
        paymentDate = this.paymentDate?.toLocalDateTime()
    )
}

private fun PaymentUpdateWebReq.toPaymentDto(id: Long): PaymentDto {
    return PaymentDto(
        id = id,
        orderId = this.orderId,
        userId = this.userId,
        paymentType = this.paymentType,
        cardNumber = this.cardNumber,
        amount = this.amount,
        paymentDate = this.paymentDate?.toLocalDateTime()
    )
}

private fun PaymentFindAllWebParam.toFilterCondition(): FilterCondition {
    return FilterCondition(
        filter = FilterCondition.QueryFilter(
            sort = this.sort,
            order = this.order,
            limit = this.limit,
            offset = this.offset
        )
    )
}

private fun PaymentFindAllProtoReq.toFilterCondition(): FilterCondition {
    return FilterCondition(
        filter = FilterCondition.QueryFilter(
            sort = this.filter.sort,
            order = this.filter.order,
            limit = this.filter.limit,
            offset = this.filter.offset
        )
    )
}
