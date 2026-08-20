package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.MandateAggregate
import com.codejsha.bookstore.payment.domain.model.command.MandateSetupCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface MandateUseCase {
    suspend fun findAllMandates(option: MandateQueryOption, pageable: Pageable, context: ActorContext): Page<MandateAggregate>
    suspend fun findMandate(uid: UUID, context: ActorContext): MandateAggregate

    suspend fun setupMandate(command: MandateSetupCommand, context: ActorContext): MandateAggregate

    suspend fun revokeMandate(uid: UUID, context: ActorContext): MandateAggregate
}
