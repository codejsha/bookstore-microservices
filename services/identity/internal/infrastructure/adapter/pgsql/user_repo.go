package pgsql

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/database/gormutils"

	"github.com/codejsha/bookstore-microservices/identity/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/identity/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/identity/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

var _ repo.UserRepo = (*userRepository)(nil)

type userRepository struct {
	db *gorm.DB
	genrepo.UsersRepo
}

func NewUserRepository(dataSource *database.DataSource) repo.UserRepo {
	db := dataSource.DB()
	return &userRepository{
		db:        db,
		UsersRepo: genrepo.NewUsersRepo(db),
	}
}

func (r userRepository) FindAll(ctx context.Context, opt option.UserQueryOption) (int64, []*repo.UserResult, error) {
	q := dao.Use(r.db)
	sortable := map[string]field.OrderExpr{
		"uid":         q.UsersEntity.Uid,
		"email":       q.UsersEntity.Email,
		"firstname":   q.UsersEntity.FirstName,
		"lastname":    q.UsersEntity.LastName,
		"phone":       q.UsersEntity.Phone,
		"role":        q.UsersEntity.Role,
		"status":      q.UsersEntity.Status,
		"lastloginat": q.UsersEntity.LastLoginAt,
		"createdat":   q.UsersEntity.CreatedAt,
		"updatedat":   q.UsersEntity.UpdatedAt,
	}

	base := q.UsersEntity.WithContext(ctx).Scopes(r.buildWhereScope(q, opt))

	total, err := base.Count()
	if err != nil {
		return 0, nil, err
	}

	entities, err := base.
		Scopes(gormutils.BuildPageScope(opt.Page(), sortable)).
		Find()
	if err != nil {
		return 0, nil, err
	}

	results := make([]*repo.UserResult, len(entities))
	for i, e := range entities {
		results[i] = toUserResult(e)
	}
	return total, results, nil
}

func (r userRepository) buildWhereScope(q *dao.Query, opt option.UserQueryOption) func(gen.Dao) gen.Dao {
	conds := make([]gen.Condition, 0)
	if v := opt.Email(); v != nil && *v != "" {
		conds = append(conds, q.UsersEntity.Email.Eq(*v))
	}
	if v := opt.Name(); v != nil && *v != "" {
		conds = append(conds, gen.Cond(gorm.Expr("(users.first_name ILIKE ? OR users.last_name ILIKE ?)", "%"+*v+"%", "%"+*v+"%"))...)
	}
	if v := opt.Phone(); v != nil && *v != "" {
		conds = append(conds, q.UsersEntity.Phone.Eq(*v))
	}

	return func(d gen.Dao) gen.Dao {
		if len(conds) > 0 {
			return d.Where(conds...)
		}
		return d
	}
}

func (r userRepository) FindByUid(ctx context.Context, uid string) (*repo.UserResult, error) {
	q := dao.Use(r.db)
	e, err := q.UsersEntity.WithContext(ctx).
		Where(q.UsersEntity.IdpUid.Eq(uid)).
		First()
	if err != nil {
		return nil, err
	}
	return toUserResult(e), nil
}

func (r userRepository) FetchByEmail(ctx context.Context, email string) ([]*repo.UserResult, error) {
	q := dao.Use(r.db)
	entities, err := q.UsersEntity.WithContext(ctx).
		Where(q.UsersEntity.Email.Eq(email)).
		Find()
	if err != nil {
		return nil, err
	}
	results := make([]*repo.UserResult, len(entities))
	for i, e := range entities {
		results[i] = toUserResult(e)
	}
	return results, nil
}

