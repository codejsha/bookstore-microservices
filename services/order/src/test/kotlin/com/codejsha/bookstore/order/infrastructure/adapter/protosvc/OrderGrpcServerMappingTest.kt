package com.codejsha.bookstore.order.infrastructure.adapter.protosvc

import com.codejsha.bookstore.generated.application.port.pb.orderpb.FindOrderRequest
import com.codejsha.bookstore.generated.application.port.pb.orderpb.ListOrdersRequest
import com.codejsha.bookstore.generated.application.port.pb.orderpb.OrderServiceGrpc
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.OrderItemEntity
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.order.infrastructure.support.auth.ACTOR_ADMIN_METADATA_KEY
import com.codejsha.bookstore.order.infrastructure.support.auth.ACTOR_UID_METADATA_KEY
import com.codejsha.bookstore.order.infrastructure.support.auth.ActorAuthInterceptor
import com.codejsha.bookstore.order.support.OrderTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import com.codejsha.platform.shared.data.buildPageRequest
import io.grpc.ManagedChannel
import io.grpc.Metadata
import io.grpc.Server
import io.grpc.ServerInterceptors
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.grpc.inprocess.InProcessChannelBuilder
import io.grpc.inprocess.InProcessServerBuilder
import io.grpc.stub.MetadataUtils
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.BDDMockito.willThrow
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import com.codejsha.bookstore.generated.application.port.pb.orderpb.OrderStatus as ProtoOrderStatus

class OrderGrpcServerMappingTest {

    private val auditContext = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val ownerUid: UUID = OrderTestFixtures.USER_UID
    private val useCase = mock(OrderUseCase::class.java)

    private lateinit var server: Server
    private lateinit var channel: ManagedChannel

    @BeforeEach
    fun startServer() {
        val name = InProcessServerBuilder.generateName()
        server = InProcessServerBuilder.forName(name)
            .directExecutor()
            .addService(ServerInterceptors.intercept(OrderGrpcServer(useCase), ActorAuthInterceptor()))
            .build()
            .start()
        channel = InProcessChannelBuilder.forName(name).directExecutor().build()
    }

    @AfterEach
    fun stopServer() {
        channel.shutdownNow()
        server.shutdownNow()
    }

    @Test
    fun `findOrder_everyDomainStatus_mapsToItsProtoStatus`() {
        val expected = mapOf(
            OrderStatus.PENDING to ProtoOrderStatus.ORDER_STATUS_PENDING,
            OrderStatus.PAID to ProtoOrderStatus.ORDER_STATUS_CONFIRMED,
            OrderStatus.SHIPPED to ProtoOrderStatus.ORDER_STATUS_SHIPPED,
            OrderStatus.DELIVERED to ProtoOrderStatus.ORDER_STATUS_DELIVERED,
            OrderStatus.CANCELLED to ProtoOrderStatus.ORDER_STATUS_CANCELLED,
            OrderStatus.REFUNDED to ProtoOrderStatus.ORDER_STATUS_REFUNDED,
        )
        assertEquals(OrderStatus.entries.toSet(), expected.keys)

        expected.forEach { (domain, proto) ->
            given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
                .willReturn(orderAggregate(status = domain))

            val response = stub().findOrder(findRequest(OrderTestFixtures.ORDER_UID.toString()))

            assertEquals(proto, response.order.status, "domain status $domain")
        }
    }

