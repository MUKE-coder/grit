package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The users endpoints, brought up to what the rest of the API already does.
//
// Three things, all of them in code Grit wrote and every project therefore has:
//
//   - Registration and the admin's Create User both SELECTed the email and
//     INSERTed when nothing came back. Two copies of a check that has a window
//     in it: a duplicate address arriving twice at once is answered with a 500.
//     Both now go through services.CreateUser and let the unique index decide.
//
//   - Update wrote the role assignment, committed it, and then wrote the row.
//     A failure in between left the user holding the permissions of the new role
//     while users.role still said the old one. Both writes are one transaction
//     now, and the grant cache is dropped after it commits rather than inside it.
//
//   - List was sixty-five lines of hand-rolled paging whose Count error was
//     discarded and which read none of the query parameters the admin's own
//     users page sends. ?role=, ?active= and ?provider= did nothing, so the
//     "Active Users" stat card counted every user.
//
// Each edit anchors on the text the scaffold wrote. Where it has been changed
// the edit is skipped and said out loud, so a handler somebody rewrote is never
// half-patched in silence.

const (
	userCreateCheckAnchor = `	// Check email uniqueness
	var existing models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "EMAIL_EXISTS",
				"message": "A user with this email already exists",
			},
		})
		return
	}

`

	userCreateInsertAnchor = `	if err := h.DB.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}
`

	userCreateInsertNew = `	// The unique index on users.email decides, not a SELECT before the INSERT.
	// Two requests for the same address arriving together both found nothing
	// and both inserted; one then failed on the constraint and was answered
	// with a 500 that said "Failed to create user".
	if err := services.CreateUser(c.Request.Context(), h.DB, &user); err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "EMAIL_EXISTS",
					"message": "A user with this email already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}
`

	registerCheckAnchor = `	// Check email uniqueness
	var existingUser models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "EMAIL_EXISTS",
				"message": "A user with this email already exists",
			},
		})
		return
	}

`

	registerInsertNew = `	// Same insert as the admin's Create User, through the same helper. The
	// check-then-insert this replaces was written out twice, and both copies
	// had the same gap between the SELECT and the INSERT: two signups for one
	// address in the same instant both passed the check, and the loser was
	// told "Failed to create user" with a 500.
	if err := services.CreateUser(c.Request.Context(), h.DB, &user); err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "EMAIL_EXISTS",
					"message": "A user with this email already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}
`

	userRoleSyncAnchor = `	txErr := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		if role.ID == "" {
			return nil // unknown name — fall back to the legacy string
		}
		return tx.Create(&models.UserRole{UserID: userID, RoleID: role.ID}).Error
	})
	if txErr != nil {
		return txErr
	}

	// Permissions just changed for this user; drop the cached grants.
	authz.Invalidate()
	return nil
}
`

	userRoleSyncNew = `	if err := db.Transaction(func(tx *gorm.DB) error {
		return writeUserRoleAssignment(tx, userID, roleName)
	}); err != nil {
		return err
	}

	// Permissions just changed for this user; drop the cached grants. After the
	// commit, never inside it: a request that resolved grants while the
	// transaction was still open would cache rows that had yet to exist.
	authz.Invalidate()
	return nil
}

// writeUserRoleAssignment is the same work with no transaction of its own, for
// a caller that already has one. Update needs it: the assignment and the row
// are one change and have to succeed or fail together.
//
// It does not invalidate the grant cache. Only the caller knows when its
// transaction commits, and that is the moment the cache is stale.
func writeUserRoleAssignment(tx *gorm.DB, userID, roleName string) error {
	var role models.Role
	err := tx.Where("name = ?", roleName).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}
	if role.ID == "" {
		return nil // unknown name — fall back to the legacy string
	}
	return tx.Create(&models.UserRole{UserID: userID, RoleID: role.ID}).Error
}
`

	userUpdateWritesAnchor = `	if req.Role != "" {
		if err := syncUserRoleAssignment(h.DB, user.ID, req.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to update role assignment",
				},
			})
			return
		}
	}

	if err := h.DB.WithContext(c.Request.Context()).Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update user",
			},
		})
		return
	}

	// Reload to get updated values
	h.DB.WithContext(c.Request.Context()).Where("id = ?", id).First(&user)
`

	userUpdateWritesNew = `	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if req.Role != "" {
			if err := writeUserRoleAssignment(tx, user.ID, req.Role); err != nil {
				return err
			}
		}
		return tx.Model(&user).Updates(updates).Error
	}); err != nil {
		if errors.Is(err, services.ErrEmailExists) || services.IsDuplicateKey(err) {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "EMAIL_EXISTS",
					"message": "A user with this email already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update user",
			},
		})
		return
	}
	if req.Role != "" {
		authz.Invalidate()
	}

	// Reload to get updated values. A failure here is reported rather than
	// ignored: the response used to carry whatever the struct happened to hold
	// when the read failed, which is the row as it was before the update.
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to reload the updated user",
			},
		})
		return
	}
`

	profileReloadAnchor = `	h.DB.WithContext(c.Request.Context()).Where("id = ?", userID).First(&user)
`

	profileReloadNew = `	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to reload your profile",
			},
		})
		return
	}
`

	userListNew = `// userListConfig is the one description of what the users list may be
// searched, sorted and filtered by.
//
// The filters are the point. The admin's users page ships a Role dropdown, a
// Status toggle and a Provider dropdown, and its "Active Users" stat card asks
// for /api/users?active=true&page_size=1. The hand-rolled list read none of
// them: every filter returned the whole table and the Active Users card
// reported the total user count.
var userListConfig = paginate.Config{
	Searchable: []string{"first_name", "last_name", "email"},
	Sortable: map[string]bool{
		"id": true, "first_name": true, "last_name": true,
		"email": true, "role": true, "created_at": true,
	},
	Filterable:   map[string]bool{"role": true, "active": true, "provider": true},
	DefaultSort:  "created_at",
	DefaultOrder: "desc",
}

// List returns a paginated list of users.
//
// Sixty-five lines of page clamping, sort whitelisting, search building and
// page arithmetic used to live here, with the Count error dropped on the
// floor: a count that failed left total at 0 and the table reported "0 of 0"
// over a full page of rows. paginate.List is the same logic every generated
// resource already uses, so a fix to it now reaches this endpoint too.
func (h *UserHandler) List(c *gin.Context) {
	res, err := paginate.List[models.User](
		h.DB.WithContext(c.Request.Context()).Model(&models.User{}),
		paginate.Bind(c),
		userListConfig,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch users",
			},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
`
)

