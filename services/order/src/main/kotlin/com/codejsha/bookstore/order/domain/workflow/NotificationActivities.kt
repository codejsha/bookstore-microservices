package com.codejsha.bookstore.order.domain.workflow

import io.temporal.activity.ActivityInterface
import io.temporal.activity.ActivityMethod

@ActivityInterface
interface NotificationActivities {
    @ActivityMethod
    fun sendShipmentUpdate(request: ShipmentNotificationRequest)
}
