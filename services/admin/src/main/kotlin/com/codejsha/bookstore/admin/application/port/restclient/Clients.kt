package com.codejsha.bookstore.admin.application.port.restclient

import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.Subject
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
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface CatalogClient {
    fun countWorks(): Long

    fun findAllWorks(option: WorkQueryOption, pageable: Pageable): Page<Work>

    fun findWork(uid: String): Work

    fun createWork(command: WorkCreateCommand)

    fun updateWork(uid: String, command: WorkUpdateCommand): Work

    fun findAllAuthors(name: String?, pageable: Pageable): Page<Author>

    fun findAllSubjects(name: String?, pageable: Pageable): Page<Subject>
}

interface IdentityClient {
    fun findAllUsers(option: UserQueryOption, pageable: Pageable): Page<User>

    fun listRisk(): List<RiskEntry>

    fun flagRisk(uid: String, level: String, reason: String, ttlSeconds: Long?): RiskEntry

    fun unflagRisk(uid: String)

    fun findUser(uid: String): User

    fun updateRoles(uid: String, roles: List<String>): User

    fun suspendUser(uid: String): User

    fun reactivateUser(uid: String): User

    fun deactivateUser(uid: String): User
}

interface OrderClient {
    fun countOrders(): Long

    fun findAllOrders(option: OrderQueryOption, pageable: Pageable): Page<Order>

    fun findOrder(uid: String): Order

    fun cancelOrder(uid: String): Order
}

interface InventoryClient {
    fun countWarehouses(): Long

    fun findAllWarehouses(name: String?, pageable: Pageable): Page<Warehouse>

    fun findWarehouse(uid: String): Warehouse

    fun findAllStocks(option: StockQueryOption, pageable: Pageable): Page<Stock>
}

interface PaymentClient {
    fun findAllPayments(option: PaymentQueryOption, pageable: Pageable): Page<Payment>

    fun findPayment(uid: String): Payment

    fun findAllRefunds(option: RefundQueryOption, pageable: Pageable): Page<Refund>

    fun findRefund(uid: String): Refund
}

interface SettlementClient {
    fun findAllSettlements(option: SettlementQueryOption, pageable: Pageable): Page<SettlementBucket>

    fun findSettlement(uid: String): SettlementBucket

    fun triggerSettlementRun(command: TriggerSettlementRunCommand): SettlementRunAck
}
