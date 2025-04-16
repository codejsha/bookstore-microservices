package com.codejsha.bookstore.payment.application.port.repo

import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import org.springframework.data.jpa.repository.JpaRepository

interface PaymentRepo : JpaRepository<PaymentEntity, Long>, PaymentExtendedRepo

interface PaymentExtendedRepo
