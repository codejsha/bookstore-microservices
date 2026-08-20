package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface RefundUseCase {
    suspend fun findAllRefunds(option: RefundQueryOption, pageable: Pageable, context: ActorContext): Page<RefundAggregate>
    suspend fun findRefund(uid: UUID, context: ActorContext): RefundAggregate
    suspend fun createRefund(command: RefundCreateCommand, context: ActorContext): RefundAggregate
}
