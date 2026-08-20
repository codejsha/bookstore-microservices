package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.Dashboard
import com.codejsha.platform.shared.data.ActorContext

interface DashboardUseCase {
    suspend fun loadDashboard(context: ActorContext): Dashboard
}
