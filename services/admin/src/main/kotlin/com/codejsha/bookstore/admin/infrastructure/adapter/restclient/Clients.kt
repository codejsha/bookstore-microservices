package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.application.port.restclient.PaymentClient
import com.codejsha.bookstore.admin.application.port.restclient.SettlementClient
import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.external.OrderLine
import com.codejsha.bookstore.admin.domain.model.external.OrderShipping
import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementDetailLine
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.StockAtWarehouse
import com.codejsha.bookstore.admin.domain.model.external.Subject
import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.external.Warehouse
import com.codejsha.bookstore.admin.domain.model.external.Work
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.AuthorApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.SubjectApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.WorkApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.AuthorItem
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.SubjectItem
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkCreateRequest
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkItem
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkUpdateRequest
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkUpdateResponse
import com.codejsha.bookstore.generated.application.port.restclient.identity.api.RiskApi
import com.codejsha.bookstore.generated.application.port.restclient.identity.api.UserApi
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.AuthRole
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.RiskEntryResponse
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.RiskFlagRequest
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.RiskLevel
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.UserFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.UserRolesRequest
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.UserStatus
import com.codejsha.bookstore.generated.application.port.restclient.inventory.api.StockApi
import com.codejsha.bookstore.generated.application.port.restclient.inventory.api.WarehouseApi
import com.codejsha.bookstore.generated.application.port.restclient.inventory.model.StockFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.inventory.model.StockWarehouseItem
import com.codejsha.bookstore.generated.application.port.restclient.inventory.model.WarehouseFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.api.OrderApi
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderItemFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderShippingFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.order.model.OrderStatus
import com.codejsha.bookstore.generated.application.port.restclient.payment.api.PaymentApi
import com.codejsha.bookstore.generated.application.port.restclient.payment.api.RefundApi
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.PaymentFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.PaymentStatus
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.RefundFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.RefundStatus
import com.codejsha.bookstore.generated.application.port.restclient.settlement.api.SettlementApi
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementDetailResponse
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementRunRequest
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementRunResponse
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementStatus
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

private const val COUNT_ONLY_PAGE_SIZE = 1

