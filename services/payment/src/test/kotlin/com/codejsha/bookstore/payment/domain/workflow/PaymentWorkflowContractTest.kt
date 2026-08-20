package com.codejsha.bookstore.payment.domain.workflow

import com.codejsha.bookstore.payment.infrastructure.support.TemporalConfig
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
    fun `ProcessPaymentRequest deserializes the order workflow's canonical payload`() {
        val json = """
            {
              "orderUid": "44444444-4444-4444-4444-444444444444",
              "userUid": "22222222-2222-2222-2222-222222222222",
              "amount": 25.00,
              "currency": "USD"
            }
        """.trimIndent()
        val payload = Payload.newBuilder()
            .putMetadata("encoding", ByteString.copyFromUtf8("json/plain"))
            .setData(ByteString.copyFromUtf8(json))
            .build()

        val request = converter.fromPayload(payload, ProcessPaymentRequest::class.java, ProcessPaymentRequest::class.java)

        assertEquals("44444444-4444-4444-4444-444444444444", request.orderUid)
        assertEquals("22222222-2222-2222-2222-222222222222", request.userUid)
        assertEquals(0, BigDecimal("25.00").compareTo(request.amount))
        assertEquals("USD", request.currency)
    }

    @Test
    fun `ProcessPaymentResult serializes with the wire field names the order workflow expects`() {
        val payload = converter.toPayload(
            ProcessPaymentResult(
                paymentUid = "99999999-9999-9999-9999-999999999999",
                status = "succeeded",
                gatewayPaymentId = "pay_hyper_1",
            ),
        ).get()

        val node = mapper.readTree(payload.data.toStringUtf8())
        assertEquals(
            setOf("paymentUid", "status", "gatewayPaymentId"),
            node.fieldNames().asSequence().toSet(),
        )
        assertEquals("99999999-9999-9999-9999-999999999999", node.get("paymentUid").asText())
    }
}