func (r userRepository) Upsert(ctx context.Context, profile repo.UserUpsert) (*repo.UserResult, error) {
	q := dao.Use(r.db)
	existing, err := q.UsersEntity.WithContext(ctx).
		Where(q.UsersEntity.IdpUid.Eq(profile.IdpUid)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	rolesJSON := encodeRoles(profile.Roles)
	now := time.Now()

	if existing == nil {
		ent := &entity.UsersEntity{
			Uid:       uuid.Must(uuid.NewV7()).String(),
			IdpUid:    profile.IdpUid,
			Email:     profile.Email,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Phone:     profile.Phone,
			Role:      rolesJSON,
			Status:    profile.Status,
			CreatedAt: now,
			Version:   1,
		}
		if err := q.UsersEntity.WithContext(ctx).Create(ent); err != nil {
			return nil, err
		}
		return toUserResult(ent), nil
	}

	existing.Email = profile.Email
	existing.FirstName = profile.FirstName
	existing.LastName = profile.LastName
	existing.Phone = profile.Phone
	if rolesJSON != "" {
		existing.Role = rolesJSON
	}
	if profile.Status != "" {
		existing.Status = profile.Status
	}
	existing.UpdatedAt = &now
	if err := q.UsersEntity.WithContext(ctx).Save(existing); err != nil {
		return nil, err
	}
	return toUserResult(existing), nil
}

func (r userRepository) UpdateProfile(ctx context.Context, idpUid string, patch repo.UserProfileUpdate) error {
	q := dao.Use(r.db)
	u := q.UsersEntity
	assigns := []field.AssignExpr{u.UpdatedAt.Value(time.Now())}
	if patch.FirstName != nil {
		assigns = append(assigns, u.FirstName.Value(*patch.FirstName))
	}
	if patch.LastName != nil {
		assigns = append(assigns, u.LastName.Value(*patch.LastName))
	}
	if patch.Phone != nil {
		assigns = append(assigns, u.Phone.Value(*patch.Phone))
	}
	_, err := u.WithContext(ctx).Where(u.IdpUid.Eq(idpUid)).UpdateSimple(assigns...)
	return err
}

func (r userRepository) UpdateStatus(ctx context.Context, idpUid string, status string) error {
	q := dao.Use(r.db)
	u := q.UsersEntity
	_, err := u.WithContext(ctx).
		Where(u.IdpUid.Eq(idpUid)).
		UpdateSimple(u.Status.Value(status), u.UpdatedAt.Value(time.Now()))
	return err
}

func (r userRepository) UpdateRoles(ctx context.Context, idpUid string, roles []string) error {
	q := dao.Use(r.db)
	u := q.UsersEntity
	_, err := u.WithContext(ctx).
		Where(u.IdpUid.Eq(idpUid)).
		UpdateSimple(u.Role.Value(encodeRoles(roles)), u.UpdatedAt.Value(time.Now()))
	return err
}

func (r userRepository) SoftDelete(ctx context.Context, idpUid string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		q := dao.Use(tx)
		u := q.UsersEntity
		if _, err := u.WithContext(ctx).
			Where(u.IdpUid.Eq(idpUid)).
			UpdateSimple(u.Status.Value("DEACTIVATED"), u.UpdatedAt.Value(time.Now())); err != nil {
			return err
		}
		if _, err := u.WithContext(ctx).Where(u.IdpUid.Eq(idpUid)).Delete(); err != nil {
			return err
		}
		return nil
	})
}

func encodeRoles(roles []string) string {
	if roles == nil {
		return ""
	}
	values := make([]constant.AuthRoleValue, len(roles))
	for i, r := range roles {
		values[i] = constant.AuthRoleValue(r)
	}
	b, err := json.Marshal(constant.AuthRoleJson{Values: values})
	if err != nil {
		return ""
	}
	return string(b)
}

func toUserResult(e *entity.UsersEntity) *repo.UserResult {
	var idpId *string
	if e.IdpUid != "" {
		v := e.IdpUid
		idpId = &v
	}
	return &repo.UserResult{
		Id:          e.Id,
		IdpId:       idpId,
		Email:       e.Email,
		FirstName:   e.FirstName,
		LastName:    e.LastName,
		Phone:       e.Phone,
		Roles:       parseRoles(e.Role),
		Status:      e.Status,
		LastLoginAt: e.LastLoginAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func parseRoles(s string) []string {
	if s == "" {
		return nil
	}
	var rj constant.AuthRoleJson
	if err := json.Unmarshal([]byte(s), &rj); err != nil {
		return nil
	}
	out := make([]string, len(rj.Values))
	for i, v := range rj.Values {
		out[i] = string(v)
	}
	return out
}
