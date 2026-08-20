package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.application.usecase.WebhookUseCase
import com.codejsha.bookstore.payment.domain.constant.WebhookOutcome
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.slf4j.LoggerFactory
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestHeader
import org.springframework.web.bind.annotation.RestController

@RestController
class HyperswitchWebhookController(
    private val webhookUseCase: WebhookUseCase,
) {

    private val log = LoggerFactory.getLogger(HyperswitchWebhookController::class.java)

    @PostMapping(WEBHOOK_PATH)
    fun receive(
        @RequestBody payload: ByteArray,
        @RequestHeader(name = SIGNATURE_HEADER, required = false) signature: String?,
    ): ResponseEntity<WebhookAckResponse> = runBlocking {
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        val outcome = webhookUseCase.handleHyperswitchWebhook(payload, signature, context)
        ResponseEntity.ok(WebhookAckResponse(outcome.name.lowercase()))
    }

    @ExceptionHandler(HyperswitchWebhookException::class)
    fun handleWebhookRejected(e: HyperswitchWebhookException): ResponseEntity<WebhookAckResponse> {
        log.warn("Rejected Hyperswitch webhook: {}", e.message)
        val status = if (e.message?.contains("signature") == true) HttpStatus.UNAUTHORIZED else HttpStatus.BAD_REQUEST
        return ResponseEntity.status(status).body(WebhookAckResponse("rejected"))
    }

    data class WebhookAckResponse(val result: String)

    companion object {
        const val WEBHOOK_PATH = "/internal/webhooks/hyperswitch"
        const val SIGNATURE_HEADER = "x-webhook-signature-512"
    }
}
