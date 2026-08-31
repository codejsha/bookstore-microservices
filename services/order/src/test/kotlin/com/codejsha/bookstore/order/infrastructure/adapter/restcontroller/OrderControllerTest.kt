package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.OrderAdjustmentCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.OrderAdjustmentType
import com.codejsha.bookstore.generated.application.port.openapi.model.OrderItemCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.OrderItemUpdateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.OrderShippingCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.OrderShippingUpdateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.PlaceOrderRequest
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAdjustmentEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.OrderItemEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderShippingEntity
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.command.OrderAdjustmentCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemUpdateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingUpdateCommand
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.order.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.order.infrastructure.support.auth.ROLE_ADMIN
import com.codejsha.bookstore.order.support.OrderTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.math.BigDecimal
import java.time.LocalDateTime
import kotlin.test.assertEquals
import kotlin.test.assertNotNull

class OrderControllerTest {

    private val subject = OrderTestFixtures.USER_UID.toString()
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val resolver = HttpPrincipalResolver(ObjectMapper())

    @BeforeEach
    fun bindPrincipal() {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", subject)
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    private fun bindAdmin() {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", subject)
        request.addHeader("X-User-Roles", ROLE_ADMIN)
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @Test
    fun `ordersGetAll_whenFiltersGiven_forwardsThemAndMapsPage`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        val unpaged = Pageable.unpaged()
        val expectedOption = OrderQueryOption(userUid = OrderTestFixtures.USER_UID, status = "PENDING")
        given(useCase.findAllOrders(expectedOption, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(orderAggregate()), unpaged, 1L))

        val response = controller.ordersGetAll(
            userUid = OrderTestFixtures.USER_UID.toString(),
            status = com.codejsha.bookstore.generated.application.port.openapi.model.OrderStatus.PENDING,
            pageable = null,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals("ORD-20260101-0001", body.items[0].orderNumber)
    }

    @Test
    fun `ordersPlace_whenRequestValid_buildsOrderAndItemCommands`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        bindAdmin()

        val request = PlaceOrderRequest(
            userUid = OrderTestFixtures.USER_UID.toString(),
            currency = "KRW",
            itemsAmount = 10_000.0,
            discountAmount = null,
            shippingAmount = null,
            taxAmount = null,
            totalAmount = 10_000.0,
            idempotencyKey = "idem_001",
            items = listOf(
                OrderItemCreateRequest(
                    productId = 200L,
                    productName = "Book A",
                    sku = "SKU-A",
                    options = null,
                    quantity = 2,
                    currency = "KRW",
                    price = 5_000.0,
                    taxRate = null,
                ),
            ),
            shipping = OrderShippingCreateRequest(
                recipientName = "Alice",
                recipientPhone = "010-1",
                addressLine1 = "1 Foo St",
                addressLine2 = null,
                city = "Seoul",
                state = "Seoul",
                postalCode = "12345",
                country = "KR",
                shippingMethod = "STANDARD",
            ),
        )
        val expectedOrderCmd = OrderCreateCommand(
            userUid = OrderTestFixtures.USER_UID,
            currency = "KRW",
            itemsAmount = BigDecimal("10000.0"),
            discountAmount = BigDecimal.ZERO,
            shippingAmount = BigDecimal.ZERO,
            taxAmount = BigDecimal.ZERO,
            totalAmount = BigDecimal("10000.0"),
            idempotencyKey = "idem_001",
        )
        val expectedItems = listOf(
            OrderItemCreateCommand(
                productId = 200L,
                sku = "SKU-A",
                productName = "Book A",
                options = null,
                quantity = 2,
                currency = "KRW",
                price = BigDecimal("5000.0"),
                taxRate = BigDecimal.ZERO,
            ),
        )
        val expectedShipping = OrderShippingCreateCommand(
            recipientName = "Alice",
            recipientPhone = "010-1",
            addressLine1 = "1 Foo St",
            addressLine2 = null,
            city = "Seoul",
            state = "Seoul",
            postalCode = "12345",
            country = "KR",
            shippingMethod = "STANDARD",
        )
        given(useCase.placeOrder(expectedOrderCmd, expectedItems, expectedShipping, controllerContext))
            .willReturn(orderAggregate())

        val response = controller.ordersPlace(request)

        assertEquals(HttpStatus.CREATED, response.statusCode)
        assertEquals(
            "/api/v1/orders/${OrderTestFixtures.ORDER_UID}",
            response.headers.location?.toString(),
        )
        verify(useCase).placeOrder(expectedOrderCmd, expectedItems, expectedShipping, controllerContext)
    }

    @Test
    fun `ordersRead_whenOrderExists_mapsItemsShippingAndAdjustments`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext))
            .willReturn(orderAggregate(withItems = true, withShipping = true, withAdjustments = true))

        val response = controller.ordersRead(OrderTestFixtures.ORDER_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1, assertNotNull(body.items).size)
        assertNotNull(body.shipping)
        assertEquals(1, assertNotNull(body.adjustments).size)
    }

    @Test
    fun `ordersCancel_whenOrderCancellable_returnsMappedAggregate`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext))
            .willReturn(orderAggregate())
        given(useCase.cancelOrder(OrderTestFixtures.ORDER_UID, controllerContext))
            .willReturn(orderAggregate(status = OrderStatus.CANCELLED))

        val response = controller.ordersCancel(OrderTestFixtures.ORDER_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.OrderStatus.CANCELLED,
            assertNotNull(response.body).status,
        )
    }

    @Test
    fun `orderItemsAdd_whenRequestValid_buildsItemCommandAndReturnsCreated`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        val request = OrderItemCreateRequest(
            productId = 200L,
            productName = "Book A",
            sku = null,
            options = null,
            quantity = 1,
            currency = "KRW",
            price = 5_000.0,
            taxRate = null,
        )
        val expectedCmd = OrderItemCreateCommand(
            productId = 200L,
            sku = null,
            productName = "Book A",
            options = null,
            quantity = 1,
            currency = "KRW",
            price = BigDecimal("5000.0"),
            taxRate = BigDecimal.ZERO,
        )
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())
        given(useCase.addItem(OrderTestFixtures.ORDER_UID, expectedCmd, controllerContext))
            .willReturn(orderItemEntity())

        val response = controller.orderItemsAdd(OrderTestFixtures.ORDER_UID.toString(), request)

        assertEquals(HttpStatus.CREATED, response.statusCode)
        assertEquals(
            "/api/v1/orders/${OrderTestFixtures.ORDER_UID}/items/${OrderTestFixtures.ORDER_ITEM_UID}",
            response.headers.location?.toString(),
        )
    }

    @Test
    fun `orderItemsUpdate_whenRequestValid_buildsUpdateCommandAndReturnsOk`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        val request = OrderItemUpdateRequest(
            productId = null,
            sku = null,
            productName = null,
            options = null,
            quantity = 5,
            currency = null,
            price = null,
            taxRate = null,
        )
        val expectedCmd = OrderItemUpdateCommand(
            productId = null, sku = null, productName = null, options = null,
            quantity = 5, currency = null, price = null, taxRate = null,
        )
        given(
            useCase.updateItem(
                OrderTestFixtures.ORDER_UID,
                OrderTestFixtures.ORDER_ITEM_UID,
                expectedCmd,
                controllerContext,
            )
        ).willReturn(orderItemEntity(quantity = 5))
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())

        val response = controller.orderItemsUpdate(
            OrderTestFixtures.ORDER_UID.toString(),
            OrderTestFixtures.ORDER_ITEM_UID.toString(),
            request,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals(5, assertNotNull(response.body).quantity)
    }

    @Test
    fun `orderItemsRemove_whenRequestValid_forwardsUidsAndReturnsNoContent`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())

        val response = controller.orderItemsRemove(
            OrderTestFixtures.ORDER_UID.toString(),
            OrderTestFixtures.ORDER_ITEM_UID.toString(),
        )

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).removeItem(
            OrderTestFixtures.ORDER_UID,
            OrderTestFixtures.ORDER_ITEM_UID,
            controllerContext,
        )
    }

    @Test
    fun `orderShippingSet_whenRequestValid_buildsCommandAndReturnsOk`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        val request = OrderShippingCreateRequest(
            recipientName = "Alice",
            recipientPhone = "010-1",
            addressLine1 = "1 Foo St",
            addressLine2 = null,
            city = "Seoul",
            state = "Seoul",
            postalCode = "12345",
            country = "KR",
            shippingMethod = "STANDARD",
        )
        val expectedCmd = OrderShippingCreateCommand(
            recipientName = "Alice",
            recipientPhone = "010-1",
            addressLine1 = "1 Foo St",
            addressLine2 = null,
            city = "Seoul",
            state = "Seoul",
            postalCode = "12345",
            country = "KR",
            shippingMethod = "STANDARD",
        )
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())
        given(useCase.setShipping(OrderTestFixtures.ORDER_UID, expectedCmd, controllerContext))
            .willReturn(orderShippingEntity())

        val response = controller.orderShippingSet(OrderTestFixtures.ORDER_UID.toString(), request)

        assertEquals(HttpStatus.OK, response.statusCode)
    }

    @Test
    fun `orderShippingUpdate_whenRequestValid_buildsUpdateCommandAndReturnsOk`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        val request = OrderShippingUpdateRequest(
            recipientName = "Bob",
            recipientPhone = null,
            addressLine1 = null,
            addressLine2 = null,
            city = null,
            state = null,
            postalCode = null,
            country = null,
            shippingMethod = null,
        )
        val expectedCmd = OrderShippingUpdateCommand(
            recipientName = "Bob",
            recipientPhone = null,
            addressLine1 = null,
            addressLine2 = null,
            city = null,
            state = null,
            postalCode = null,
            country = null,
            shippingMethod = null,
        )
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())
        given(useCase.updateShipping(OrderTestFixtures.ORDER_UID, expectedCmd, controllerContext))
            .willReturn(orderShippingEntity())

        val response = controller.orderShippingUpdate(OrderTestFixtures.ORDER_UID.toString(), request)

        assertEquals(HttpStatus.OK, response.statusCode)
    }

    @Test
    fun `orderAdjustmentsApply_whenRequestValid_mapsBodyToCommandAndReturnsCreated`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        val request = OrderAdjustmentCreateRequest(
            type = OrderAdjustmentType.COUPON,
            label = "WELCOME",
            amount = -1000.0,
            meta = null,
        )
        val expectedCmd = OrderAdjustmentCreateCommand(
            type = "COUPON",
            label = "WELCOME",
            amount = BigDecimal("-1000.0"),
            meta = null,
        )
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())
        given(useCase.applyAdjustment(OrderTestFixtures.ORDER_UID, expectedCmd, controllerContext))
            .willReturn(orderAdjustmentEntity())

        val response = controller.orderAdjustmentsApply(OrderTestFixtures.ORDER_UID.toString(), request)

        assertEquals(HttpStatus.CREATED, response.statusCode)
        assertEquals(
            "/api/v1/orders/${OrderTestFixtures.ORDER_UID}/adjustments/${OrderTestFixtures.ORDER_ADJUSTMENT_UID}",
            response.headers.location?.toString(),
        )
    }

    @Test
    fun `orderAdjustmentsRemove_whenRequestValid_returnsNoContent`(): Unit = runBlocking {
        val useCase = mock(OrderUseCase::class.java)
        val controller = OrderController(useCase, resolver)

        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, controllerContext)).willReturn(orderAggregate())

        val response = controller.orderAdjustmentsRemove(
            OrderTestFixtures.ORDER_UID.toString(),
            OrderTestFixtures.ORDER_ADJUSTMENT_UID.toString(),
        )

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).removeAdjustment(
            OrderTestFixtures.ORDER_UID,
            OrderTestFixtures.ORDER_ADJUSTMENT_UID,
            controllerContext,
        )
    }

    private fun orderAggregate(
        status: OrderStatus = OrderStatus.PENDING,
        withItems: Boolean = false,
        withShipping: Boolean = false,
        withAdjustments: Boolean = false,
    ): OrderAggregate = OrderAggregate(
        id = 1L,
        uid = OrderTestFixtures.ORDER_UID,
        userUid = OrderTestFixtures.USER_UID,
        orderNumber = "ORD-20260101-0001",
        status = status,
        currency = "KRW",
        itemsAmount = BigDecimal("10000.00"),
        discountAmount = BigDecimal.ZERO,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = BigDecimal("10000.00"),
        idempotencyKey = "idem_001",
        paymentUid = null,
        items = if (withItems) listOf(orderItemEntity()) else emptyList(),
        adjustments = if (withAdjustments) listOf(orderAdjustmentEntity()) else emptyList(),
        shipping = if (withShipping) orderShippingEntity() else null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    private fun orderItemEntity(quantity: Int = 1): OrderItemEntity = OrderItemEntity(
        id = 1L,
        uid = OrderTestFixtures.ORDER_ITEM_UID,
        orderId = 1L,
        productId = 200L,
        sku = "SKU-A",
        productName = "Book A",
        options = null,
        quantity = quantity,
        currency = "KRW",
        price = BigDecimal("5000.00"),
        taxRate = BigDecimal.ZERO,
        subtotal = BigDecimal("5000.00") * BigDecimal(quantity),
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    private fun orderShippingEntity(): OrderShippingEntity = OrderShippingEntity(
        id = 1L,
        uid = OrderTestFixtures.ORDER_SHIPPING_UID,
        orderId = 1L,
        recipientName = "Alice",
        recipientPhone = "010-1",
        addressLine1 = "1 Foo St",
        addressLine2 = null,
        city = "Seoul",
        state = "Seoul",
        postalCode = "12345",
        country = "KR",
        shippingMethod = "STANDARD",
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    private fun orderAdjustmentEntity(): OrderAdjustmentEntity = OrderAdjustmentEntity(
        id = 1L,
        uid = OrderTestFixtures.ORDER_ADJUSTMENT_UID,
        orderId = 1L,
        type = com.codejsha.bookstore.order.domain.constant.OrderAdjustmentType.COUPON,
        label = "WELCOME",
        amount = BigDecimal("-1000.00"),
        meta = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )
}
