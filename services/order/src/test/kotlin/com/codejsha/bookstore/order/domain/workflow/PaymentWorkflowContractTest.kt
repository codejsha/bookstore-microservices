package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.infrastructure.support.TemporalConfig
import com.fasterxml.jackson.databind.ObjectMapper
import com.google.protobuf.ByteString
import io.temporal.api.common.v1.Payload
import org.junit.jupiter.api.Test
import java.math.BigDecimal
import kotlin.test.assertEquals

class PaymentWorkflowContractTest {

    private val converter = TemporalConfig().temporalDataConverter()
    private val mapper = ObjectMapper()

    @Test
    fun `ProcessPaymentRequest serializes with the wire field names the payment worker expects`() {
        val payload = converter.toPayload(
            ProcessPaymentRequest(
                orderUid = "44444444-4444-4444-4444-444444444444",
                userUid = "22222222-2222-2222-2222-222222222222",
                amount = BigDecimal("25.00"),
                currency = "USD",
            ),
        ).get()

        val node = mapper.readTree(payload.data.toStringUtf8())
        assertEquals(
            setOf("orderUid", "userUid", "amount", "currency"),
            node.fieldNames().asSequence().toSet(),
        )
    }

    @Test
    fun `ProcessPaymentResult deserializes the payment worker's canonical response, tolerating unknown fields`() {
        val json = """
            {
              "paymentUid": "99999999-9999-9999-9999-999999999999",
              "status": "succeeded",
              "gatewayPaymentId": "pay_hyper_1",
              "someFieldAddedByANewerPaymentService": true
            }
        """.trimIndent()
        val payload = Payload.newBuilder()
            .putMetadata("encoding", ByteString.copyFromUtf8("json/plain"))
            .setData(ByteString.copyFromUtf8(json))
            .build()

        val result = converter.fromPayload(payload, ProcessPaymentResult::class.java, ProcessPaymentResult::class.java)

        assertEquals("99999999-9999-9999-9999-999999999999", result.paymentUid)
        assertEquals("succeeded", result.status)
        assertEquals("pay_hyper_1", result.gatewayPaymentId)
    }

    @Test
    fun `refund activity results deserialize as booleans`() {
        val payload = Payload.newBuilder()
            .putMetadata("encoding", ByteString.copyFromUtf8("json/plain"))
            .setData(ByteString.copyFromUtf8("true"))
            .build()

        val refunded = converter.fromPayload(payload, Boolean::class.java, Boolean::class.java)

        assertEquals(true, refunded)
    }
}
