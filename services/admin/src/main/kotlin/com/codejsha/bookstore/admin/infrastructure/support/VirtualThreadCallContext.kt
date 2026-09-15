package com.codejsha.bookstore.admin.infrastructure.support

import com.codejsha.bookstore.admin.application.port.support.ConcurrentCallContext
import jakarta.annotation.PreDestroy
import kotlinx.coroutines.ExecutorCoroutineDispatcher
import kotlinx.coroutines.ThreadContextElement
import kotlinx.coroutines.asCoroutineDispatcher
import org.springframework.stereotype.Component
import org.springframework.web.context.request.RequestAttributes
import org.springframework.web.context.request.RequestContextHolder
import java.util.concurrent.Executors
import kotlin.coroutines.AbstractCoroutineContextElement
import kotlin.coroutines.CoroutineContext

@Component
class VirtualThreadCallContext : ConcurrentCallContext {

    private val dispatcher: ExecutorCoroutineDispatcher =
        Executors.newVirtualThreadPerTaskExecutor().asCoroutineDispatcher()

    override fun capture(): CoroutineContext =
        dispatcher + RequestAttributesElement(RequestContextHolder.getRequestAttributes())

    @PreDestroy
    fun close() {
        dispatcher.close()
    }

    private class RequestAttributesElement(
        private val attributes: RequestAttributes?,
    ) : AbstractCoroutineContextElement(Key), ThreadContextElement<RequestAttributes?> {

        companion object Key : CoroutineContext.Key<RequestAttributesElement>

        override fun updateThreadContext(context: CoroutineContext): RequestAttributes? {
            val previous = RequestContextHolder.getRequestAttributes()
            RequestContextHolder.setRequestAttributes(attributes)
            return previous
        }

        override fun restoreThreadContext(context: CoroutineContext, oldState: RequestAttributes?) {
            RequestContextHolder.setRequestAttributes(oldState)
        }
    }
}
