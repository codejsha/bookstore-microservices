package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.repo.MandateRepo
import com.codejsha.bookstore.payment.application.port.repo.MandateResult
import com.codejsha.bookstore.payment.application.port.support.TransactionRunner
import com.codejsha.bookstore.payment.application.usecase.MandateUseCase
import com.codejsha.bookstore.payment.domain.aggregate.MandateAggregate
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.constant.MandateType
import com.codejsha.bookstore.payment.domain.model.command.MandateCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.MandateSetupCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service
import java.util.UUID

@Service
class MandateService(
    private val mandateRepo: MandateRepo,
    private val hyperswitchClient: HyperswitchClient,
    private val txRunner: TransactionRunner,
) : MandateUseCase {

    private companion object {
        const val ACCEPTANCE_TYPE_ONLINE = "online"
    }

    @WithSpan
    override suspend fun findAllMandates(
        option: MandateQueryOption, pageable: Pageable, context: ActorContext
    ): Page<MandateAggregate> = txRunner.tx {
        mandateRepo.findAll(option, pageable, context).map { toMandateAggregate(it) }
    }

    @WithSpan
    override suspend fun findMandate(uid: UUID, context: ActorContext): MandateAggregate = txRunner.tx {
        toMandateAggregate(mandateRepo.findOne(uid, context))
    }

    @WithSpan
    override suspend fun setupMandate(command: MandateSetupCommand, context: ActorContext): MandateAggregate {
        val hyperswitchResult = hyperswitchClient.setupMandate(
            HyperswitchSetupMandateCommand(
                idempotencyKey = "mandate:${command.customerId}:${command.paymentMethodToken}",
                customerId = command.customerId,
                paymentMethodToken = command.paymentMethodToken,
                currency = command.currency,
                mandateAmountMinor = command.mandateAmountMinor,
            ),
        )

        return txRunner.tx {
            toMandateAggregate(
                mandateRepo.create(
                    MandateCreateCommand(
                        mandateId = hyperswitchResult.gatewayMandateId,
                        customerId = command.customerId,
                        paymentMethodId = hyperswitchResult.gatewayPaymentMethodId,
                        mandateType = MandateType.MULTI_USE.value,
                        mandateStatus = hyperswitchResult.status,
                        mandateAmount = command.mandateAmountMinor,
                        mandateCurrency = command.currency,
                        setupFutureUsage = FutureUsage.OFF_SESSION.value,
                        customerAcceptanceType = ACCEPTANCE_TYPE_ONLINE,
                    ),
                    context,
                ),
            )
        }
    }

    @WithSpan
    override suspend fun revokeMandate(uid: UUID, context: ActorContext): MandateAggregate = txRunner.tx {
        toMandateAggregate(mandateRepo.revoke(uid, context))
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toMandateAggregate(result: MandateResult) = MandateAggregate(
        id = result.id,
        uid = result.uid,
        mandateId = result.mandateId,
        customerId = result.customerId,
        paymentMethodId = result.paymentMethodId,
        mandateType = MandateType.fromValue(result.mandateType),
        mandateStatus = MandateStatus.fromValue(result.mandateStatus),
        mandateAmount = result.mandateAmount,
        mandateCurrency = result.mandateCurrency,
        startDate = result.startDate,
        endDate = result.endDate,
        setupFutureUsage = result.setupFutureUsage?.let { FutureUsage.fromValue(it) },
        customerAcceptanceType = result.customerAcceptanceType,
        customerAcceptedAt = result.customerAcceptedAt,
        metadata = result.metadata,
        createdAt = result.createdAt,
        updatedAt = result.updatedAt,
    )
}
