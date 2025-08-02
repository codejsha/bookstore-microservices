package com.codejsha.bookstore.payment.config

import org.springframework.context.annotation.ComponentScan
import org.springframework.context.annotation.Configuration

@Configuration
@ComponentScan(basePackages = ["com.netflix.conductor"])
class ConductorConfig