    @Test
    fun `findOrder_orderWithItems_mapsAmountsAndLineItems`() {
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext)).willReturn(
            orderAggregate(
                items = listOf(
                    item(productId = 200L, quantity = 2, price = BigDecimal("5000.00")),
                    item(productId = 300L, quantity = 1, price = BigDecimal("1250.50")),
                ),
            ),
        )

        val order = stub().findOrder(findRequest(OrderTestFixtures.ORDER_UID.toString())).order

        assertEquals("ORD-20260101-0001", order.orderNumber)
        assertEquals("KRW", order.currency)
        assertEquals(11250.5, order.itemsAmount)
        assertEquals(1000.0, order.discountAmount)
        assertEquals(2500.0, order.shippingAmount)
        assertEquals(125.25, order.taxAmount)
        assertEquals(12875.75, order.totalAmount)
        assertEquals(listOf("200", "300"), order.lineItemsList.map { it.productUid })
        assertEquals(listOf(2, 1), order.lineItemsList.map { it.quantity })
        assertEquals(listOf(5000.0, 1250.5), order.lineItemsList.map { it.price })
        assertEquals(listOf("KRW", "KRW"), order.lineItemsList.map { it.currency })
    }

    @Test
    fun `findOrder_uidNotUuid_failsInvalidArgumentWithoutCallingUseCase`() {
        val ex = assertFailsWith<StatusRuntimeException> { stub().findOrder(findRequest("not-a-uuid")) }

        assertEquals(Status.Code.INVALID_ARGUMENT, ex.status.code)
        verifyNoInteractions(useCase)
    }

    @Test
    fun `findOrder_useCaseRejectsArgument_failsInvalidArgumentWithMessage`() {
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
            .willThrow(IllegalArgumentException("bad order reference"))

        val ex = assertFailsWith<StatusRuntimeException> {
            stub().findOrder(findRequest(OrderTestFixtures.ORDER_UID.toString()))
        }

        assertEquals(Status.Code.INVALID_ARGUMENT, ex.status.code)
        assertEquals("bad order reference", ex.status.description)
    }

    @Test
    fun `findOrder_useCaseFailsUnexpectedly_failsInternalWithoutLeakingCause`() {
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
            .willThrow(IllegalStateException("jdbc:postgresql://db/secret refused"))

        val ex = assertFailsWith<StatusRuntimeException> {
            stub().findOrder(findRequest(OrderTestFixtures.ORDER_UID.toString()))
        }

        assertEquals(Status.Code.INTERNAL, ex.status.code)
        assertEquals("internal error", ex.status.description)
    }

    @Test
    fun `listOrders_everyProtoStatusFilter_passesItsDomainStatusValue`() {
        val expected = mapOf(
            ProtoOrderStatus.ORDER_STATUS_PENDING to "PENDING",
            ProtoOrderStatus.ORDER_STATUS_CONFIRMED to "PAID",
            ProtoOrderStatus.ORDER_STATUS_PROCESSING to "PAID",
            ProtoOrderStatus.ORDER_STATUS_SHIPPED to "SHIPPED",
            ProtoOrderStatus.ORDER_STATUS_DELIVERED to "DELIVERED",
            ProtoOrderStatus.ORDER_STATUS_CANCELLED to "CANCELLED",
            ProtoOrderStatus.ORDER_STATUS_REFUNDED to "REFUNDED",
        )
        val pageable = buildPageRequest(PAGE_SIZE, 0, "")

        expected.forEach { (proto, domainValue) ->
            given(useCase.findAllOrders(OrderQueryOption(userUid = ownerUid, status = domainValue), pageable, auditContext))
                .willReturn(PageImpl(listOf(orderAggregate(status = OrderStatus.fromValue(domainValue))), pageable, 1L))

            val response = stub().listOrders(
                ListOrdersRequest.newBuilder().setStatus(proto).setPageSize(PAGE_SIZE).build(),
            )

            assertEquals(1, response.totalSize, "proto status $proto")
        }
    }

    @Test
    fun `listOrders_morePagesRemain_setsNextPageToken`() {
        val pageable = buildPageRequest(PAGE_SIZE, 1, "")
        given(useCase.findAllOrders(OrderQueryOption(userUid = ownerUid), pageable, auditContext))
            .willReturn(PageImpl(listOf(orderAggregate(), orderAggregate()), pageable, 45L))

        val response = stub().listOrders(
            ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).setPageToken("1").build(),
        )

        assertEquals(45, response.totalSize)
        assertEquals(2, response.ordersCount)
        assertEquals("2", response.nextPageToken)
    }

    @Test
    fun `listOrders_lastPage_omitsNextPageToken`() {
        val pageable = buildPageRequest(PAGE_SIZE, 2, "")
        given(useCase.findAllOrders(OrderQueryOption(userUid = ownerUid), pageable, auditContext))
            .willReturn(PageImpl(listOf(orderAggregate()), pageable, 41L))

        val response = stub().listOrders(
            ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).setPageToken("2").build(),
        )

        assertEquals("", response.nextPageToken)
    }

    @Test
    fun `listOrders_pageTokenNotNumeric_startsFromFirstPage`() {
        val pageable = buildPageRequest(PAGE_SIZE, 0, "")
        given(useCase.findAllOrders(OrderQueryOption(userUid = ownerUid), pageable, auditContext))
            .willReturn(PageImpl(listOf(orderAggregate()), pageable, 30L))

        val response = stub().listOrders(
            ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).setPageToken("abc").build(),
        )

        assertEquals(1, response.ordersCount)
        assertEquals("1", response.nextPageToken)
    }

    @Test
    fun `listOrders_adminUserUidNotUuid_failsInvalidArgumentWithoutCallingUseCase`() {
        val request = ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).setUserUid("not-a-uuid").build()

        val ex = assertFailsWith<StatusRuntimeException> { stub(admin = true).listOrders(request) }

        assertEquals(Status.Code.INVALID_ARGUMENT, ex.status.code)
        verifyNoInteractions(useCase)
    }

    @Test
    fun `listOrders_statusNotInEnum_failsInvalidArgumentWithoutCallingUseCase`() {
        val request = ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).setStatusValue(999).build()

        val ex = assertFailsWith<StatusRuntimeException> { stub().listOrders(request) }

        assertEquals(Status.Code.INVALID_ARGUMENT, ex.status.code)
        verifyNoInteractions(useCase)
    }

    @Test
    fun `listOrders_useCaseFailures_mapToGrpcStatusCodes`() {
        val pageable = buildPageRequest(PAGE_SIZE, 0, "")
        val request = ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).build()
        val cases = listOf(
            NoSuchElementException("gone") to Status.Code.NOT_FOUND,
            IllegalArgumentException("bad sort") to Status.Code.INVALID_ARGUMENT,
            IllegalStateException("pool exhausted") to Status.Code.INTERNAL,
        )

        cases.forEach { (failure, code) ->
            willThrow(failure).given(useCase)
                .findAllOrders(OrderQueryOption(userUid = ownerUid), pageable, auditContext)

            val ex = assertFailsWith<StatusRuntimeException> { stub().listOrders(request) }

            assertEquals(code, ex.status.code, "failure ${failure::class.simpleName}")
            assertFalse(ex.status.description.orEmpty().contains("pool exhausted"))
        }
    }

    private fun stub(admin: Boolean = false): OrderServiceGrpc.OrderServiceBlockingStub {
        val md = Metadata()
        md.put(ACTOR_UID_METADATA_KEY, ownerUid.toString())
        if (admin) md.put(ACTOR_ADMIN_METADATA_KEY, "true")
        return OrderServiceGrpc.newBlockingStub(channel)
            .withInterceptors(MetadataUtils.newAttachHeadersInterceptor(md))
    }

    private fun findRequest(uid: String): FindOrderRequest = FindOrderRequest.newBuilder().setUid(uid).build()

    private fun item(productId: Long, quantity: Int, price: BigDecimal) = OrderItemEntity(
        id = productId,
        uid = UUID.randomUUID(),
        orderId = 1L,
        productId = productId,
        sku = null,
        productName = null,
        options = null,
        quantity = quantity,
        currency = "KRW",
        price = price,
        taxRate = BigDecimal.ZERO,
        subtotal = price * BigDecimal(quantity),
        createdAt = OrderTestFixtures.FIXED_TIME,
        updatedAt = null,
    )

    private fun orderAggregate(
        status: OrderStatus = OrderStatus.PENDING,
        items: List<OrderItemEntity> = emptyList(),
    ): OrderAggregate = OrderAggregate(
        id = 1L,
        uid = OrderTestFixtures.ORDER_UID,
        userUid = ownerUid,
        orderNumber = "ORD-20260101-0001",
        status = status,
        currency = "KRW",
        itemsAmount = BigDecimal("11250.50"),
        discountAmount = BigDecimal("1000.00"),
        shippingAmount = BigDecimal("2500.00"),
        taxAmount = BigDecimal("125.25"),
        totalAmount = BigDecimal("12875.75"),
        idempotencyKey = "idem_001",
        paymentUid = null,
        items = items,
        adjustments = emptyList(),
        shipping = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    companion object {
        private const val PAGE_SIZE = 20
    }
}
