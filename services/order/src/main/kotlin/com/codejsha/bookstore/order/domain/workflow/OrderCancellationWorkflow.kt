package com.codejsha.bookstore.order.domain.workflow

import io.temporal.workflow.WorkflowInterface
import io.temporal.workflow.WorkflowMethod

@WorkflowInterface
interface OrderCancellationWorkflow {
    @WorkflowMethod
    fun cancelOrder(request: CancelOrderWorkflowRequest): CancelOrderWorkflowResult
}
