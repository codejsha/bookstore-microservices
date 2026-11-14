package com.codejsha.bookstore.order.infrastructure.support

import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoMoreInteractions
import org.springframework.data.redis.core.StringRedisTemplate
import org.springframework.data.redis.core.ValueOperations
import java.time.Duration
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNull
import kotlin.test.assertTrue

class RedisIdempotencyServiceTest {

    private val template = mock(StringRedisTemplate::class.java)

    @Suppress("UNCHECKED_CAST")
    private val values = mock(ValueOperations::class.java) as ValueOperations<String, String>

    private val service = RedisIdempotencyService(template)

    @Test
    fun `isProcessed_keyStored_reportsProcessed`() {
        given(template.hasKey("idempotency:checkout-1")).willReturn(true)

        assertTrue(service.isProcessed("checkout-1"))
    }

    @Test
    fun `isProcessed_keyAbsent_reportsUnprocessed`() {
        given(template.hasKey("idempotency:checkout-1")).willReturn(false)

        assertFalse(service.isProcessed("checkout-1"))
    }

    @Test
    fun `isProcessed_redisAnswersNull_reportsUnprocessed`() {
        given(template.hasKey("idempotency:checkout-1")).willReturn(null)

        assertFalse(service.isProcessed("checkout-1"))
    }

    @Test
    fun `getResult_keyStored_yieldsStoredResultUid`() {
        given(template.opsForValue()).willReturn(values)
        given(values.get("idempotency:checkout-1")).willReturn("order-uid-1")

        assertEquals("order-uid-1", service.getResult("checkout-1"))
    }

    @Test
    fun `getResult_keyAbsent_yieldsNull`() {
        given(template.opsForValue()).willReturn(values)
        given(values.get("idempotency:checkout-1")).willReturn(null)

        assertNull(service.getResult("checkout-1"))
    }

    @Test
    fun `tryMarkProcessed_firstWrite_storesResultFor24Hours`() {
        given(template.opsForValue()).willReturn(values)
        given(values.setIfAbsent("idempotency:checkout-1", "order-uid-1", Duration.ofHours(24))).willReturn(true)

        assertTrue(service.tryMarkProcessed("checkout-1", "order-uid-1"))

        verify(values).setIfAbsent("idempotency:checkout-1", "order-uid-1", Duration.ofHours(24))
        verifyNoMoreInteractions(values)
    }

    @Test
    fun `tryMarkProcessed_keyAlreadyStored_reportsLostRace`() {
        given(template.opsForValue()).willReturn(values)
        given(values.setIfAbsent("idempotency:checkout-1", "order-uid-2", Duration.ofHours(24))).willReturn(false)

        assertFalse(service.tryMarkProcessed("checkout-1", "order-uid-2"))
    }

    @Test
    fun `tryMarkProcessed_redisAnswersNull_reportsNotStored`() {
        given(template.opsForValue()).willReturn(values)
        given(values.setIfAbsent("idempotency:checkout-1", "order-uid-1", Duration.ofHours(24))).willReturn(null)

        assertFalse(service.tryMarkProcessed("checkout-1", "order-uid-1"))
    }

    @Test
    fun `markProcessed_firstWriteThenRepeat_storesOnceAndRejectsRepeat`() {
        given(template.opsForValue()).willReturn(values)
        given(values.setIfAbsent("idempotency:checkout-1", "order-uid-1", Duration.ofHours(24)))
            .willReturn(true, false)

        assertTrue(service.markProcessed("checkout-1", "order-uid-1"))
        assertFalse(service.markProcessed("checkout-1", "order-uid-1"))
    }
}
