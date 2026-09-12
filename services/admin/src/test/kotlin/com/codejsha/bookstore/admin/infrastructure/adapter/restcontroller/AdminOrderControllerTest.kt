package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.OrderUseCase
import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.external.OrderLine
import com.codejsha.bookstore.admin.domain.model.external.OrderShipping
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull

class AdminOrderControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `every order endpoint rejects a caller without the STAFF role`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminOrdersListOrders(null, null, null) }
        assertFailsWith<ForbiddenException> { controller.adminOrdersReadOrder(ORDER_UID) }
        assertFailsWith<ForbiddenException> { controller.adminOrdersCancelOrder(ORDER_UID) }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `an unfiltered list forwards no owner and no paging`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        given(useCase.findAllOrders(OrderQueryOption(), unpaged, controllerContext))
            .willReturn(PageImpl(listOf(order()), unpaged, 1))

        val body = controller.adminOrdersListOrders(null, null, null).body!!

        assertEquals(1L, body.total)
        assertEquals(ORDER_UID, body.items.first().uid)
        assertEquals(USER_UID, body.items.first().userUid)
        assertEquals("PAID", body.items.first().status)
        assertEquals(42.0, body.items.first().totalAmount)
    }

    @Test
    fun `a filtered page forwards the owner, the status and the requested page`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)
        val option = OrderQueryOption(userUid = USER_UID, status = "PAID")
        val pageable = PageRequest.of(2, 20, Sort.by("created_at").descending())
        given(useCase.findAllOrders(option, pageable, controllerContext))
            .willReturn(PageImpl(emptyList(), pageable, 0))

        val body = controller.adminOrdersListOrders(USER_UID, "PAID", pageable).body!!

        assertEquals(0L, body.total)
    }

    @Test
    fun `readOrder maps the lines and the shipping destination`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)
        given(useCase.findOrder(ORDER_UID, controllerContext)).willReturn(order(withLines = true))

        val body = controller.adminOrdersReadOrder(ORDER_UID).body!!

        assertEquals(1, body.items?.size)
        val line = body.items!!.first()
        assertEquals("Structure and Interpretation", line.productName)
        assertEquals(2, line.quantity)
        assertEquals("Ada Lovelace", body.shipping?.recipientName)
        assertEquals("Seoul", body.shipping?.city)
    }

    @Test
    fun `an order without lines maps items to null`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)
        given(useCase.findOrder(ORDER_UID, controllerContext)).willReturn(order(withLines = false))

        val body = controller.adminOrdersReadOrder(ORDER_UID).body!!

        assertNull(body.items)
        assertNull(body.shipping)
    }

    @Test
    fun `cancelOrder answers with the order as the order service left it`(): Unit = runBlocking {
        bindPrincipal(roles = "MANAGE,STAFF,USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)
        given(useCase.cancelOrder(ORDER_UID, controllerContext)).willReturn(order(status = "CANCELLED"))

        val response = controller.adminOrdersCancelOrder(ORDER_UID)

        assertEquals(200, response.statusCode.value())
        assertEquals("CANCELLED", response.body!!.status)
    }

    @Test
    fun `a staff caller cannot cancel an order`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(OrderUseCase::class.java)
        val controller = AdminOrderController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminOrdersCancelOrder(ORDER_UID) }
        verifyNoInteractions(useCase)
    }

    private companion object {
        private const val ORDER_UID = "33333333-3333-3333-3333-333333333333"
        private const val USER_UID = "22222222-2222-2222-2222-222222222222"
        private val CREATED_AT: OffsetDateTime = OffsetDateTime.of(2026, 3, 1, 9, 0, 0, 0, ZoneOffset.UTC)

        private fun order(
            status: String = "PAID",
            withLines: Boolean = false,
        ) = Order(
            uid = ORDER_UID,
            orderNumber = "ORD-0001",
            userUid = USER_UID,
            status = status,
            currency = "USD",
            itemsAmount = 40.0,
            discountAmount = 0.0,
            shippingAmount = 2.0,
            taxAmount = 0.0,
            totalAmount = 42.0,
            lines = if (withLines) {
                listOf(
                    OrderLine(
                        uid = "44444444-4444-4444-4444-444444444444",
                        productId = 7L,
                        productName = "Structure and Interpretation",
                        sku = "SICP-1",
                        quantity = 2,
                        price = 20.0,
                        subtotal = 40.0,
                        taxRate = 0.0,
                        currency = "USD",
                        options = null,
                    )
                )
            } else {
                emptyList()
            },
            shipping = if (withLines) {
                OrderShipping(
                    recipientName = "Ada Lovelace",
                    recipientPhone = "+82-10-0000-0000",
                    addressLine1 = "1 Main Street",
                    addressLine2 = null,
                    city = "Seoul",
                    state = "Seoul",
                    postalCode = "04524",
                    country = "KR",
                    shippingMethod = "STANDARD",
                )
            } else {
                null
            },
            createdAt = CREATED_AT,
            updatedAt = null,
        )
    }
}
