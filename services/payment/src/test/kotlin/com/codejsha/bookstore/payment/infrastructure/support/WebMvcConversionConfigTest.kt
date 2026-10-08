package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.generated.application.port.openapi.model.MandateStatus
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentMethodEnum
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentStatus
import com.codejsha.bookstore.generated.application.port.openapi.model.RefundStatus
import org.junit.jupiter.api.Test
import org.springframework.core.convert.ConversionFailedException
import org.springframework.format.support.DefaultFormattingConversionService
import tools.jackson.databind.json.JsonMapper
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class WebMvcConversionConfigTest {

    private val conversionService = DefaultFormattingConversionService().also {
        WebMvcConversionConfig(JsonMapper.builder().build()).addFormatters(it)
    }

    @Test
    fun convert_specWireValue_resolvesConstant() {
        assertEquals(
            PaymentStatus.REQUIRES_PAYMENT_METHOD,
            conversionService.convert("requires_payment_method", PaymentStatus::class.java),
        )
        assertEquals(RefundStatus.REVIEW, conversionService.convert("review", RefundStatus::class.java))
        assertEquals(
            PaymentMethodEnum.BANK_TRANSFER,
            conversionService.convert("bank_transfer", PaymentMethodEnum::class.java),
        )
        assertEquals(MandateStatus.REVOKED, conversionService.convert("revoked", MandateStatus::class.java))
    }

    @Test
    fun convert_constantName_resolvesConstant() {
        assertEquals(PaymentStatus.SUCCEEDED, conversionService.convert("SUCCEEDED", PaymentStatus::class.java))
    }

    @Test
    fun convert_unknownValue_fails() {
        assertFailsWith<ConversionFailedException> {
            conversionService.convert("bogus", PaymentStatus::class.java)
        }
    }
}
