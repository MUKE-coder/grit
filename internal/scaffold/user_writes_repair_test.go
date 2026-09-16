package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The users endpoints after M31 and the users half of M32: no check-then-insert,
// one transaction for the two writes an update makes, and a list that goes
// through paginate.
func TestUserWritesAreAtomicInTheTemplates(t *testing.T) {
	user := apiUserHandlerGo()
	for _, want := range []string{
		"services.CreateUser(c.Request.Context(), h.DB, &user)",
		"errors.Is(err, services.ErrEmailExists)",
		"writeUserRoleAssignment(tx, user.ID, req.Role)",
		"paginate.List[models.User](",
	} {
		if !strings.Contains(user, want) {
			t.Errorf("the user handler is missing %q", want)
		}
	}
	for _, gone := range []string{
		"// Check email uniqueness",
		"syncUserRoleAssignment(h.DB, user.ID, req.Role)",
		"\th.DB.WithContext(c.Request.Context()).Where(\"id = ?\", id).First(&user)\n",
		"\th.DB.WithContext(c.Request.Context()).Where(\"id = ?\", userID).First(&user)\n",
	} {
		if strings.Contains(user, gone) {
			t.Errorf("the user handler still has %q", gone)
		}
	}
	if strings.Contains(user, "\t\"math\"\n") || strings.Contains(user, "\t\"strconv\"\n") {
		t.Error("the user handler still imports the packages only the hand-rolled list used")
	}

	auth := apiAuthHandlerGo()
	if !strings.Contains(auth, "services.CreateUser(c.Request.Context(), h.DB, &user)") {
		t.Error("register still inserts the user itself")
	}
	if strings.Contains(auth, "var existingUser models.User") {
		t.Error("register still reads the email before inserting it")
	}

	write := apiUserWriteServiceGo()
	for _, want := range []string{"var ErrEmailExists", "func CreateUser(", "func IsDuplicateKey("} {
		if !strings.Contains(write, want) {
			t.Errorf("services/user_write.go is missing %q", want)
		}
	}
	for name, src := range map[string]string{
		"user_write.go":      apiUserWriteServiceGo(),
		"user_write_test.go": apiUserWriteServiceTestGo(),
	} {
		if _, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "example.com/app"))); err != nil {
			t.Errorf("%s is not valid Go: %v", name, err)
		}
	}
}

