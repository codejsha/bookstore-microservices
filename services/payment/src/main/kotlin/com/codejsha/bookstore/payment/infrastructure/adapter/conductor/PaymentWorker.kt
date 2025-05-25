package com.codejsha.bookstore.payment.infrastructure.adapter.conductor

import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.handler.PaymentHandler
import com.codejsha.bookstore.payment.domain.model.MakePaymentPayload
import com.netflix.conductor.common.metadata.tasks.TaskResult
import com.netflix.conductor.sdk.workflow.task.WorkerTask
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.stereotype.Component

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "command")
class PaymentWorker(
    private val paymentHandler: PaymentHandler,
    private val paymentUseCase: PaymentUseCase
) {

    @WorkerTask(value = "make_payment")
    @WithSpan
    fun makePaymentTask(payload: MakePaymentPayload): TaskResult {
        val payment = paymentUseCase.createPayment(payload.toPaymentDto())
        val output = payment.toMap()
        val result = TaskResult()
        result.outputData = output
        result.status = TaskResult.Status.COMPLETED
        return result
    }
}
