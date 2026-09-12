package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.order.application.usecase.CartUseCase
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.CartItemEntity
import com.codejsha.bookstore.order.domain.model.command.InvalidCommandException
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.domain.workflow.CancelOrderWorkflowRequest
import com.codejsha.bookstore.order.domain.workflow.FulfillOrderWorkflowRequest
import com.codejsha.bookstore.order.domain.workflow.OrderCancellationWorkflow
import com.codejsha.bookstore.order.domain.workflow.OrderFulfillmentWorkflow
import com.codejsha.bookstore.order.domain.workflow.OrderPlacementWorkflow
import com.codejsha.bookstore.order.domain.workflow.PlaceOrderWorkflowRequest
import com.codejsha.bookstore.order.domain.workflow.ShipmentDestination
import com.codejsha.bookstore.order.infrastructure.support.auth.AuthPrincipal
import com.codejsha.bookstore.order.infrastructure.support.auth.BadRequestException
import com.codejsha.bookstore.order.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.order.infrastructure.support.auth.Principal
import com.codejsha.bookstore.order.infrastructure.support.auth.assertOrderOwner
import com.codejsha.bookstore.order.infrastructure.support.auth.isStaff
import com.codejsha.bookstore.order.infrastructure.support.auth.orUnauthorized
import com.codejsha.bookstore.order.infrastructure.support.auth.subjectUserUid
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import io.temporal.client.WorkflowClient
import io.temporal.client.WorkflowExecutionAlreadyStarted
import io.temporal.api.enums.v1.WorkflowIdReusePolicy
import io.temporal.client.WorkflowOptions
import kotlinx.coroutines.runBlocking
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RestController
import java.math.BigDecimal
import java.util.UUID

@RestController
class OrderWorkflowController(
    private val workflowClient: WorkflowClient,
    private val orderUseCase: OrderUseCase,
    private val cartUseCase: CartUseCase,
) {

    companion object {
        private const val ORDER_TASK_QUEUE = "order-task-queue"
    }

    @PostMapping("/api/v1/orders/saga/place")
    fun placeOrderSaga(
        @RequestBody body: PlaceOrderWorkflowRequest,
        @AuthPrincipal principal: Principal?,
    ): ResponseEntity<WorkflowTriggerResponse> {
        val actor = principal.orUnauthorized()
        val ownerUserUid: String =
            if (actor.isStaff()) body.userUid
            else actor.subjectUserUid()?.toString()
                ?: throw ForbiddenException("subject '${actor.sub}' is not a valid user uid")
        val cartItems = runBlocking {
            val context = ActorContext(actorId = 0L, ActorType.USER)
            cartUseCase.getCart(UUID.fromString(ownerUserUid), context).items
        }
        val request = repriceFromCart(body.copy(userUid = ownerUserUid), cartItems)
        validatePlaceOrderRequest(request)
        val workflowId = "place-${request.idempotencyKey}"
        val stub = workflowClient.newWorkflowStub(
            OrderPlacementWorkflow::class.java,
            WorkflowOptions.newBuilder()
                .setTaskQueue(ORDER_TASK_QUEUE)
                .setWorkflowId(workflowId)
                .setWorkflowIdReusePolicy(WorkflowIdReusePolicy.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY)
                .build(),
        )
        return acceptedOrExisting(workflowId) { WorkflowClient.start(stub::placeOrder, request) }
    }

    @PostMapping("/api/v1/orders/{uid}/cancellation")
    fun cancelOrderSaga(
        @PathVariable uid: String,
        @RequestBody body: CancelOrderTriggerRequest,
        @AuthPrincipal principal: Principal?,
    ): ResponseEntity<WorkflowTriggerResponse> {
        val actor = principal.orUnauthorized()
        actor.assertOrderOwner(loadOrderOwner(uid, actor))
        val workflowId = "cancel-$uid"
        val stub = workflowClient.newWorkflowStub(
            OrderCancellationWorkflow::class.java,
            WorkflowOptions.newBuilder()
                .setTaskQueue(ORDER_TASK_QUEUE)
                .setWorkflowId(workflowId)
                .setWorkflowIdReusePolicy(WorkflowIdReusePolicy.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY)
                .build(),
        )
        val request = CancelOrderWorkflowRequest(
            orderUid = uid,
            reason = body.reason,
        )
        return acceptedOrExisting(workflowId) { WorkflowClient.start(stub::cancelOrder, request) }
    }

    @PostMapping("/api/v1/orders/{uid}/fulfill")
    fun fulfillOrderSaga(
        @PathVariable uid: String,
        @RequestBody body: FulfillOrderTriggerRequest,
        @AuthPrincipal principal: Principal?,
    ): ResponseEntity<WorkflowTriggerResponse> {
        val actor = principal.orUnauthorized()
        val ownerUserUid = loadOrderOwner(uid, actor)
        actor.assertOrderOwner(ownerUserUid)
        val workflowId = "fulfill-$uid"
        val stub = workflowClient.newWorkflowStub(
            OrderFulfillmentWorkflow::class.java,
            WorkflowOptions.newBuilder()
                .setTaskQueue(ORDER_TASK_QUEUE)
                .setWorkflowId(workflowId)
                .setWorkflowIdReusePolicy(WorkflowIdReusePolicy.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY)
                .build(),
        )
        val request = FulfillOrderWorkflowRequest(
            orderUid = uid,
            userUid = ownerUserUid.toString(),
            destination = body.destination,
        )
        return acceptedOrExisting(workflowId) { WorkflowClient.start(stub::fulfillOrder, request) }
    }

    private fun loadOrderOwner(uid: String, principal: Principal): UUID = runBlocking {
        val context = ActorContext(actorId = 0L, ActorType.USER)
        orderUseCase.findOrder(UUID.fromString(uid), context).userUid
    }

    private inline fun acceptedOrExisting(
        workflowId: String,
        start: () -> io.temporal.api.common.v1.WorkflowExecution,
    ): ResponseEntity<WorkflowTriggerResponse> {
        val id = try {
            start().workflowId
        } catch (e: WorkflowExecutionAlreadyStarted) {
            workflowId
        }
        return ResponseEntity.accepted().body(WorkflowTriggerResponse(workflowId = id))
    }
}