// repairUserWrites applies M31 and the users half of M32 to an existing
// project. It runs after repairListCounts, which is a no-op once the list goes
// through paginate.List.
func repairUserWrites(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileContains(filepath.Join(apiRoot, "internal", "services", "user_write.go"), "func CreateUser(") {
		// writeFrameworkOwnedFiles delivers it, and a repair that calls a
		// function the project has not got is a project that will not compile.
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()
	files := []struct {
		path string
		fix  func(string, string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "handlers", "user.go"), repairUserHandlerWritesSource},
		{filepath.Join(apiRoot, "internal", "handlers", "auth.go"), repairRegisterWritesSource},
	}
	for _, f := range files {
		if !fileExists(f.path) {
			continue
		}
		fix := f.fix
		if err := repairSourceFile(root, m, f.path, func(src string) (string, []string, []string) {
			return fix(src, module)
		}); err != nil {
			return err
		}
	}
	return nil
}

// repairUserHandlerWritesSource fixes Create, List and Update in user.go.
func repairUserHandlerWritesSource(src, module string) (out string, fixed, warn []string) {
	if !strings.Contains(src, "func (h *UserHandler) Update(c *gin.Context) {") {
		return src, nil, nil
	}
	out = src

	// Create: let the unique index answer.
	if !strings.Contains(out, "services.CreateUser(") {
		if strings.Count(out, userCreateCheckAnchor) != 1 || strings.Count(out, userCreateInsertAnchor) != 1 {
			warn = append(warn, "Create still reads the email before inserting it, and is not the handler Grit wrote: insert through services.CreateUser, or two signups for one address in the same instant answer one of them with a 500")
		} else {
			out = strings.Replace(out, userCreateCheckAnchor, "", 1)
			out = strings.Replace(out, userCreateInsertAnchor, userCreateInsertNew, 1)
			fixed = append(fixed, "a duplicate email on Create is a 409 rather than a race with a 500")
		}
	}

	// Update: one transaction for the assignment and the row.
	if !strings.Contains(out, "writeUserRoleAssignment(") {
		if strings.Count(out, userRoleSyncAnchor) != 1 || strings.Count(out, userUpdateWritesAnchor) != 1 {
			warn = append(warn, "Update writes the role assignment and the user row separately, and is not the handler Grit wrote: put both in one h.DB.Transaction, or a failed update leaves the permissions changed and users.role behind")
		} else {
			out = strings.Replace(out, userRoleSyncAnchor, userRoleSyncNew, 1)
			out = strings.Replace(out, userUpdateWritesAnchor, userUpdateWritesNew, 1)
			fixed = append(fixed, "a user's role assignment and row are updated in one transaction")
		}
	}

	// UpdateProfile: the reload whose error was dropped.
	if strings.Count(out, profileReloadAnchor) == 1 {
		out = strings.Replace(out, profileReloadAnchor, profileReloadNew, 1)
		fixed = append(fixed, "a profile reload that fails is reported rather than answered with the old row")
	}

	// List: through paginate, with the filters the admin page sends.
	if !strings.Contains(out, "paginate.List[models.User]") {
		listed, ok := replaceUserList(out)
		if !ok {
			warn = append(warn, "the users list clamps and sorts by hand and reads neither ?role= nor ?active=, and is not the handler Grit wrote: list through paginate.List, or the users page's filters and its Active Users card are wrong")
		} else {
			out = listed
			fixed = append(fixed, "the users list goes through paginate.List, so its search, sort and filters work")
		}
	}

	if out == src {
		return src, fixed, warn
	}
	for _, imp := range []string{"errors", module + "/internal/paginate", module + "/internal/services"} {
		if added, ok := addImportGroup(out, imp); ok {
			out = added
		}
	}
	// gofmt formats, it does not prune: the hand-rolled list was the only
	// caller of these two, and a handler that imports what it does not use
	// stops compiling.
	for _, imp := range []string{"math", "strconv"} {
		if pruned, ok := dropUnusedImport(out, imp); ok {
			out = pruned
		}
	}
	return out, fixed, warn
}

