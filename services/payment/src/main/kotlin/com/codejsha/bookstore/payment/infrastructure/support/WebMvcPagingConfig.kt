package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.platform.shared.data.PageSizeLimitCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.web.config.PageableHandlerMethodArgumentResolverCustomizer

@Configuration
class WebMvcPagingConfig {

    @Bean
    fun pageableCustomizer(): PageableHandlerMethodArgumentResolverCustomizer = PageSizeLimitCustomizer()
}
