import "@typespec/http";

using Http;

namespace IdentityService;

@route("/api/v1/users")
@tag("User")
@useAuth(BearerAuth)
interface Users {
  @doc("Register a new user (creates in Keycloak + local DB)")
  @post
  register(@body body: UserRegisterRequest): {
    @statusCode statusCode: 201;
    @body body: UserFindResponse;
  } | BadRequestError;

  @doc("Find all users")
  @get
  getAll(
    @query email?: string,
    @query name?: string,
    @query phone?: string,
    @query status?: UserStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): UserFindAllResponse;

  @doc("Find user by uid")
  @route("{uid}")
  @get
  read(@path uid: uuid): UserFindResponse | NotFoundError;

  @doc("Find user by email")
  @route("by-email/{email}")
  @get
  readByEmail(@path email: string): UserFindResponse | NotFoundError;

  @doc("Update user profile")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: UserUpdateRequest): UserFindResponse | NotFoundError;

  @doc("Update user roles")
  @route("{uid}/roles")
  @put
  updateRoles(@path uid: uuid, @body body: UserRolesRequest): UserFindResponse | NotFoundError;

  @doc("Suspend user account (disable in Keycloak)")
  @route("{uid}/suspend")
  @post
  suspend(@path uid: uuid): UserFindResponse | BadRequestError | NotFoundError;

  @doc("Reactivate suspended user account (enable in Keycloak)")
  @route("{uid}/reactivate")
  @post
  reactivate(@path uid: uuid): UserFindResponse | BadRequestError | NotFoundError;

  @doc("Deactivate user account (delete in Keycloak, soft-delete locally)")
  @route("{uid}/deactivate")
  @post
  deactivate(@path uid: uuid): UserFindResponse | BadRequestError | NotFoundError;

  @doc("Sync user data from Keycloak to local DB")
  @route("sync/{idp_uid}")
  @post
  syncFromIdp(@path idp_uid: uuid): UserFindResponse | NotFoundError;
}