// The repair has to take the handler as it stood before this release to the
// handler the scaffold now writes, and say so rather than half-patch one that
// has been edited.
func TestRepairUserWritesReachesAnExistingProject(t *testing.T) {
	fresh := apiUserHandlerGo()
	if out, fixed, warn := repairUserHandlerWritesSource(fresh, "{{MODULE}}"); out != fresh || len(fixed)+len(warn) != 0 {
		t.Errorf("a fresh user.go needs no repair, got fixed %v warnings %v", fixed, warn)
	}

	old := previousUserHandler(t)
	out, fixed, warn := repairUserHandlerWritesSource(old, "{{MODULE}}")
	if len(warn) != 0 || len(fixed) != 4 {
		t.Fatalf("repairing the previous user.go: fixed %v warnings %v", fixed, warn)
	}
	for _, want := range []string{
		"services.CreateUser(c.Request.Context(), h.DB, &user)",
		"writeUserRoleAssignment(tx, user.ID, req.Role)",
		"paginate.List[models.User](",
		"\t\"errors\"",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired user.go is missing %q", want)
		}
	}
	if strings.Contains(out, "\t\"math\"\n") || strings.Contains(out, "\t\"strconv\"\n") {
		t.Error("the repaired user.go still imports math or strconv, so it will not compile")
	}
	if _, err := format.Source([]byte(strings.ReplaceAll(out, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("the repaired user.go is not valid Go: %v", err)
	}

	edited := strings.Replace(old, "Failed to update role assignment", "Could not update role assignment", 1)
	if _, _, warn := repairUserHandlerWritesSource(edited, "{{MODULE}}"); len(warn) != 1 {
		t.Errorf("an edited Update should be named and left alone, got %v", warn)
	}
}

func TestRepairRegisterWrites(t *testing.T) {
	fresh := apiAuthHandlerGo()
	if out, fixed, warn := repairRegisterWritesSource(fresh, "{{MODULE}}"); out != fresh || len(fixed)+len(warn) != 0 {
		t.Errorf("a fresh auth.go needs no repair, got fixed %v warnings %v", fixed, warn)
	}

	old := strings.Replace(fresh, registerInsertNew, userCreateInsertAnchor, 1)
	old = strings.Replace(old, "\tuser := models.User{\n\t\tFirstName:  req.FirstName,",
		registerCheckAnchor+"\tuser := models.User{\n\t\tFirstName:  req.FirstName,", 1)
	if !strings.Contains(old, "var existingUser models.User") || strings.Contains(old, "services.CreateUser(") {
		t.Fatal("could not rebuild auth.go as it was before")
	}
	out, fixed, warn := repairRegisterWritesSource(old, "{{MODULE}}")
	if len(warn) != 0 || len(fixed) != 1 {
		t.Fatalf("repairing the previous auth.go: fixed %v warnings %v", fixed, warn)
	}
	if !strings.Contains(out, "services.CreateUser(") || strings.Contains(out, "var existingUser models.User") {
		t.Error("register was not moved onto services.CreateUser")
	}
	if _, err := format.Source([]byte(strings.ReplaceAll(out, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("the repaired auth.go is not valid Go: %v", err)
	}
}

// previousUserHandler rebuilds user.go as it stood in v3.280.0 from the current
// template, so the repair is exercised against the text it will actually meet.
func previousUserHandler(t *testing.T) string {
	t.Helper()
	src := apiUserHandlerGo()

	src = strings.Replace(src, userCreateInsertNew, userCreateInsertAnchor, 1)
	src = strings.Replace(src, "\tuser := models.User{\n\t\tFirstName: req.FirstName,",
		userCreateCheckAnchor+"\tuser := models.User{\n\t\tFirstName: req.FirstName,", 1)
	src = strings.Replace(src, userRoleSyncNew, userRoleSyncAnchor, 1)
	src = strings.Replace(src, userUpdateWritesNew, userUpdateWritesAnchor, 1)
	src = strings.Replace(src, profileReloadNew, profileReloadAnchor, 1)

	start := strings.Index(src, "// userListConfig is the one description")
	end := strings.Index(src, "// GetByID returns a single user by ID.")
	if start < 0 || end < 0 {
		t.Fatal("could not find the users list in the template")
	}
	src = src[:start] + oldUserListForTest + src[end:]

	src = strings.Replace(src, "\t\"errors\"\n", "", 1)
	src = strings.Replace(src, "\t\"log\"\n", "\t\"log\"\n\t\"math\"\n", 1)
	src = strings.Replace(src, "\t\"net/http\"\n", "\t\"net/http\"\n\t\"strconv\"\n", 1)

	for _, marker := range []string{"// Check email uniqueness", "allowedSorts", "\t\"math\"\n", "\t\"strconv\"\n"} {
		if !strings.Contains(src, marker) {
			t.Fatalf("could not rebuild user.go as it was before: %q is missing", marker)
		}
	}
	return src
}

// oldUserListForTest is the hand-rolled list as v3.280.0 wrote it, counts
// splice included.
const oldUserListForTest = `// List returns a paginated list of users.
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	allowedSorts := map[string]bool{
		"id": true, "first_name": true, "last_name": true, "email": true, "role": true, "created_at": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	query := h.DB.WithContext(c.Request.Context()).Model(&models.User{})

	if search != "" {
		query = query.Where("LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	var total int64
	query.Count(&total)

	// The admin's stat cards ask for their counts on this request, over the
	// same search as the total.
	counts, err := paginate.Counts(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to count users",
			},
		})
		return
	}

	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Order(sortBy + " " + sortOrder).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch users",
			},
		})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
			"counts":    counts,
		},
	})
}

`