@Component
class CatalogRestClient(
    private val workApi: WorkApi,
    private val authorApi: AuthorApi,
    private val subjectApi: SubjectApi,
) : CatalogClient {

    override fun countWorks(): Long =
        workApi.worksSearch(null, null, null, null, COUNT_ONLY_PAGE_SIZE, null, null).total

    override fun findAllWorks(option: WorkQueryOption, pageable: Pageable): Page<Work> {
        val response = workApi.worksSearch(
            option.title,
            option.authorUid,
            option.subjectUid,
            null,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toWork() }, pageable, response.total)
    }

    override fun findWork(uid: String): Work = workApi.worksRead(uid).toWork()

    override fun createWork(command: WorkCreateCommand) {
        workApi.worksCreate(
            WorkCreateRequest(
                title = command.title,
                description = command.description,
                firstPublishDate = command.firstPublishDate,
                olKey = command.olKey,
                authorUids = command.authorUids,
                subjectNames = command.subjectNames,
                coverUids = null,
            )
        )
    }

    override fun updateWork(uid: String, command: WorkUpdateCommand): Work =
        workApi.worksUpdate(
            uid,
            WorkUpdateRequest(
                title = command.title,
                description = command.description,
                firstPublishDate = command.firstPublishDate,
                olKey = command.olKey,
                authorUids = command.authorUids,
                subjectNames = command.subjectNames,
                coverUids = null,
            ),
        ).toWork()

    override fun findAllAuthors(name: String?, pageable: Pageable): Page<Author> {
        val response = authorApi.authorsGetAll(
            name,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toAuthor() }, pageable, response.total)
    }

    override fun findAllSubjects(name: String?, pageable: Pageable): Page<Subject> {
        val response = subjectApi.subjectsGetAll(
            name,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toSubject() }, pageable, response.total)
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun AuthorItem.toAuthor() = Author(uid = uid, name = name)

    private fun SubjectItem.toSubject() = Subject(uid = uid, name = name)

    private fun WorkItem.toWork() = Work(
        uid = uid,
        title = title,
        description = description,
        firstPublishDate = firstPublishDate,
        olKey = olKey,
        authors = authors.map { it.toAuthor() },
        subjects = subjects.map { it.toSubject() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun WorkFindResponse.toWork() = Work(
        uid = uid,
        title = title,
        description = description,
        firstPublishDate = firstPublishDate,
        olKey = olKey,
        authors = authors.map { it.toAuthor() },
        subjects = subjects.map { it.toSubject() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun WorkUpdateResponse.toWork() = Work(
        uid = uid,
        title = title,
        description = description,
        firstPublishDate = firstPublishDate,
        olKey = olKey,
        authors = authors.map { it.toAuthor() },
        subjects = subjects.map { it.toSubject() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )
}

@Component
class IdentityRestClient(
    private val userApi: UserApi,
    private val riskApi: RiskApi,
) : IdentityClient {

    override fun listRisk(): List<RiskEntry> = riskApi.riskGetAll().entries.map { it.toRiskEntry() }

    override fun flagRisk(uid: String, level: String, reason: String, ttlSeconds: Long?): RiskEntry =
        riskApi.riskFlagPrincipal(
            uid,
            RiskFlagRequest(level = level.toRiskLevel(), reason = reason, ttlSeconds = ttlSeconds),
        ).toRiskEntry()

    override fun unflagRisk(uid: String) = riskApi.riskUnflagPrincipal(uid)

    override fun findAllUsers(option: UserQueryOption, pageable: Pageable): Page<User> {
        val response = userApi.usersGetAll(
            option.email,
            option.name,
            option.phone,
            option.status?.toUserStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toUser() }, pageable, response.total)
    }

    override fun findUser(uid: String): User = userApi.usersRead(uid).toUser()

    override fun updateRoles(uid: String, roles: List<String>): User =
        userApi.usersUpdateRoles(uid, UserRolesRequest(roles = roles.map { it.toAuthRole() })).toUser()

    override fun suspendUser(uid: String): User = userApi.usersSuspend(uid).toUser()

    override fun reactivateUser(uid: String): User = userApi.usersReactivate(uid).toUser()

    override fun deactivateUser(uid: String): User = userApi.usersDeactivate(uid).toUser()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun RiskEntryResponse.toRiskEntry(): RiskEntry = RiskEntry(
        userUid = userUid,
        level = level.value,
        reason = reason,
        flaggedBy = flaggedBy,
        flaggedAt = flaggedAt,
        expiresAt = expiresAt,
    )

    private fun String.toRiskLevel(): RiskLevel =
        try {
            RiskLevel.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown risk level: $this")
        } catch (_: IllegalArgumentException) {
            throw UnsupportedValueException("unknown risk level: $this")
        }

    private fun String.toUserStatus(): UserStatus =
        try {
            UserStatus.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown user status: $this")
        }

    private fun String.toAuthRole(): AuthRole {
        val role = try {
            AuthRole.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown role: $this")
        }
        if (role == AuthRole.UNKNOWN) throw UnsupportedValueException("unknown role: $this")
        return role
    }

    private fun UserFindResponse.toUser() = User(
        uid = uid,
        email = email,
        firstName = firstName,
        lastName = lastName,
        phone = phone,
        status = status.value,
        roles = roles.map { it.value },
        lastLoginAt = lastLoginAt,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )
}

@Component
class OrderRestClient(
    private val orderApi: OrderApi,
) : OrderClient {
    override fun countOrders(): Long =
        orderApi.ordersGetAll(null, null, COUNT_ONLY_PAGE_SIZE, null, null).total

    override fun findAllOrders(option: OrderQueryOption, pageable: Pageable): Page<Order> {
        val response = orderApi.ordersGetAll(
            option.userUid,
            option.status?.toOrderStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toOrder() }, pageable, response.total)
    }

    override fun findOrder(uid: String): Order = orderApi.ordersRead(uid).toOrder()

    override fun cancelOrder(uid: String): Order = orderApi.ordersCancel(uid).toOrder()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun String.toOrderStatus(): OrderStatus =
        try {
            OrderStatus.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown order status: $this")
        }

    private fun OrderFindResponse.toOrder() = Order(
        uid = uid,
        orderNumber = orderNumber,
        userUid = userUid,
        status = status.value,
        currency = currency,
        itemsAmount = itemsAmount,
        discountAmount = discountAmount,
        shippingAmount = shippingAmount,
        taxAmount = taxAmount,
        totalAmount = totalAmount,
        lines = items.orEmpty().map { it.toLine() },
        shipping = shipping?.toShipping(),
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun OrderItemFindResponse.toLine() = OrderLine(
        uid = uid,
        productId = productId,
        productName = productName,
        sku = sku,
        quantity = quantity,
        price = price,
        subtotal = subtotal,
        taxRate = taxRate,
        currency = currency,
        options = options,
    )

    private fun OrderShippingFindResponse.toShipping() = OrderShipping(
        recipientName = recipientName,
        recipientPhone = recipientPhone,
        addressLine1 = addressLine1,
        addressLine2 = addressLine2,
        city = city,
        state = state,
        postalCode = postalCode,
        country = country,
        shippingMethod = shippingMethod,
    )
}

@Component
class InventoryRestClient(
    private val warehouseApi: WarehouseApi,
    private val stockApi: StockApi,
) : InventoryClient {
    override fun countWarehouses(): Long =
        warehouseApi.warehousesGetAll(null, COUNT_ONLY_PAGE_SIZE, null, null).total

    override fun findAllWarehouses(name: String?, pageable: Pageable): Page<Warehouse> {
        val response = warehouseApi.warehousesGetAll(
            name,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toWarehouse() }, pageable, response.total)
    }

    override fun findWarehouse(uid: String): Warehouse = warehouseApi.warehousesRead(uid).toWarehouse()

    override fun findAllStocks(option: StockQueryOption, pageable: Pageable): Page<Stock> {
        val response = stockApi.stocksGetAll(
            option.editionUid,
            option.warehouseUid,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toStock() }, pageable, response.total)
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun WarehouseFindResponse.toWarehouse() = Warehouse(
        uid = uid,
        name = name,
        address = address,
        capacity = capacity,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun StockFindResponse.toStock() = Stock(
        uid = uid,
        editionUid = editionUid,
        totalQuantity = totalQuantity,
        warehouses = warehouses.map { it.toStockAtWarehouse() },
    )

    private fun StockWarehouseItem.toStockAtWarehouse() = StockAtWarehouse(
        warehouseUid = warehouseUid,
        warehouseName = warehouseName,
        quantity = quantity,
    )
}

@Component
class PaymentRestClient(
    private val paymentApi: PaymentApi,
    private val refundApi: RefundApi,
) : PaymentClient {

    override fun findAllPayments(option: PaymentQueryOption, pageable: Pageable): Page<Payment> {
        val response = paymentApi.paymentsGetAll(
            option.customerId,
            option.status?.toPaymentStatus(),
            option.connector,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toPayment() }, pageable, response.total)
    }

    override fun findPayment(uid: String): Payment = paymentApi.paymentsRead(uid).toPayment()

    override fun findAllRefunds(option: RefundQueryOption, pageable: Pageable): Page<Refund> {
        val response = refundApi.refundsGetAll(
            option.paymentId,
            option.status?.toRefundStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toRefund() }, pageable, response.total)
    }

    override fun findRefund(uid: String): Refund = refundApi.refundsRead(uid).toRefund()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun String.toPaymentStatus(): PaymentStatus =
        try {
            PaymentStatus.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown payment status: $this")
        }

    private fun String.toRefundStatus(): RefundStatus =
        try {
            RefundStatus.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown refund status: $this")
        }

    private fun PaymentFindResponse.toPayment() = Payment(
        uid = uid,
        paymentId = paymentId,
        customerId = customerId,
        status = status.value,
        amount = amount,
        amountCaptured = amountCaptured,
        amountCapturable = amountCapturable,
        currency = currency,
        paymentMethod = paymentMethod?.value,
        connector = connector,
        errorCode = errorCode,
        errorMessage = errorMessage,
        confirmedAt = confirmedAt,
        capturedAt = capturedAt,
        cancelledAt = cancelledAt,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun RefundFindResponse.toRefund() = Refund(
        uid = uid,
        refundId = refundId,
        paymentId = paymentId,
        status = status.value,
        refundType = refundType.value,
        amount = amount,
        currency = currency,
        reason = reason,
        connector = connector,
        errorCode = errorCode,
        errorMessage = errorMessage,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )
}

@Component
class SettlementRestClient(
    private val settlementApi: SettlementApi,
) : SettlementClient {

    override fun findAllSettlements(option: SettlementQueryOption, pageable: Pageable): Page<SettlementBucket> {
        val response = settlementApi.settlementsGetAll(
            option.settlementDate,
            option.status?.toSettlementStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toBucket() }, pageable, response.total)
    }

    override fun findSettlement(uid: String): SettlementBucket =
        settlementApi.settlementsRead(uid).toBucket()

    override fun triggerSettlementRun(command: TriggerSettlementRunCommand): SettlementRunAck =
        settlementApi.settlementRunsCreate(
            SettlementRunRequest(targetDate = command.targetDate, rerun = command.rerun),
        ).toAck()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun String.toSettlementStatus(): SettlementStatus =
        try {
            SettlementStatus.fromValue(this)
        } catch (_: NoSuchElementException) {
            throw UnsupportedValueException("unknown settlement status: $this")
        }

    private fun SettlementFindResponse.toBucket() = SettlementBucket(
        uid = uid,
        settlementDate = settlementDate,
        currency = currency,
        paymentMethod = paymentMethod,
        grossAmount = grossAmount,
        refundAmount = refundAmount,
        feeAmount = feeAmount,
        netAmount = netAmount,
        paymentCount = paymentCount,
        refundCount = refundCount,
        status = status.value,
        details = details.orEmpty().map { it.toLine() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun SettlementDetailResponse.toLine() = SettlementDetailLine(
        uid = uid,
        sourceType = sourceType.value,
        sourceId = sourceId,
        paymentId = paymentId,
        amount = amount,
        currency = currency,
        paymentMethod = paymentMethod,
        occurredAt = occurredAt,
    )

    private fun SettlementRunResponse.toAck() = SettlementRunAck(
        uid = uid,
        targetDate = targetDate,
        status = status,
        startedAt = startedAt,
    )
}
