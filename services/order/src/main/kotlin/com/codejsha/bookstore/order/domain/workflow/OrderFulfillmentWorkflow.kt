package com.codejsha.bookstore.order.domain.workflow

import io.temporal.workflow.WorkflowInterface
import io.temporal.workflow.WorkflowMethod

@WorkflowInterface
interface OrderFulfillmentWorkflow {
    @WorkflowMethod
    fun fulfillOrder(request: FulfillOrderWorkflowRequest): FulfillOrderWorkflowResult
}
