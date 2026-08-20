import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin/users")
@tag("Admin User")
@useAuth(BearerAuth)
interface AdminUsers {
  @doc("List user accounts")
  @get
  listUsers(
    @query email?: string,
    @query name?: string,
    @query phone?: string,

    @doc("ACTIVE, SUSPENDED, or DEACTIVATED")
    @query
    status?: string,

    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminUserFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Find a user by uid")
  @route("{uid}")
  @get
  readUser(@path uid: uuid): AdminUserResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Replace a user's roles. An administrator may not change their own roles.")
  @route("{uid}/roles")
  @put
  updateUserRoles(
    @path uid: uuid,
    @body body: AdminUserRolesRequest,
  ): AdminUserResponse | BadRequestError | UnauthorizedError | ForbiddenError | NotFoundError;

  @doc("Suspend an account (disabled in Keycloak). An administrator may not suspend themselves.")
  @route("{uid}/suspend")
  @post
  suspendUser(
    @path uid: uuid,
  ): AdminUserResponse | BadRequestError | UnauthorizedError | ForbiddenError | NotFoundError;

  @doc("Reactivate a suspended account")
  @route("{uid}/reactivate")
  @post
  reactivateUser(
    @path uid: uuid,
  ): AdminUserResponse | BadRequestError | UnauthorizedError | ForbiddenError | NotFoundError;

  @doc("Deactivate an account (deleted in Keycloak, soft-deleted locally). An administrator may not deactivate themselves.")
  @route("{uid}/deactivate")
  @post
  deactivateUser(
    @path uid: uuid,
  ): AdminUserResponse | BadRequestError | UnauthorizedError | ForbiddenError | NotFoundError;
}
