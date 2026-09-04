package com.codejsha.bookstore.order.infrastructure.adapter.protosvc

import com.codejsha.bookstore.generated.application.port.pb.orderpb.FindOrderRequest
import com.codejsha.bookstore.generated.application.port.pb.orderpb.ListOrdersRequest
import com.codejsha.bookstore.generated.application.port.pb.orderpb.OrderServiceGrpc
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
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
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.springframework.data.domain.PageImpl
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class OrderGrpcServerTest {

    private val auditContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private val ownerUid: UUID = OrderTestFixtures.USER_UID
    private val otherUid: UUID = UUID.fromString("99999999-9999-9999-9999-999999999999")

    private val expectedPageable = buildPageRequest(PAGE_SIZE, 0, "")

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

    private fun stub(actorUid: String?, admin: Boolean = false): OrderServiceGrpc.OrderServiceBlockingStub {
        val md = Metadata()
        if (actorUid != null) md.put(ACTOR_UID_METADATA_KEY, actorUid)
        if (admin) md.put(ACTOR_ADMIN_METADATA_KEY, "true")
        return OrderServiceGrpc.newBlockingStub(channel)
            .withInterceptors(MetadataUtils.newAttachHeadersInterceptor(md))
    }

    // ─── findOrder ─────────────────────────────────────────────────────────────

    @Test
    fun `findOrder_whenOrderOwnedByActor_returnsOrder`(): Unit = runBlocking {
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
            .willReturn(orderAggregate(userUid = ownerUid))

        val response = stub(ownerUid.toString()).findOrder(findRequest())

        assertEquals(OrderTestFixtures.ORDER_UID.toString(), response.order.uid)
        assertEquals(ownerUid.toString(), response.order.userUid)
    }

    @Test
    fun `findOrder_whenActorIsAdmin_returnsAnotherUsersOrder`(): Unit = runBlocking {
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
            .willReturn(orderAggregate(userUid = ownerUid))

        val response = stub(otherUid.toString(), admin = true).findOrder(findRequest())

        assertEquals(ownerUid.toString(), response.order.userUid)
    }

    @Test
    fun `findOrder_whenOrderOwnedByAnotherUser_returnsNotFound`(): Unit = runBlocking {
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
            .willReturn(orderAggregate(userUid = ownerUid))
        val foreign = assertFailsWith<StatusRuntimeException> {
            stub(otherUid.toString()).findOrder(findRequest())
        }
        assertEquals(Status.Code.NOT_FOUND, foreign.status.code)

        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, auditContext))
            .willThrow(NoSuchElementException("Order with uid ${OrderTestFixtures.ORDER_UID} not found"))
        val missing = assertFailsWith<StatusRuntimeException> {
            stub(ownerUid.toString()).findOrder(findRequest())
        }
        assertEquals(Status.Code.NOT_FOUND, missing.status.code)
        assertEquals(foreign.status.description, missing.status.description)
    }

    @Test
    fun `findOrder_whenActorMetadataMissing_returnsUnauthenticated`(): Unit = runBlocking {
        val ex = assertFailsWith<StatusRuntimeException> {
            stub(actorUid = null).findOrder(findRequest())
        }
        assertEquals(Status.Code.UNAUTHENTICATED, ex.status.code)
        assertEquals("missing or invalid actor identity", ex.status.description)
    }

    @Test
    fun `findOrder_whenActorUidNotUuid_returnsUnauthenticated`(): Unit = runBlocking {
        val ex = assertFailsWith<StatusRuntimeException> {
            stub("not-a-uuid").findOrder(findRequest())
        }
        assertEquals(Status.Code.UNAUTHENTICATED, ex.status.code)
        assertEquals("missing or invalid actor identity", ex.status.description)
    }

    // ─── listOrders ─────────────────────────────────────────────────────────────

    @Test
    fun `listOrders_whenActorIsNotAdmin_pinsFilterToActor`(): Unit = runBlocking {
        given(useCase.findAllOrders(OrderQueryOption(userUid = ownerUid), expectedPageable, auditContext))
            .willReturn(PageImpl(listOf(orderAggregate(userUid = ownerUid)), expectedPageable, 1L))

        val response = stub(ownerUid.toString()).listOrders(listRequest(userUid = otherUid.toString()))

        assertEquals(1, response.totalSize)
        assertEquals(ownerUid.toString(), response.getOrders(0).userUid)
    }

    @Test
    fun `listOrders_whenActorIsAdmin_honorsRequestFilter`(): Unit = runBlocking {
        given(useCase.findAllOrders(OrderQueryOption(userUid = otherUid), expectedPageable, auditContext))
            .willReturn(PageImpl(listOf(orderAggregate(userUid = otherUid)), expectedPageable, 1L))

        val response = stub(ownerUid.toString(), admin = true)
            .listOrders(listRequest(userUid = otherUid.toString()))

        assertEquals(1, response.totalSize)
        assertEquals(otherUid.toString(), response.getOrders(0).userUid)
    }

    @Test
    fun `listOrders_whenActorIsAdminAndUserUidOmitted_returnsEveryOwnersOrders`(): Unit = runBlocking {
        given(useCase.findAllOrders(OrderQueryOption(userUid = null), expectedPageable, auditContext))
            .willReturn(PageImpl(emptyList(), expectedPageable, 0L))

        val response = stub(ownerUid.toString(), admin = true).listRequestAll()

        assertEquals(0, response.totalSize)
    }

    @Test
    fun `listOrders_whenActorMetadataMissing_returnsUnauthenticated`(): Unit = runBlocking {
        val ex = assertFailsWith<StatusRuntimeException> {
            stub(actorUid = null).listOrders(listRequest(userUid = ""))
        }
        assertEquals(Status.Code.UNAUTHENTICATED, ex.status.code)
    }

    // ─── Request/fixture helpers ──────────────────────────────────────────────

    private fun findRequest(): FindOrderRequest =
        FindOrderRequest.newBuilder().setUid(OrderTestFixtures.ORDER_UID.toString()).build()

    private fun listRequest(userUid: String): ListOrdersRequest =
        ListOrdersRequest.newBuilder().setUserUid(userUid).setPageSize(PAGE_SIZE).build()

    private fun OrderServiceGrpc.OrderServiceBlockingStub.listRequestAll() =
        listOrders(ListOrdersRequest.newBuilder().setPageSize(PAGE_SIZE).build())

    private fun orderAggregate(
        userUid: UUID,
        status: OrderStatus = OrderStatus.PENDING,
    ): OrderAggregate = OrderAggregate(
        id = 1L,
        uid = OrderTestFixtures.ORDER_UID,
        userUid = userUid,
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
        items = emptyList(),
        adjustments = emptyList(),
        shipping = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    companion object {
        private const val PAGE_SIZE = 20
    }
}
