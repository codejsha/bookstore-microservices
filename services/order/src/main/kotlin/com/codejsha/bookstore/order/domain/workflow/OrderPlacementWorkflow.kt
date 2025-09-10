package com.codejsha.bookstore.order.domain.workflow

import io.temporal.workflow.WorkflowInterface
import io.temporal.workflow.WorkflowMethod

@WorkflowInterface
interface OrderPlacementWorkflow {
    @WorkflowMethod
    fun placeOrder(request: PlaceOrderWorkflowRequest): PlaceOrderWorkflowResult
}
