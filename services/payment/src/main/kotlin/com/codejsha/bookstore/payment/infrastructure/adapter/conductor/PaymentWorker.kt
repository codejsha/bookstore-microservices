package com.codejsha.bookstore.payment.infrastructure.adapter.conductor

import com.codejsha.bookstore.payment.application.port.saga.MakePaymentPayload
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.handler.PaymentHandler
import com.codejsha.bookstore.payment.domain.model.PaymentMapper
import com.netflix.conductor.common.metadata.tasks.TaskResult
import com.netflix.conductor.sdk.workflow.task.WorkerTask
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.stereotype.Component

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "command")
class PaymentWorker(
    private val paymentUseCase: PaymentUseCase,
    private val paymentHandler: PaymentHandler
) {

    @WorkerTask(value = "make_payment")
    @WithSpan
    fun makePaymentTask(payload: MakePaymentPayload): TaskResult {
        val dto = PaymentMapper.toPaymentDto(payload)
        val payment = paymentUseCase.createPayment(dto)
        val output: Map<String, Any?> = mapOf(
            "order_id" to payment.orderId,
        )

        val result = TaskResult()
        result.outputData = output
        result.status = TaskResult.Status.COMPLETED
        return result
    }
}
