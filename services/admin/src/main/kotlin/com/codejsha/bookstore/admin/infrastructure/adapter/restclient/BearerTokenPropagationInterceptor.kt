package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.infrastructure.support.auth.RequestTokenResolver
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpRequest
import org.springframework.http.client.ClientHttpRequestExecution
import org.springframework.http.client.ClientHttpRequestInterceptor
import org.springframework.http.client.ClientHttpResponse
import org.springframework.stereotype.Component

@Component
class BearerTokenPropagationInterceptor(
    private val tokenResolver: RequestTokenResolver,
) : ClientHttpRequestInterceptor {

    override fun intercept(
        request: HttpRequest,
        body: ByteArray,
        execution: ClientHttpRequestExecution,
    ): ClientHttpResponse {
        tokenResolver.currentAuthorization()?.let { token ->
            request.headers.set(HttpHeaders.AUTHORIZATION, token)
        }
        return execution.execute(request, body)
    }
}
