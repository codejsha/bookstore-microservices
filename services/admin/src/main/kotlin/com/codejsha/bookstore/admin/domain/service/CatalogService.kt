package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.application.usecase.CatalogUseCase
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class CatalogService(
    private val catalogClient: CatalogClient,
) : CatalogUseCase {

    override fun findAllWorks(option: WorkQueryOption, pageable: Pageable, context: ActorContext) =
        catalogClient.findAllWorks(option, pageable)

    override fun findWork(uid: String, context: ActorContext) = catalogClient.findWork(uid)

    override fun createWork(command: WorkCreateCommand, context: ActorContext) =
        catalogClient.createWork(command)

    override fun updateWork(uid: String, command: WorkUpdateCommand, context: ActorContext) =
        catalogClient.updateWork(uid, command)

    override fun findAllAuthors(name: String?, pageable: Pageable, context: ActorContext) =
        catalogClient.findAllAuthors(name, pageable)

    override fun findAllSubjects(name: String?, pageable: Pageable, context: ActorContext) =
        catalogClient.findAllSubjects(name, pageable)
}
