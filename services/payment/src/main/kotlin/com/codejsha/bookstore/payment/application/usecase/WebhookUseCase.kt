package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.constant.WebhookOutcome
import com.codejsha.platform.shared.data.ActorContext

interface WebhookUseCase {
    suspend fun handleHyperswitchWebhook(payload: ByteArray, signature: String?, context: ActorContext): WebhookOutcome
}
