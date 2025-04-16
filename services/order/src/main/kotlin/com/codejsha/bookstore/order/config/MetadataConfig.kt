package com.codejsha.bookstore.order.config

import org.springframework.core.io.ClassPathResource
import org.springframework.stereotype.Component
import java.util.*

@Component
class MetadataConfig {
    val version: String by lazy { loadProjectProperties().getProperty("version") }

    private fun loadProjectProperties(): Properties {
        return Properties().apply {
            val resource = ClassPathResource("gradle.properties")
            resource.inputStream.use { load(it) }
        }
    }
}
