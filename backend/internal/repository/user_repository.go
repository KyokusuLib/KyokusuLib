package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/models/db"
	"github.com/lib/pq"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUsers(ctx context.Context, search string, limit int, offset int) ([]*db.User, error) {
	var query string
	var rows *sql.Rows
	var err error

	columns := `
		u.id, u.email,
		u.role, u.status, u.last_login,
		u.is_verified, u.isPublic,
		u.create_at,
		p.name, p.picture, p.banner, p.about, p.birthday, p.gender,
		t.tag,
		ups.is_show_tag, ups.is_show_bookmark`

	if search != "" {
		query = `
			SELECT` + columns + `
			FROM users u
			LEFT JOIN user_profiles p ON p.user_id = u.id
			LEFT JOIN user_tags t ON t.id = p.tag_id
			LEFT JOIN user_profile_settings ups ON ups.user_id = u.id
			WHERE p.name ILIKE $1::text OR u.email ILIKE $1::text
			ORDER BY u.id DESC LIMIT $2 OFFSET $3`
		rows, err = r.DB.QueryContext(ctx, query, "%"+search+"%", limit, offset)
	} else {
		query = `
			SELECT` + columns + `
			FROM users u
			LEFT JOIN user_profiles p ON p.user_id = u.id
			LEFT JOIN user_tags t ON t.id = p.tag_id
			LEFT JOIN user_profile_settings ups ON ups.user_id = u.id
			ORDER BY u.id DESC LIMIT $1 OFFSET $2`
		rows, err = r.DB.QueryContext(ctx, query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*db.User
	for rows.Next() {
		var u db.User
		var (
			name                       sql.NullString
			picture                    sql.NullString
			banner                     sql.NullString
			about                      sql.NullString
			birthday                   sql.NullTime
			gender                     sql.NullString
			tag                        sql.NullString
			isShowTag                  sql.NullBool
			isShowBookmark             sql.NullBool
		)

		if err := rows.Scan(
			&u.ID, &u.Email,
			&u.Role, &u.Status, &u.LastLogin,
			&u.IsVerified, &u.IsPublic,
			&u.CreateAt,
			&name, &picture, &banner, &about, &birthday, &gender,
			&tag,
			&isShowTag, &isShowBookmark,
		); err != nil {
			return nil, err
		}

		if name.Valid {
			u.Name = name.String
		}
		if picture.Valid {
			u.Picture = picture.String
		}
		if banner.Valid {
			u.Banner = banner.String
		}
		if about.Valid {
			u.About = about.String
		}
		if birthday.Valid {
			u.Birthday = &birthday.Time
		}
		if gender.Valid {
			u.Gender = db.UserGenere(gender.String)
		}
		if tag.Valid {
			u.Tag = tag.String
		}
		if isShowTag.Valid {
			u.IsShowTag = isShowTag.Bool
		}
		if isShowBookmark.Valid {
			u.IsShowBookmark = isShowBookmark.Bool
		}

		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByEmail(email string) (*db.User, error) {
	return r.findOne("email", email)
}

func (r *UserRepository) GetByID(id int) (*db.User, error) {
	return r.findOne("id", id)
}

func (r *UserRepository) GetByVerificationToken(token string) (*db.User, error) {
	return r.findOne("verification_token", token)
}

func (r *UserRepository) GetByResetToken(token string) (*db.User, error) {
	return r.findOne("reset_token", token)
}

func (r *UserRepository) findOne(field string, value interface{}) (*db.User, error) {
	u := &db.User{}

	var (
		passwordHash                 sql.NullString
		verificationToken            sql.NullString
		verificationTokenExpiresAt   sql.NullTime
		resetToken                   sql.NullString
		resetTokenExpiresAt          sql.NullTime
		name                         sql.NullString
		picture                      sql.NullString
		banner                       sql.NullString
		about                        sql.NullString
		birthday                     sql.NullTime
		gender                       sql.NullString
		isShowTag                    sql.NullBool
		isShowBookmark               sql.NullBool
	)

	query := `
		SELECT 
			u.id, u.email, u.password_hash, u.role, u.status, u.last_login,
			u.is_verified, u.isPublic,
			u.verification_token, u.verification_token_expires_at,
			u.reset_token, u.reset_token_expires_at,
			u.create_at,
			p.name, p.picture, p.banner, p.about, p.birthday, p.gender,
			t.tag,
			ups.is_show_tag, ups.is_show_bookmark
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		LEFT JOIN user_tags t ON t.id = p.tag_id
		LEFT JOIN user_profile_settings ups ON ups.user_id = u.id
		WHERE u.` + field + ` = $1`

	err := r.DB.QueryRow(query, value).Scan(
		&u.ID,
		&u.Email,
		&passwordHash,
		&u.Role,
		&u.Status,
		&u.LastLogin,
		&u.IsVerified,
		&u.IsPublic,
		&verificationToken,
		&verificationTokenExpiresAt,
		&resetToken,
		&resetTokenExpiresAt,
		&u.CreateAt,
		&name,
		&picture,
		&banner,
		&about,
		&birthday,
		&gender,
		&u.Tag,
		&isShowTag,
		&isShowBookmark,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if passwordHash.Valid {
		u.PasswordHash = passwordHash.String
	}
	if verificationToken.Valid {
		u.VerificationToken = verificationToken.String
	}
	if verificationTokenExpiresAt.Valid {
		u.VerificationTokenExpiresAt = &verificationTokenExpiresAt.Time
	}
	if resetToken.Valid {
		u.ResetToken = resetToken.String
	}
	if resetTokenExpiresAt.Valid {
		u.ResetTokenExpiresAt = &resetTokenExpiresAt.Time
	}
	if name.Valid {
		u.Name = name.String
	}
	if picture.Valid {
		u.Picture = picture.String
	}
	if about.Valid {
		u.About = about.String
	}
	if birthday.Valid {
		u.Birthday = &birthday.Time
	}
	if gender.Valid {
		u.Gender = db.UserGenere(gender.String)
	}
	if banner.Valid {
		u.Banner = banner.String
	}
	if isShowTag.Valid {
		u.IsShowTag = isShowTag.Bool
	}
	if isShowBookmark.Valid {
		u.IsShowBookmark = isShowBookmark.Bool
	}
    
	return u, nil
}

func (r *UserRepository) Create(u *db.User) error {
	if u.Role == "" {
		u.Role = "user"
	}

	var passwordHash interface{} = nil
	if u.PasswordHash != "" {
		passwordHash = u.PasswordHash
	}

	var verificationToken interface{} = nil
	if u.VerificationToken != "" {
		verificationToken = u.VerificationToken
	}

	var verificationTokenExpiresAt interface{} = nil
	if u.VerificationTokenExpiresAt != nil {
		verificationTokenExpiresAt = u.VerificationTokenExpiresAt
	}

	query := `
		INSERT INTO users 
			(email, password_hash, role, status, last_login, is_verified, verification_token, verification_token_expires_at) 
		VALUES ($1, $2, $3, 'online', $4, $5, $6, $7) 
		RETURNING id`

	err := r.DB.QueryRow(query,
		u.Email,
		passwordHash,
		u.Role,
		time.Now(),
		u.IsVerified,
		verificationToken,
		verificationTokenExpiresAt,
	).Scan(&u.ID)
	
	if err != nil {
		return err
	}

	return r.createOrUpdateProfile(u)
}

func (r *UserRepository) Update(u *db.User) error {
	var resetToken interface{} = nil
	if u.ResetToken != "" {
		resetToken = u.ResetToken
	}

	var resetTokenExpiresAt interface{} = nil
	if u.ResetTokenExpiresAt != nil {
		resetTokenExpiresAt = u.ResetTokenExpiresAt
	}

	query := `
		UPDATE users 
		SET password_hash = $1, 
		    reset_token = $2, 
		    reset_token_expires_at = $3,
		    status = $4,
		    last_login = $5,
		    is_verified = $6
		WHERE id = $7`

	_, err := r.DB.Exec(query,
		u.PasswordHash,
		resetToken,
		resetTokenExpiresAt,
		u.Status,
		u.LastLogin,
		u.IsVerified,
		u.ID,
	)
	return err
}

func (r *UserRepository) UpdateProfile(u *db.User) error {
	return r.createOrUpdateProfile(u)
}

func (r *UserRepository) createOrUpdateProfile(u *db.User) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var tagID int
	if u.Tag != "" {
		err := tx.QueryRow("SELECT id FROM user_tags WHERE tag = $1", u.Tag).Scan(&tagID)	
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				tagID = 1
			} else {
				return err
			}
		}
	} else {
		tagID = 1 
	}

	query := `
		INSERT INTO user_profiles (user_id, name, picture, about, birthday, gender, tag_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE
		SET name = EXCLUDED.name,
		    picture = EXCLUDED.picture,
		    about = EXCLUDED.about,
		    birthday = EXCLUDED.birthday,
		    gender = EXCLUDED.gender,
		    tag_id = EXCLUDED.tag_id`

	_, err = tx.Exec(query,
		u.ID,
		u.Name,
		u.Picture,
		u.About,
		u.Birthday,
		u.Gender,
		tagID,
	)
	if err != nil {
		return err
	}

	query_public := `UPDATE users SET ispublic = $1 WHERE id = $2`
	_, err = tx.Exec(query_public, u.IsPublic, u.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) DeleteExpiredUnverified() error {
	query := `
		DELETE FROM users 
		WHERE is_verified = FALSE 
		  AND verification_token_expires_at < $1`
	_, err := r.DB.Exec(query, time.Now())
	return err
}

func (r *UserRepository) UpdateStatus(userID int, status string) error {
	query := `UPDATE users SET status = $1, last_login = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, status, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateDtoStatus(
    ctx context.Context,
    userID int,
    status constants.UserStatus,
    lastActive time.Time,
) error {
    query := `UPDATE users SET status = $1, last_login = $2 WHERE id = $3`
    _, err := r.DB.ExecContext(ctx, query, status.String(), lastActive, userID)
    return err
}

func (r *UserRepository) UpdateAvatar(userID int, avatar string) error {
	query := `UPDATE user_profiles SET picture = $1 WHERE user_id = $2`
	_, err := r.DB.Exec(query, avatar, userID)
	return err
}

func (r *UserRepository) UpdateBanner(userID int, banner string) error {
	query := `UPDATE user_profiles SET banner = $1 WHERE user_id = $2`
	_, err := r.DB.Exec(query, banner, userID)
	return err
}


func (r *UserRepository) GetAvatarUrl(userID int) (string, error) {
	var picture sql.NullString

	query := `
		SELECT p.picture
		FROM user_profiles p
		WHERE p.user_id = $1`

	err := r.DB.QueryRow(query, userID).Scan(&picture)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	if picture.Valid {
		return picture.String, nil
	}
	return "", nil
}

func (r *UserRepository) GetBannerUrl(userID int) (string, error) {
	var picture sql.NullString

	query := `
		SELECT p.banner
		FROM user_profiles p
		WHERE p.user_id = $1`

	err := r.DB.QueryRow(query, userID).Scan(&picture)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	if picture.Valid {
		return picture.String, nil
	}
	return "", nil
}

func (r *UserRepository) MarkUserVerified(userID int) error {
	query := `
		UPDATE users 
		SET is_verified = TRUE, 
		    verification_token = NULL, 
		    verification_token_expires_at = NULL 
		WHERE id = $1`
	_, err := r.DB.Exec(query, userID)
	return err
}

func (r *UserRepository) GetUserTags(ctx context.Context, userID int) ([]*db.UserTag, error) {
    const query = `
        SELECT ut.id, ut.tag
        FROM user_tags ut
        JOIN users_user_tags uut ON ut.id = uut.tag_id
        WHERE uut.user_id = $1
        ORDER BY ut.tag
    `

    rows, err := r.DB.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tags []*db.UserTag
    for rows.Next() {
        var t db.UserTag
        if err := rows.Scan(&t.TagID, &t.Tag); err != nil {
            return nil, err
        }
        tags = append(tags, &t)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return tags, nil
}

func (r *UserRepository) UpdateUserTag(ctx context.Context, userID int, tagID int) error {
	query := `UPDATE user_profiles
		SET tag_id = $2
		WHERE user_id = $1`
		
	_, err := r.DB.Exec(query, tagID, userID)
	return err
}

func (r *UserRepository) HasUserTag(ctx context.Context, userID int, tagID int) (bool, error) {
	query := `SELECT 1 FROM users_user_tags WHERE user_id = $1 AND tag_id = $2`
	var exists int
	err := r.DB.QueryRowContext(ctx, query, userID, tagID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UserRepository) GrantUserTag(ctx context.Context, userID int, tagID int) error {
	query := `INSERT INTO users_user_tags (user_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.DB.ExecContext(ctx, query, userID, tagID)
	return err
}

type UpdateUserParams struct {
	Name           string
	About          string
	Gender         string
	Birthday       *time.Time
	IsPublic       bool
	Role           string
	Status         string
	IsShowTag      bool
	IsShowBookmark bool
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID int, p UpdateUserParams) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET role = $1, status = $2, ispublic = $3 WHERE id = $4`,
		p.Role, p.Status, p.IsPublic, userID,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_profiles (user_id, name, about, birthday, gender)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
		SET name = EXCLUDED.name,
		    about = EXCLUDED.about,
		    birthday = EXCLUDED.birthday,
		    gender = EXCLUDED.gender`,
		userID, p.Name, p.About, p.Birthday, p.Gender,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_profile_settings (user_id, is_show_tag, is_show_bookmark)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET is_show_tag = EXCLUDED.is_show_tag,
		    is_show_bookmark = EXCLUDED.is_show_bookmark`,
		userID, p.IsShowTag, p.IsShowBookmark,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetReadChaptersCount(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM read_chapters WHERE user_id = $1`
	var count int
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *UserRepository) GetUserStats(ctx context.Context, userID int) (int, int, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM novela_comments WHERE user_id = $1),
			(SELECT COUNT(*) FROM read_chapters WHERE user_id = $1)`
	var totalComments, readChapters int
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&totalComments, &readChapters)
	return totalComments, readChapters, err
}

func (r *UserRepository) GetUserStatsBatch(ctx context.Context, userIDs []int) (map[int]db.UserStats, error) {
	result := make(map[int]db.UserStats, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT
			u.id,
			(SELECT COUNT(*) FROM novela_comments nc WHERE nc.user_id = u.id),
			(SELECT COUNT(*) FROM read_chapters rc WHERE rc.user_id = u.id)
		FROM users u
		WHERE u.id = ANY($1)`

	rows, err := r.DB.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var stats db.UserStats
		if err := rows.Scan(&id, &stats.TotalComments, &stats.ReadChapters); err != nil {
			return nil, err
		}
		result[id] = stats
	}

	return result, rows.Err()
}

func (r *UserRepository) GetUserTagsBatch(ctx context.Context, userIDs []int) (map[int][]*db.UserTag, error) {
	result := make(map[int][]*db.UserTag, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT uut.user_id, ut.id, ut.tag
		FROM users_user_tags uut
		JOIN user_tags ut ON ut.id = uut.tag_id
		WHERE uut.user_id = ANY($1)
		ORDER BY ut.tag`

	rows, err := r.DB.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var t db.UserTag
		if err := rows.Scan(&userID, &t.TagID, &t.Tag); err != nil {
			return nil, err
		}
		result[userID] = append(result[userID], &t)
	}

	return result, rows.Err()
}

func (r *UserRepository) IsExist(ctx context.Context, userID int) (bool, error) {
    query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

    var exists bool
    err := r.DB.QueryRowContext(ctx, query, userID).Scan(&exists)
    if err != nil {
        return false, err
    }

    return exists, nil
}