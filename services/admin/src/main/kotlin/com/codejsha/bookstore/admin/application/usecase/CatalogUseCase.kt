package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.external.Subject
import com.codejsha.bookstore.admin.domain.model.external.Work
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface CatalogUseCase {
    suspend fun findAllWorks(option: WorkQueryOption, pageable: Pageable, context: ActorContext): Page<Work>

    suspend fun findWork(uid: String, context: ActorContext): Work

    suspend fun createWork(command: WorkCreateCommand, context: ActorContext)

    suspend fun updateWork(uid: String, command: WorkUpdateCommand, context: ActorContext): Work

    suspend fun findAllAuthors(name: String?, pageable: Pageable, context: ActorContext): Page<Author>

    suspend fun findAllSubjects(name: String?, pageable: Pageable, context: ActorContext): Page<Subject>
}
