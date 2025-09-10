package com.codejsha.bookstore.settlement.domain.service

import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementRepo
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailRepo
import com.codejsha.bookstore.settlement.application.port.support.TransactionRunner
import com.codejsha.bookstore.settlement.application.usecase.SettlementUseCase
import com.codejsha.bookstore.settlement.domain.aggregate.DailySettlement
import com.codejsha.bookstore.settlement.domain.aggregate.toAggregate
import com.codejsha.bookstore.settlement.domain.aggregate.toEntity
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service
import java.util.UUID

@Service
class SettlementService(
    private val dailySettlementRepo: DailySettlementRepo,
    private val settlementDetailRepo: SettlementDetailRepo,
    private val txRunner: TransactionRunner,
) : SettlementUseCase {

    @WithSpan
    override suspend fun findAllSettlements(
        option: SettlementQueryOption,
        pageable: Pageable,
        context: ActorContext,
    ): Page<DailySettlement> = txRunner.tx {
        dailySettlementRepo.findAll(option, pageable, context).map { it.toAggregate() }
    }

    @WithSpan
    override suspend fun findSettlement(uid: UUID, context: ActorContext): DailySettlement = txRunner.tx {
        val settlement = dailySettlementRepo.findOne(uid, context).toAggregate()
        val details = settlementDetailRepo.findAllByBucket(
            settlement.settlementDate,
            settlement.currency,
            settlement.paymentMethod,
            context,
        ).map { it.toEntity() }
        settlement.copy(details = details)
    }
}
