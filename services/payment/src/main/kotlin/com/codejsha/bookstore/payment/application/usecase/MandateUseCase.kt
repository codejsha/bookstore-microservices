package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.MandateAggregate
import com.codejsha.bookstore.payment.domain.model.command.MandateSetupCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface MandateUseCase {
    fun findAllMandates(option: MandateQueryOption, pageable: Pageable, context: ActorContext): Page<MandateAggregate>
    fun findMandate(uid: UUID, context: ActorContext): MandateAggregate

    fun setupMandate(command: MandateSetupCommand, context: ActorContext): MandateAggregate

    fun revokeMandate(uid: UUID, context: ActorContext): MandateAggregate
}