// repairRegisterWritesSource fixes Register in auth.go.
func repairRegisterWritesSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *AuthHandler) Register(c *gin.Context) {") {
		return src, nil, nil
	}
	if strings.Contains(src, "services.CreateUser(") {
		return src, nil, nil
	}
	if strings.Count(src, registerCheckAnchor) != 1 || strings.Count(src, userCreateInsertAnchor) != 1 {
		return src, nil, []string{"Register still reads the email before inserting it, and is not the handler Grit wrote: insert through services.CreateUser, or two signups for one address in the same instant answer one of them with a 500"}
	}
	out := strings.Replace(src, registerCheckAnchor, "", 1)
	out = strings.Replace(out, userCreateInsertAnchor, registerInsertNew, 1)
	for _, imp := range []string{"errors", module + "/internal/services"} {
		if added, ok := addImportGroup(out, imp); ok {
			out = added
		}
	}
	return out, []string{"a duplicate email on register is a 409 rather than a race with a 500"}, nil
}

// replaceUserList swaps the hand-rolled List for the paginate one. It reports
// false unless the body is the one the scaffold wrote, markers and all.
func replaceUserList(src string) (string, bool) {
	head := "// List returns a paginated list of users.\nfunc (h *UserHandler) List(c *gin.Context) {\n"
	start := strings.Index(src, head)
	if start < 0 {
		return src, false
	}
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return src, false
	}
	end += start + len("\n}\n")
	body := src[start:end]
	for _, marker := range []string{
		"allowedSorts",
		"offset := (page - 1) * pageSize",
		"query.Count(&total)",
	} {
		if !strings.Contains(body, marker) {
			return src, false
		}
	}
	return src[:start] + userListNew + src[end:], true
}
