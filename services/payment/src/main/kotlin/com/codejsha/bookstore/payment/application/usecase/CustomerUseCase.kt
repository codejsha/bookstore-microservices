package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.CustomerAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentMethodEntity
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface CustomerUseCase {
    // ─── Customer (Aggregate Root) ──────────────────────────────────────────
    fun findAllCustomers(option: CustomerQueryOption, pageable: Pageable, context: ActorContext): Page<CustomerAggregate>
    fun findCustomer(uid: UUID, context: ActorContext): CustomerAggregate
    fun createCustomer(command: CustomerCreateCommand, context: ActorContext): CustomerAggregate
    fun updateCustomer(uid: UUID, command: CustomerUpdateCommand, context: ActorContext): CustomerAggregate
    fun deleteCustomer(uid: UUID, context: ActorContext)

    // ─── PaymentMethod (Internal Entity) ────────────────────────────────────
    fun findAllPaymentMethods(customerUid: UUID, option: PaymentMethodQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentMethodEntity>
    fun findPaymentMethod(customerUid: UUID, uid: UUID, context: ActorContext): PaymentMethodEntity
    fun createPaymentMethod(customerUid: UUID, command: PaymentMethodCreateCommand, context: ActorContext): PaymentMethodEntity
    fun updatePaymentMethod(customerUid: UUID, uid: UUID, command: PaymentMethodUpdateCommand, context: ActorContext): PaymentMethodEntity
    fun deletePaymentMethod(customerUid: UUID, uid: UUID, context: ActorContext)
}
