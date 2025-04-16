package com.codejsha.bookstore.payment.infrastructure.adapter.conductor

import com.codejsha.bookstore.payment.domain.model.PaymentDto
import com.codejsha.bookstore.payment.domain.service.PaymentService
import com.netflix.conductor.common.metadata.tasks.TaskResult
import com.netflix.conductor.sdk.workflow.task.WorkerTask
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.stereotype.Component

@Component
class PaymentWorker(
    private val paymentService: PaymentService
) {

    @WorkerTask(value = "make_payment")
    @WithSpan
    fun makePayment(paymentDto: PaymentDto): TaskResult {
        val payment = paymentService.createPayment(paymentDto)
        val output = payment.toMap()
        val result = TaskResult()
        result.outputData = output
        result.status = TaskResult.Status.COMPLETED
        return result
    }
}
