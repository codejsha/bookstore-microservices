package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.OrderUseCase
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminOrderApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderResponse
import org.springframework.data.domain.Pageable
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminOrderController(
    private val orderUseCase: OrderUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminOrderApi {

    override fun adminOrdersListOrders(
        userUid: String?,
        status: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminOrderFindAllResponse> {
        principalResolver.requireStaff()
        val option = OrderQueryOption(userUid = userUid, status = status)
        val context = buildContext()
        val result = orderUseCase.findAllOrders(option, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminOrderFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminOrderItem(it) },
            )
        )
    }

    override fun adminOrdersReadOrder(uid: String): ResponseEntity<AdminOrderResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        return ResponseEntity.ok(toAdminOrderResponse(orderUseCase.findOrder(uid, context)))
    }

    override fun adminOrdersCancelOrder(uid: String): ResponseEntity<AdminOrderResponse> {
        principalResolver.requireManager()
        val context = buildContext()
        return ResponseEntity.ok(toAdminOrderResponse(orderUseCase.cancelOrder(uid, context)))
    }
}
