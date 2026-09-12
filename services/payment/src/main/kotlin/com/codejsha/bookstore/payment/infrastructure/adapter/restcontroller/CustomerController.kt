package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.CustomerApi
import com.codejsha.bookstore.generated.application.port.openapi.api.CustomerPaymentMethodApi
import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.payment.application.usecase.CustomerUseCase
import com.codejsha.bookstore.payment.domain.aggregate.CustomerAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentMethodEntity
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.payment.infrastructure.support.auth.Principal
import com.codejsha.bookstore.payment.infrastructure.support.auth.isStaff
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import org.springframework.web.server.ResponseStatusException
import java.net.URI
import java.time.ZoneOffset
import java.util.*

@RestController
class CustomerController(
    private val customerUseCase: CustomerUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : CustomerApi, CustomerPaymentMethodApi {

    // ─── CustomerApi ────────────────────────────────────────────────────────

    override fun customersGetAll(
        customerId: String?,
        email: String?,
        pageable: Pageable?
    ): ResponseEntity<CustomerFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val ownerFilter = if (principal.isStaff()) customerId else principal.sub
        val option = CustomerQueryOption(customerId = ownerFilter, email = email)
        val context = buildContext()

        val result = customerUseCase.findAllCustomers(option, pageable ?: Pageable.unpaged(), context)
        val response = CustomerFindAllResponse(
            total = result.totalElements,
            items = result.content.map { toCustomerFindResponse(it) }
        )
        ResponseEntity.ok(response)
    }

    override fun customersCreate(requestBody: CustomerCreateRequest): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        val ownerCustomerId = if (principal.isStaff()) requestBody.customerId else principal.sub
        val command = CustomerCreateCommand(
            customerId = ownerCustomerId,
            name = requestBody.name,
            email = requestBody.email,
            phone = requestBody.phone,
            phoneCountryCode = requestBody.phoneCountryCode,
            description = requestBody.description,
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        val customer = customerUseCase.createCustomer(command, context)
        ResponseEntity.created(URI.create("/api/v1/customers/${customer.uid}")).build()
    }

    override fun customersRead(uid: String): ResponseEntity<CustomerFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        val customer = customerUseCase.findCustomer(UUID.fromString(uid), context)
        assertOwnerOrStaff(principal, customer)
        ResponseEntity.ok(toCustomerFindResponse(customer))
    }

    override fun customersUpdate(
        uid: String,
        requestBody: CustomerUpdateRequest
    ): ResponseEntity<CustomerUpdateResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        assertOwnerOrStaff(principal, customerUseCase.findCustomer(UUID.fromString(uid), context))
        val command = CustomerUpdateCommand(
            name = requestBody.name,
            email = requestBody.email,
            phone = requestBody.phone,
            phoneCountryCode = requestBody.phoneCountryCode,
            description = requestBody.description,
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        val customer = customerUseCase.updateCustomer(UUID.fromString(uid), command, context)
        ResponseEntity.ok(toCustomerUpdateResponse(customer))
    }

    override fun customersDelete(uid: String): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        assertOwnerOrStaff(principal, customerUseCase.findCustomer(UUID.fromString(uid), context))
        customerUseCase.deleteCustomer(UUID.fromString(uid), context)
        ResponseEntity.noContent().build()
    }

    // ─── CustomerPaymentMethodApi ───────────────────────────────────────────

    override fun customerPaymentMethodsGetAll(
        customerUid: String,
        paymentMethod: PaymentMethodEnum?,
        pageable: Pageable?
    ): ResponseEntity<PaymentMethodFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val option = PaymentMethodQueryOption(paymentMethod = paymentMethod?.value)
        val context = buildContext()
        authorizeCustomer(customerUid, principal, context)

        val result = customerUseCase.findAllPaymentMethods(
            UUID.fromString(customerUid),
            option,
            pageable ?: Pageable.unpaged(),
            context
        )
        val response = PaymentMethodFindAllResponse(
            total = result.totalElements,
            items = result.content.map { toPaymentMethodFindResponse(it) }
        )
        ResponseEntity.ok(response)
    }

    override fun customerPaymentMethodsCreate(
        customerUid: String,
        requestBody: PaymentMethodCreateRequest
    ): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        authorizeCustomer(customerUid, principal, context)
        val command = PaymentMethodCreateCommand(
            paymentMethod = requestBody.paymentMethod.value,
            paymentMethodType = requestBody.paymentMethodType,
            paymentMethodIssuer = requestBody.paymentMethodIssuer,
            cardNetwork = requestBody.cardNetwork,
            cardLast4 = requestBody.cardLast4,
            cardExpMonth = requestBody.cardExpMonth,
            cardExpYear = requestBody.cardExpYear,
            cardHolderName = requestBody.cardHolderName,
            isDefault = requestBody.isDefault,
            metadata = null,
        )
        val pm = customerUseCase.createPaymentMethod(UUID.fromString(customerUid), command, context)
        ResponseEntity.created(URI.create("/api/v1/customers/$customerUid/payment-methods/${pm.uid}")).build()
    }

    override fun customerPaymentMethodsRead(
        customerUid: String,
        uid: String
    ): ResponseEntity<PaymentMethodFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        authorizeCustomer(customerUid, principal, context)
        val pm = customerUseCase.findPaymentMethod(UUID.fromString(customerUid), UUID.fromString(uid), context)
        ResponseEntity.ok(toPaymentMethodFindResponse(pm))
    }

    override fun customerPaymentMethodsUpdate(
        customerUid: String,
        uid: String,
        requestBody: PaymentMethodUpdateRequest
    ): ResponseEntity<PaymentMethodUpdateResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        authorizeCustomer(customerUid, principal, context)
        val command = PaymentMethodUpdateCommand(
            paymentMethod = requestBody.paymentMethod?.value,
            paymentMethodType = requestBody.paymentMethodType,
            paymentMethodIssuer = requestBody.paymentMethodIssuer,
            cardNetwork = requestBody.cardNetwork,
            cardLast4 = requestBody.cardLast4,
            cardExpMonth = requestBody.cardExpMonth,
            cardExpYear = requestBody.cardExpYear,
            cardHolderName = requestBody.cardHolderName,
            isDefault = requestBody.isDefault,
            metadata = null,
        )
        val pm = customerUseCase.updatePaymentMethod(
            UUID.fromString(customerUid),
            UUID.fromString(uid),
            command,
            context,
        )
        ResponseEntity.ok(toPaymentMethodUpdateResponse(pm))
    }

    override fun customerPaymentMethodsDelete(customerUid: String, uid: String): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        authorizeCustomer(customerUid, principal, context)
        customerUseCase.deletePaymentMethod(UUID.fromString(customerUid), UUID.fromString(uid), context)
        ResponseEntity.noContent().build()
    }

    // ─── Authorization ──────────────────────────────────────────────────────

    private suspend fun authorizeCustomer(customerUid: String, principal: Principal, context: ActorContext) {
        assertOwnerOrStaff(principal, customerUseCase.findCustomer(UUID.fromString(customerUid), context))
    }

    private fun assertOwnerOrStaff(principal: Principal, customer: CustomerAggregate) {
        if (principal.isStaff() || customer.customerId == principal.sub) return
        throw ResponseStatusException(HttpStatus.NOT_FOUND, "Customer not found")
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toCustomerFindResponse(agg: CustomerAggregate) = CustomerFindResponse(
        uid = agg.uid.toString(),
        customerId = agg.customerId,
        name = agg.name,
        email = agg.email,
        phone = agg.phone,
        phoneCountryCode = agg.phoneCountryCode,
        description = agg.description,
        createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun toCustomerUpdateResponse(agg: CustomerAggregate) = CustomerUpdateResponse(
        uid = agg.uid.toString(),
        customerId = agg.customerId,
        name = agg.name,
        email = agg.email,
        phone = agg.phone,
        phoneCountryCode = agg.phoneCountryCode,
        description = agg.description,
        createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun toPaymentMethodFindResponse(entity: PaymentMethodEntity) = PaymentMethodFindResponse(
        uid = entity.uid.toString(),
        paymentMethodId = entity.paymentMethodId,
        customerId = entity.customerId,
        paymentMethod = PaymentMethodEnum.fromValue(entity.paymentMethod.value),
        isDefault = entity.isDefault,
        paymentMethodType = entity.paymentMethodType,
        paymentMethodIssuer = entity.paymentMethodIssuer,
        cardNetwork = entity.cardNetwork,
        cardLast4 = entity.cardLast4,
        cardExpMonth = entity.cardExpMonth,
        cardExpYear = entity.cardExpYear,
        cardHolderName = entity.cardHolderName,
        createdAt = entity.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = entity.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun toPaymentMethodUpdateResponse(entity: PaymentMethodEntity) = PaymentMethodUpdateResponse(
        uid = entity.uid.toString(),
        paymentMethodId = entity.paymentMethodId,
        customerId = entity.customerId,
        paymentMethod = PaymentMethodEnum.fromValue(entity.paymentMethod.value),
        isDefault = entity.isDefault,
        paymentMethodType = entity.paymentMethodType,
        paymentMethodIssuer = entity.paymentMethodIssuer,
        cardNetwork = entity.cardNetwork,
        cardLast4 = entity.cardLast4,
        cardExpMonth = entity.cardExpMonth,
        cardExpYear = entity.cardExpYear,
        cardHolderName = entity.cardHolderName,
        createdAt = entity.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = entity.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun buildContext() = ActorContext(actorId = 0L, ActorType.USER)
}
