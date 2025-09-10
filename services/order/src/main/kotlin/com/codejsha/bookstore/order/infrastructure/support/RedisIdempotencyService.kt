package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.application.port.support.IdempotencyService
import org.springframework.data.redis.core.StringRedisTemplate
import org.springframework.stereotype.Component
import java.time.Duration

@Component
class RedisIdempotencyService(
    private val redisTemplate: StringRedisTemplate,
) : IdempotencyService {

    companion object {
        private val IDEMPOTENCY_TTL = Duration.ofHours(24)
    }

    override fun isProcessed(idempotencyKey: String): Boolean {
        return redisTemplate.hasKey("idempotency:$idempotencyKey") == true
    }

    override fun getResult(idempotencyKey: String): String? {
        return redisTemplate.opsForValue().get("idempotency:$idempotencyKey")
    }

    override fun tryMarkProcessed(idempotencyKey: String, resultUid: String): Boolean {
        return redisTemplate.opsForValue()
            .setIfAbsent("idempotency:$idempotencyKey", resultUid, IDEMPOTENCY_TTL) == true
    }

    override fun markProcessed(idempotencyKey: String, resultUid: String): Boolean =
        tryMarkProcessed(idempotencyKey, resultUid)
}
