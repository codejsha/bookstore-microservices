package com.codejsha.bookstore.settlement.config

import org.springframework.boot.context.properties.ConfigurationPropertiesScan
import org.springframework.context.annotation.Configuration

@Configuration
@ConfigurationPropertiesScan(basePackages = ["com.codejsha.bookstore.settlement.config.properties"])
class PropertyScan
