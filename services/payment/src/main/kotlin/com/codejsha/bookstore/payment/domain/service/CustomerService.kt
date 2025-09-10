package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.repo.CustomerRepo
import com.codejsha.bookstore.payment.application.port.repo.CustomerResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentMethodRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentMethodResult
import com.codejsha.bookstore.payment.application.port.support.TransactionRunner
import com.codejsha.bookstore.payment.application.usecase.CustomerUseCase
import com.codejsha.bookstore.payment.domain.aggregate.CustomerAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentMethodEntity
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service
import java.util.UUID

@Service
class CustomerService(
    private val customerRepo: CustomerRepo,
    private val paymentMethodRepo: PaymentMethodRepo,
    private val txRunner: TransactionRunner,
) : CustomerUseCase {

    @WithSpan
    override suspend fun findAllCustomers(
        option: CustomerQueryOption, pageable: Pageable, context: ActorContext
    ): Page<CustomerAggregate> = txRunner.tx {
        customerRepo.findAll(option, pageable, context).map { toCustomerAggregate(it) }
    }

    @WithSpan
    override suspend fun findCustomer(uid: UUID, context: ActorContext): CustomerAggregate = txRunner.tx {
        toCustomerAggregate(customerRepo.findOne(uid, context))
    }

    @WithSpan
    override suspend fun createCustomer(
        command: CustomerCreateCommand, context: ActorContext
    ): CustomerAggregate = txRunner.tx {
        toCustomerAggregate(customerRepo.create(command, context))
    }

    @WithSpan
    override suspend fun updateCustomer(
        uid: UUID, command: CustomerUpdateCommand, context: ActorContext
    ): CustomerAggregate = txRunner.tx {
        toCustomerAggregate(customerRepo.update(uid, command, context))
    }

    @WithSpan
    override suspend fun deleteCustomer(uid: UUID, context: ActorContext) {
        txRunner.tx {
            customerRepo.delete(uid, context)
        }
    }

    @WithSpan
    override suspend fun findAllPaymentMethods(
        customerUid: UUID, option: PaymentMethodQueryOption, pageable: Pageable, context: ActorContext
    ): Page<PaymentMethodEntity> = txRunner.tx {
        paymentMethodRepo.findAllByCustomer(customerUid, option, pageable, context).map { toPaymentMethodEntity(it) }
    }

    @WithSpan
    override suspend fun findPaymentMethod(
        customerUid: UUID, uid: UUID, context: ActorContext
    ): PaymentMethodEntity = txRunner.tx {
        toPaymentMethodEntity(paymentMethodRepo.findOne(customerUid, uid, context))
    }

    @WithSpan
    override suspend fun createPaymentMethod(
        customerUid: UUID, command: PaymentMethodCreateCommand, context: ActorContext
    ): PaymentMethodEntity = txRunner.tx {
        toPaymentMethodEntity(paymentMethodRepo.create(customerUid, command, context))
    }

    @WithSpan
    override suspend fun updatePaymentMethod(
        customerUid: UUID, uid: UUID, command: PaymentMethodUpdateCommand, context: ActorContext
    ): PaymentMethodEntity = txRunner.tx {
        toPaymentMethodEntity(paymentMethodRepo.update(customerUid, uid, command, context))
    }

    @WithSpan
    override suspend fun deletePaymentMethod(customerUid: UUID, uid: UUID, context: ActorContext) {
        txRunner.tx {
            paymentMethodRepo.delete(customerUid, uid, context)
        }
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toCustomerAggregate(result: CustomerResult) = CustomerAggregate(
        id = result.id,
        uid = result.uid,
        customerId = result.customerId,
        name = result.name,
        email = result.email,
        phone = result.phone,
        phoneCountryCode = result.phoneCountryCode,
        description = result.description,
        metadata = result.metadata,
        defaultBillingAddress = null,
        defaultShippingAddress = null,
        paymentMethods = emptyList(),
        createdAt = result.createdAt,
        updatedAt = result.updatedAt,
    )

    private fun toPaymentMethodEntity(result: PaymentMethodResult) = PaymentMethodEntity(
        id = result.id,
        uid = result.uid,
        paymentMethodId = result.paymentMethodId,
        customerId = result.customerId,
        paymentMethod = PaymentMethodType.fromValue(result.paymentMethod),
        paymentMethodType = result.paymentMethodType,
        paymentMethodIssuer = result.paymentMethodIssuer,
        cardNetwork = result.cardNetwork,
        cardLast4 = result.cardLast4,
        cardExpMonth = result.cardExpMonth,
        cardExpYear = result.cardExpYear,
        cardHolderName = result.cardHolderName,
        isDefault = result.isDefault,
        metadata = result.metadata,
        createdAt = result.createdAt,
        updatedAt = result.updatedAt,
    )
}
