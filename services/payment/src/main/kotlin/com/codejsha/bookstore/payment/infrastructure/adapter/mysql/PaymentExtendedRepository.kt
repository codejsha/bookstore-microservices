package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.application.port.repo.PaymentExtendedRepo
import com.querydsl.jpa.impl.JPAQueryFactory
import org.springframework.stereotype.Repository

@Repository
class PaymentExtendedRepository(
    private val jpaQueryFactory: JPAQueryFactory
) : PaymentExtendedRepo
