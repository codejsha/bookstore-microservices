package com.codejsha.bookstore.payment.support

import com.codejsha.bookstore.payment.application.port.support.ActivityHeartbeat

class RecordingActivityHeartbeat : ActivityHeartbeat {

    val beats = mutableListOf<String>()

    override fun beat(detail: String) {
        beats += detail
    }
}

class CancellingActivityHeartbeat(
    private val failOn: String,
    private val error: RuntimeException,
) : ActivityHeartbeat {

    val beats = mutableListOf<String>()

    override fun beat(detail: String) {
        beats += detail
        if (detail.startsWith(failOn)) {
            throw error
        }
    }
}
