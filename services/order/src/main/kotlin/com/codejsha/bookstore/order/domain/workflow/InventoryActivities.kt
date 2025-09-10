package com.codejsha.bookstore.order.domain.workflow

import io.temporal.activity.ActivityInterface
import io.temporal.activity.ActivityMethod

@ActivityInterface
interface InventoryActivities {
    @ActivityMethod
    fun reserveStock(request: ReserveStockRequest)

    @ActivityMethod
    fun releaseStock(request: ReserveStockRequest)
}
