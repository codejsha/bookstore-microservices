package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.generated.application.port.pb.orderpb.FindOrderRequest
import com.codejsha.bookstore.generated.application.port.pb.orderpb.OrderServiceGrpc
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.config.properties.GrpcConfig
import com.codejsha.bookstore.order.config.properties.GrpcServerConfig
import com.codejsha.bookstore.order.infrastructure.adapter.protosvc.OrderGrpcServer
import com.codejsha.bookstore.order.infrastructure.support.auth.ActorAuthInterceptor
import com.codejsha.bookstore.order.support.OrderTestFixtures
import io.grpc.ManagedChannelBuilder
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.micrometer.core.instrument.simple.SimpleMeterRegistry
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import java.net.ServerSocket
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class GrpcServerTest {

    private val useCase = mock(OrderUseCase::class.java)
    private val port = ServerSocket(0).use { it.localPort }
    private val grpcServer = GrpcServer(
        GrpcConfig(GrpcServerConfig(host = "localhost", port = port)),
        OrderGrpcServer(useCase),
        ActorAuthInterceptor(),
        GrpcMetricsInterceptor(SimpleMeterRegistry()),
    )

    @AfterEach
    fun stopServer() {
        grpcServer.stop()
    }

    @Test
    fun `isRunning_neverStarted_reportsStopped`() {
        assertFalse(grpcServer.isRunning)
    }

    @Test
    fun `stop_neverStarted_staysStopped`() {
        grpcServer.stop()

        assertFalse(grpcServer.isRunning)
    }

    @Test
    fun `start_configuredPort_servesOrderServiceBehindActorInterceptor`() {
        grpcServer.start()
        val channel = ManagedChannelBuilder.forAddress("localhost", port).usePlaintext().build()

        try {
            val ex = assertFailsWith<StatusRuntimeException> {
                OrderServiceGrpc.newBlockingStub(channel).findOrder(
                    FindOrderRequest.newBuilder().setUid(OrderTestFixtures.ORDER_UID.toString()).build(),
                )
            }

            assertTrue(grpcServer.isRunning)
            assertEquals(Status.Code.UNAUTHENTICATED, ex.status.code)
            verifyNoInteractions(useCase)
        } finally {
            channel.shutdownNow()
        }
    }

    @Test
    fun `start_alreadyRunning_keepsTheBoundServer`() {
        grpcServer.start()

        grpcServer.start()

        assertTrue(grpcServer.isRunning)
    }

    @Test
    fun `stop_running_shutsDownAndAllowsRestart`() {
        grpcServer.start()

        grpcServer.stop()
        assertFalse(grpcServer.isRunning)

        grpcServer.start()
        assertTrue(grpcServer.isRunning)
    }

    @Test
    fun `getPhase_queried_startsJustBeforeTheLastPhase`() {
        assertEquals(Int.MAX_VALUE - 1, grpcServer.phase)
    }
}