internal fun repriceFromCart(
    request: PlaceOrderWorkflowRequest,
    cartItems: List<CartItemEntity>,
): PlaceOrderWorkflowRequest {
    if (request.items.isEmpty()) {
        throw BadRequestException("order must contain at least one item")
    }
    val cartByProduct = cartItems.associateBy { it.productId }
    val repriced = request.items.map { item ->
        val cartItem = cartByProduct[item.productId]
            ?: throw BadRequestException(
                "no server-side price for product ${item.productId}; add it to the cart before placing the order",
            )
        item.copy(price = cartItem.price, currency = cartItem.currency)
    }
    return request.copy(items = repriced)
}

internal fun validatePlaceOrderRequest(request: PlaceOrderWorkflowRequest) {
    val userUid = runCatching { UUID.fromString(request.userUid) }.getOrNull()
        ?: throw InvalidCommandException("user_uid must be a valid uuid, was ${request.userUid}")

    OrderCreateCommand(
        userUid = userUid,
        currency = request.currency,
        itemsAmount = BigDecimal.ZERO,
        discountAmount = BigDecimal.ZERO,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = BigDecimal.ZERO,
        idempotencyKey = request.idempotencyKey,
    )

    request.items.forEach { item ->
        OrderItemCreateCommand(
            productId = item.productId,
            sku = null,
            productName = item.productName,
            options = null,
            quantity = item.quantity,
            currency = item.currency,
            price = item.price,
            taxRate = BigDecimal.ZERO,
        )
    }

    request.shipping?.let { shipping ->
        OrderShippingCreateCommand(
            recipientName = shipping.recipientName,
            recipientPhone = shipping.recipientPhone,
            addressLine1 = shipping.addressLine1,
            addressLine2 = shipping.addressLine2,
            city = shipping.city,
            state = shipping.state,
            postalCode = shipping.postalCode,
            country = shipping.country,
            shippingMethod = shipping.shippingMethod,
        )
    }
}

data class CancelOrderTriggerRequest(
    val reason: String? = null,
)

data class FulfillOrderTriggerRequest(
    val destination: ShipmentDestination,
)

data class WorkflowTriggerResponse(
    val workflowId: String,
)
