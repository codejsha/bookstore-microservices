package com.codejsha.bookstore.order.config

import org.springframework.boot.context.properties.ConfigurationPropertiesScan
import org.springframework.context.annotation.Configuration

@Configuration
@ConfigurationPropertiesScan(basePackages = ["com.codejsha.bookstore.order.config.properties"])
class PropertyScan
