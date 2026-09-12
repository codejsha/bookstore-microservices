package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.OrderUseCase
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertManager
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertStaff
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminOrderApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminOrderResponse
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
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
    ): ResponseEntity<AdminOrderFindAllResponse> = runBlocking {
        val principal = requireStaff()
        val option = OrderQueryOption(userUid = userUid, status = status)
        val context = buildContext(principal)
        val result = orderUseCase.findAllOrders(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminOrderFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminOrderItem(it) },
            )
        )
    }

    override fun adminOrdersReadOrder(uid: String): ResponseEntity<AdminOrderResponse> = runBlocking {
        val principal = requireStaff()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminOrderResponse(orderUseCase.findOrder(uid, context)))
    }

    override fun adminOrdersCancelOrder(uid: String): ResponseEntity<AdminOrderResponse> = runBlocking {
        val principal = requireManager()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminOrderResponse(orderUseCase.cancelOrder(uid, context)))
    }

    private fun requireStaff(): Principal = principalResolver.require().also { it.assertStaff() }

    private fun requireManager(): Principal = principalResolver.require().also { it.assertManager() }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
