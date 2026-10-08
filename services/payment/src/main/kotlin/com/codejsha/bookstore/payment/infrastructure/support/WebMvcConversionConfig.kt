package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.platform.shared.web.WireValueEnumConverterFactory
import org.springframework.context.annotation.Configuration
import org.springframework.format.FormatterRegistry
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer
import tools.jackson.databind.ObjectMapper

@Configuration
class WebMvcConversionConfig(
    private val objectMapper: ObjectMapper,
) : WebMvcConfigurer {

    override fun addFormatters(registry: FormatterRegistry) {
        registry.removeConvertible(String::class.java, Enum::class.java)
        registry.addConverterFactory(WireValueEnumConverterFactory(objectMapper))
    }
}
