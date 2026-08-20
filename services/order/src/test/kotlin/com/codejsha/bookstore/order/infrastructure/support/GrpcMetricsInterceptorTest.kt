package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.generated.application.port.pb.orderpb.FindOrderRequest
import com.codejsha.bookstore.generated.application.port.pb.orderpb.OrderServiceGrpc
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.infrastructure.adapter.protosvc.OrderGrpcServer
import com.codejsha.bookstore.order.infrastructure.support.auth.ACTOR_UID_METADATA_KEY
import com.codejsha.bookstore.order.infrastructure.support.auth.ActorAuthInterceptor
import com.codejsha.bookstore.order.support.OrderTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import io.grpc.ManagedChannel
import io.grpc.Metadata
import io.grpc.Server
import io.grpc.ServerInterceptors
import io.grpc.StatusRuntimeException
import io.grpc.inprocess.InProcessChannelBuilder
import io.grpc.inprocess.InProcessServerBuilder
import io.grpc.stub.MetadataUtils
import io.micrometer.core.instrument.simple.SimpleMeterRegistry
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import kotlinx.coroutines.runBlocking
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class GrpcMetricsInterceptorTest {

    private val registry = SimpleMeterRegistry()
    private val useCase = mock(OrderUseCase::class.java)

    private lateinit var server: Server
    private lateinit var channel: ManagedChannel

    @BeforeEach
    fun startServer() {
        val name = InProcessServerBuilder.generateName()
        server = InProcessServerBuilder.forName(name)
            .directExecutor()
            .addService(
                ServerInterceptors.intercept(
                    OrderGrpcServer(useCase),
                    ActorAuthInterceptor(),
                    GrpcMetricsInterceptor(registry),
                ),
            )
            .build()
            .start()
        channel = InProcessChannelBuilder.forName(name).directExecutor().build()
    }

    @AfterEach
    fun stopServer() {
        channel.shutdownNow()
        server.shutdownNow()
    }

    private fun stub(actorUid: String?): OrderServiceGrpc.OrderServiceBlockingStub {
        val md = Metadata()
        if (actorUid != null) md.put(ACTOR_UID_METADATA_KEY, actorUid)
        return OrderServiceGrpc.newBlockingStub(channel)
            .withInterceptors(MetadataUtils.newAttachHeadersInterceptor(md))
    }

    private fun findRequest(): FindOrderRequest =
        FindOrderRequest.newBuilder().setUid(OrderTestFixtures.ORDER_UID.toString()).build()

    @Test
    fun `records rpc server duration tagged with service, method and status`(): Unit = runBlocking {
        val ownerUid = OrderTestFixtures.USER_UID
        given(useCase.findOrder(OrderTestFixtures.ORDER_UID, ActorContext(actorId = 0L, actorType = ActorType.USER)))
            .willReturn(orderAggregate(ownerUid))

        stub(ownerUid.toString()).findOrder(findRequest())

        val timer = registry.find("rpc.server.duration").timers().single()
        assertEquals(1L, timer.count())
        assertEquals("grpc", timer.id.getTag("rpc.system"))
        assertEquals("order.v1.OrderService", timer.id.getTag("rpc.service"))
        assertEquals("FindOrder", timer.id.getTag("rpc.method"))
        assertEquals("0", timer.id.getTag("rpc.grpc.status_code"))
    }

    private fun orderAggregate(userUid: UUID): OrderAggregate = OrderAggregate(
        id = 1L,
        uid = OrderTestFixtures.ORDER_UID,
        userUid = userUid,
        orderNumber = "ORD-20260101-0001",
        status = OrderStatus.PENDING,
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

    @Test
    fun `records an UNAUTHENTICATED call rejected before the service is reached`() {
        assertFailsWith<StatusRuntimeException> { stub(null).findOrder(findRequest()) }

        val timer = registry.find("rpc.server.duration").timers().single()
        assertEquals(1L, timer.count())
        assertEquals("16", timer.id.getTag("rpc.grpc.status_code"))
    }
}
