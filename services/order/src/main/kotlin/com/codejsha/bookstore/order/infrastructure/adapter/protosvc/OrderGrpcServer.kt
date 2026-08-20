package com.codejsha.bookstore.order.infrastructure.adapter.protosvc

import com.codejsha.bookstore.generated.application.port.pb.orderpb.*
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.order.infrastructure.support.auth.currentGrpcActor
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import com.codejsha.platform.shared.data.buildPageRequest
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.grpc.stub.StreamObserver
import kotlinx.coroutines.runBlocking
import org.springframework.stereotype.Component
import java.util.UUID
import com.codejsha.bookstore.order.domain.constant.OrderStatus as DomainOrderStatus

@Component
class OrderGrpcServer(
    private val orderUseCase: OrderUseCase
) : OrderServiceGrpc.OrderServiceImplBase() {
    override fun listOrders(
        request: ListOrdersRequest,
        responseObserver: StreamObserver<ListOrdersResponse>,
    ) {
        val actor = currentGrpcActor()
        val statusFilter = request.status
            .takeIf { it != OrderStatus.ORDER_STATUS_UNSPECIFIED }
            ?.let { it.toDomainStatus().value }
        val userUidFilter = if (actor.isAdmin) {
            request.userUid.takeIf { it.isNotEmpty() }?.let { UUID.fromString(it) }
        } else {
            actor.uid
        }
        val option = OrderQueryOption(userUid = userUidFilter, status = statusFilter)
        val page = request.pageToken.takeIf { it.isNotBlank() }?.toIntOrNull() ?: 0
        val pageable = buildPageRequest(request.pageSize, page, request.orderBy)
        val context = ActorContext(actorId = 0L, ActorType.USER)

        try {
            val orders = runBlocking { orderUseCase.findAllOrders(option, pageable, context) }
            val responseBuilder = ListOrdersResponse.newBuilder()
                .addAllOrders(orders.content.map { it.toOrderProto() })
                .setTotalSize(orders.totalElements.toInt())
            if (orders.hasNext()) {
                responseBuilder.nextPageToken = (page + 1).toString()
            }

            responseObserver.onNext(responseBuilder.build())
            responseObserver.onCompleted()
        } catch (t: Throwable) {
            responseObserver.onError(StatusRuntimeException(t.toGrpcStatus()))
        }
    }

    override fun findOrder(
        request: FindOrderRequest,
        responseObserver: StreamObserver<FindOrderResponse>,
    ) {
        try {
            val actor = currentGrpcActor()
            val context = ActorContext(actorId = 0L, ActorType.USER)
            val orderUid = UUID.fromString(request.uid)
            val order = runBlocking { orderUseCase.findOrder(orderUid, context) }
            if (!actor.isAdmin && order.userUid != actor.uid) {
                throw NoSuchElementException("Order with uid $orderUid not found")
            }
            val response = FindOrderResponse.newBuilder()
                .setOrder(order.toOrderProto())
                .build()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (t: Throwable) {
            responseObserver.onError(StatusRuntimeException(t.toGrpcStatus()))
        }
    }

    private fun Throwable.toGrpcStatus(): Status = when (this) {
        is NoSuchElementException -> Status.NOT_FOUND.withDescription(message)
        is IllegalArgumentException -> Status.INVALID_ARGUMENT.withDescription(message)
        else -> Status.INTERNAL.withDescription("internal error")
    }

    private fun OrderAggregate.toOrderProto(): Order =
        Order.newBuilder()
            .setUid(uid.toString())
            .setUserUid(userUid.toString())
            .setOrderNumber(orderNumber)
            .setStatus(status.toProtoStatus())
            .setCurrency(currency)
            .setItemsAmount(itemsAmount.toDouble())
            .setDiscountAmount(discountAmount.toDouble())
            .setShippingAmount(shippingAmount.toDouble())
            .setTaxAmount(taxAmount.toDouble())
            .setTotalAmount(totalAmount.toDouble())
            .addAllLineItems(items.map { item ->
                OrderLineItem.newBuilder()
                    .setProductUid(item.productId.toString())
                    .setQuantity(item.quantity)
                    .setCurrency(item.currency)
                    .setPrice(item.price.toDouble())
                    .build()
            })
            .build()

    private fun DomainOrderStatus.toProtoStatus(): OrderStatus =
        when (this) {
            DomainOrderStatus.PENDING -> OrderStatus.ORDER_STATUS_PENDING
            DomainOrderStatus.PAID -> OrderStatus.ORDER_STATUS_CONFIRMED
            DomainOrderStatus.SHIPPED -> OrderStatus.ORDER_STATUS_SHIPPED
            DomainOrderStatus.DELIVERED -> OrderStatus.ORDER_STATUS_DELIVERED
            DomainOrderStatus.CANCELLED -> OrderStatus.ORDER_STATUS_CANCELLED
            DomainOrderStatus.REFUNDED -> OrderStatus.ORDER_STATUS_REFUNDED
        }

    private fun OrderStatus.toDomainStatus(): DomainOrderStatus =
        when (this) {
            OrderStatus.ORDER_STATUS_PENDING -> DomainOrderStatus.PENDING
            OrderStatus.ORDER_STATUS_CONFIRMED -> DomainOrderStatus.PAID
            OrderStatus.ORDER_STATUS_PROCESSING -> DomainOrderStatus.PAID
            OrderStatus.ORDER_STATUS_SHIPPED -> DomainOrderStatus.SHIPPED
            OrderStatus.ORDER_STATUS_DELIVERED -> DomainOrderStatus.DELIVERED
            OrderStatus.ORDER_STATUS_CANCELLED -> DomainOrderStatus.CANCELLED
            OrderStatus.ORDER_STATUS_REFUNDED -> DomainOrderStatus.REFUNDED
            OrderStatus.ORDER_STATUS_UNSPECIFIED, OrderStatus.UNRECOGNIZED ->
                error("invalid status: $this")
        }
}
