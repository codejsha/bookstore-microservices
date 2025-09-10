package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.CatalogUseCase
import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.external.Subject
import com.codejsha.bookstore.admin.domain.model.external.Work
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkUpdateRequest
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class AdminCatalogControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `every catalog endpoint rejects a caller without the ADMIN role`() {
        bindPrincipal(roles = "VIEW")
        val useCase = mock(CatalogUseCase::class.java)
        val controller = AdminCatalogController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminCatalogListWorks(null, null, null, null) }
        assertFailsWith<ForbiddenException> { controller.adminCatalogReadWork(WORK_UID) }
        assertFailsWith<ForbiddenException> { controller.adminCatalogCreateWork(createRequest()) }
        assertFailsWith<ForbiddenException> {
            controller.adminCatalogUpdateWork(WORK_UID, AdminWorkUpdateRequest(title = "x"))
        }
        assertFailsWith<ForbiddenException> { controller.adminCatalogListAuthors(null, null) }
        assertFailsWith<ForbiddenException> { controller.adminCatalogListSubjects(null, null) }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `an absent pageable is forwarded as unpaged`(): Unit = runBlocking {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(CatalogUseCase::class.java)
        val controller = AdminCatalogController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        given(useCase.findAllWorks(WorkQueryOption(), unpaged, controllerContext))
            .willReturn(PageImpl(emptyList(), unpaged, 0))

        controller.adminCatalogListWorks(null, null, null, null)

        verify(useCase).findAllWorks(WorkQueryOption(), unpaged, controllerContext)
    }

    @Test
    fun `a sorted page is forwarded as the caller requested it`(): Unit = runBlocking {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(CatalogUseCase::class.java)
        val controller = AdminCatalogController(useCase, resolver)

        val pageable = PageRequest.of(2, 25, Sort.by(Sort.Direction.DESC, "title"))
        val option = WorkQueryOption(title = "dune")
        given(useCase.findAllWorks(option, pageable, controllerContext))
            .willReturn(PageImpl(listOf(work()), pageable, 51))

        val body = controller.adminCatalogListWorks("dune", null, null, pageable).body!!

        assertEquals(51L, body.total)
        assertEquals("Dune", body.items[0].title)
        verify(useCase).findAllWorks(option, pageable, controllerContext)
    }

    @Test
    fun `createWork answers 201 without a body`() {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(CatalogUseCase::class.java)
        val controller = AdminCatalogController(useCase, resolver)

        val response = controller.adminCatalogCreateWork(createRequest())

        assertEquals(201, response.statusCode.value())
        assertEquals(null, response.body)
    }

    @Test
    fun `updateWork forwards absent fields as null so the catalog leaves them alone`(): Unit = runBlocking {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(CatalogUseCase::class.java)
        val controller = AdminCatalogController(useCase, resolver)

        val expected = WorkUpdateCommand(
            title = "Dune Messiah",
            description = null,
            firstPublishDate = null,
            olKey = null,
            authorUids = null,
            subjectNames = null,
        )
        given(useCase.updateWork(WORK_UID, expected, controllerContext)).willReturn(work())

        controller.adminCatalogUpdateWork(WORK_UID, AdminWorkUpdateRequest(title = "Dune Messiah"))

        verify(useCase).updateWork(WORK_UID, expected, controllerContext)
    }

    private fun createRequest() = AdminWorkCreateRequest(title = "Dune", authorUids = listOf(AUTHOR_UID))

    private fun work() = Work(
        uid = WORK_UID,
        title = "Dune",
        description = null,
        firstPublishDate = "1965",
        olKey = null,
        authors = listOf(Author(uid = AUTHOR_UID, name = "Frank Herbert")),
        subjects = listOf(Subject(uid = SUBJECT_UID, name = "Science Fiction")),
        createdAt = OffsetDateTime.of(2026, 1, 1, 0, 0, 0, 0, ZoneOffset.UTC),
        updatedAt = null,
    )

    private companion object {
        const val WORK_UID = "22222222-2222-2222-2222-222222222222"
        const val AUTHOR_UID = "33333333-3333-3333-3333-333333333333"
        const val SUBJECT_UID = "44444444-4444-4444-4444-444444444444"
    }
}
