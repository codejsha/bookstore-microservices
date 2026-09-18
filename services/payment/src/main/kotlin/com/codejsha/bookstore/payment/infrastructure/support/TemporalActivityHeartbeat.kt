package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.application.port.support.ActivityHeartbeat
import io.temporal.activity.Activity
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component

@Component
class TemporalActivityHeartbeat : ActivityHeartbeat {

    override fun beat(detail: String) {
        val context = try {
            Activity.getExecutionContext()
        } catch (e: IllegalStateException) {
            log.debug("Skipping heartbeat for {} outside of an activity context: {}", detail, e.message)
            return
        }
        context.heartbeat(detail)
    }

    private companion object {
        val log = LoggerFactory.getLogger(TemporalActivityHeartbeat::class.java)
    }
}
