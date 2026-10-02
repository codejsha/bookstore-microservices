package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.external.Subject
import com.codejsha.bookstore.admin.domain.model.external.Work
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.AuthorApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.SubjectApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.WorkApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.AuthorItem
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.SubjectItem
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkCreateRequest
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkItem
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkUpdateRequest
import com.codejsha.bookstore.generated.application.port.restclient.catalog.model.WorkUpdateResponse
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

@Component
class CatalogRestClient(
    private val workApi: WorkApi,
    private val authorApi: AuthorApi,
    private val subjectApi: SubjectApi,
) : CatalogClient {

    override fun countWorks(): Long =
        workApi.worksSearch(null, null, null, null, COUNT_ONLY_PAGE_SIZE, null, null).total

    override fun findAllWorks(option: WorkQueryOption, pageable: Pageable): Page<Work> {
        val response = workApi.worksSearch(
            option.title,
            option.authorUid,
            option.subjectUid,
            null,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toWork() }, pageable, response.total)
    }

    override fun findWork(uid: String): Work = workApi.worksRead(uid).toWork()

    override fun createWork(command: WorkCreateCommand) {
        workApi.worksCreate(
            WorkCreateRequest(
                title = command.title,
                description = command.description,
                firstPublishDate = command.firstPublishDate,
                olKey = command.olKey,
                authorUids = command.authorUids,
                subjectNames = command.subjectNames,
                coverUids = null,
            )
        )
    }

    override fun updateWork(uid: String, command: WorkUpdateCommand): Work =
        workApi.worksUpdate(
            uid,
            WorkUpdateRequest(
                title = command.title,
                description = command.description,
                firstPublishDate = command.firstPublishDate,
                olKey = command.olKey,
                authorUids = command.authorUids,
                subjectNames = command.subjectNames,
                coverUids = null,
            ),
        ).toWork()

    override fun findAllAuthors(name: String?, pageable: Pageable): Page<Author> {
        val response = authorApi.authorsGetAll(
            name,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toAuthor() }, pageable, response.total)
    }

    override fun findAllSubjects(name: String?, pageable: Pageable): Page<Subject> {
        val response = subjectApi.subjectsGetAll(
            name,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toSubject() }, pageable, response.total)
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun AuthorItem.toAuthor() = Author(uid = uid, name = name)

    private fun SubjectItem.toSubject() = Subject(uid = uid, name = name)

    private fun WorkItem.toWork() = Work(
        uid = uid,
        title = title,
        description = description,
        firstPublishDate = firstPublishDate,
        olKey = olKey,
        authors = authors.map { it.toAuthor() },
        subjects = subjects.map { it.toSubject() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun WorkFindResponse.toWork() = Work(
        uid = uid,
        title = title,
        description = description,
        firstPublishDate = firstPublishDate,
        olKey = olKey,
        authors = authors.map { it.toAuthor() },
        subjects = subjects.map { it.toSubject() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun WorkUpdateResponse.toWork() = Work(
        uid = uid,
        title = title,
        description = description,
        firstPublishDate = firstPublishDate,
        olKey = olKey,
        authors = authors.map { it.toAuthor() },
        subjects = subjects.map { it.toSubject() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )
}
