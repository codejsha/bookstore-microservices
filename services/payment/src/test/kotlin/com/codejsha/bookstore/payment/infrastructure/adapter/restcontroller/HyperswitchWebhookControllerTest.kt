package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.application.usecase.WebhookUseCase
import com.codejsha.bookstore.payment.domain.constant.WebhookOutcome
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.springframework.http.HttpStatus
import kotlin.test.assertEquals
import kotlin.test.assertNotNull

class HyperswitchWebhookControllerTest {

    private val systemContext = ActorContext(actorId = 0L, ActorType.SYSTEM)

    @Test
    fun `receive_whenWebhookArrives_handsRawBodyAndSignatureToUsecaseAndAcksOutcome`(): Unit = runBlocking {
        val useCase = mock(WebhookUseCase::class.java)
        val payload = """{"event_id":"evt_1"}""".toByteArray()
        given(useCase.handleHyperswitchWebhook(payload, "abc", systemContext)).willReturn(WebhookOutcome.APPLIED)
        val controller = HyperswitchWebhookController(useCase)

        val response = controller.receive(payload, "abc")

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals("applied", assertNotNull(response.body).result)
    }

    // Both rejection paths of the same handler are asserted here.
    @Test
    fun `receive_whenSignatureInvalid_returns401AndOtherRejections400`() {
        val controller = HyperswitchWebhookController(mock(WebhookUseCase::class.java))

        val unauthorized = controller.handleWebhookRejected(HyperswitchWebhookException("Hyperswitch webhook signature is missing or invalid"))
        val badRequest = controller.handleWebhookRejected(HyperswitchWebhookException("Hyperswitch webhook payload has no event_id"))

        assertEquals(HttpStatus.UNAUTHORIZED, unauthorized.statusCode)
        assertEquals(HttpStatus.BAD_REQUEST, badRequest.statusCode)
        assertEquals("rejected", assertNotNull(badRequest.body).result)
    }
}
