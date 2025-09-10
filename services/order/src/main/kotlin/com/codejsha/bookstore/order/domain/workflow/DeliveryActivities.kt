package com.codejsha.bookstore.order.domain.workflow

import io.temporal.activity.ActivityInterface
import io.temporal.activity.ActivityMethod

@ActivityInterface
interface DeliveryActivities {
    @ActivityMethod
    fun createShipment(request: CreateShipmentRequest): CreateShipmentResult
}
