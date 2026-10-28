package plugin

// kanbanHandlerGo emits internal/handlers/board.go.
//
// The handlers are deliberately thin: bind, call the service, respond. The
// one that earns its keep is MoveCard, and even there the interesting line is
// a single Updates call. Everything that makes the move correct lives in
// BoardService and internal/kanban.
func kanbanHandlerGo(ctx Context) string {
	return goSource(ctx, `package handlers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/events"
	"{{MODULE}}/internal/kanban"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/respond"
	"{{MODULE}}/internal/services"
)

// BoardHandler serves Kanban boards.
type BoardHandler struct {
	DB      *gorm.DB
	Service *services.BoardService
}

func NewBoardHandler(db *gorm.DB) *BoardHandler {
	return &BoardHandler{DB: db, Service: services.NewBoardService(db)}
}

// boardActorID reads the signed-in user from the gin context.
//
// Named for this plugin rather than something general, because the handlers
// package is shared: a plugin that defines currentUserID collides with the
// next plugin that does the same, and the project stops compiling with an
// error that names neither of them.
func boardActorID(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ---------------------------------------------------------------- boards ---

// ListBoards returns every board, newest first. Boards are team data, so
// there is no per-user filter here; the route is behind auth.
func (h *BoardHandler) ListBoards(c *gin.Context) {
	var boards []models.Board
	if err := h.DB.Order("created_at desc").Find(&boards).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, boards)
}

// GetBoard returns one board with its columns and cards, each in rank order.
func (h *BoardHandler) GetBoard(c *gin.Context) {
	board, err := h.Service.LoadBoard(c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.NotFound(c, "Board not found")
			return
		}
		respond.Internal(c, err)
		return
	}
	respond.OK(c, board)
}

// CreateBoard makes a board with three starting columns, so that creating one
// lands on a board somebody can use rather than an empty screen.
func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var in struct {
		Name        string ~json:"name" binding:"required"~
		Description string ~json:"description"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "A board needs a name")
		return
	}

	board := models.Board{
		Name:        in.Name,
		Description: in.Description,
		CreatedBy:   boardActorID(c),
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&board).Error; err != nil {
			return err
		}
		cols, err := h.Service.DefaultColumns(board.ID)
		if err != nil {
			return err
		}
		return tx.Create(&cols).Error
	})
	if err != nil {
		respond.Internal(c, err)
		return
	}

	events.Emitted(c, "boards", "Board", "created", board.ID, board.Name, "", nil, board)
	respond.Created(c, board, "Board created")
}

// UpdateBoard renames a board or changes its description.
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	var board models.Board
	if err := h.DB.First(&board, "id = ?", c.Param("id")).Error; err != nil {
		respond.NotFound(c, "Board not found")
		return
	}

	var in struct {
		Name        *string ~json:"name"~
		Description *string ~json:"description"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Could not read the update")
		return
	}

	updates := map[string]any{}
	if in.Name != nil {
		if *in.Name == "" {
			respond.BadRequest(c, "A board needs a name")
			return
		}
		updates["name"] = *in.Name
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if len(updates) == 0 {
		respond.OK(c, board)
		return
	}

	before := board
	if err := h.DB.Model(&board).Updates(updates).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	events.Emitted(c, "boards", "Board", "updated", board.ID, board.Name, "", before, board)
	respond.OK(c, board, "Board updated")
}

// DeleteBoard removes the board and everything on it.
func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	id := c.Param("id")

	var board models.Board
	if err := h.DB.First(&board, "id = ?", id).Error; err != nil {
		respond.NotFound(c, "Board not found")
		return
	}
	if err := h.Service.DeleteBoard(id); err != nil {
		respond.Internal(c, err)
		return
	}
	events.Emitted(c, "boards", "Board", "deleted", id, board.Name, "", board, nil)
	respond.OK(c, gin.H{"id": id}, "Board deleted")
}

// --------------------------------------------------------------- columns ---

// CreateColumn appends a column to the right-hand end of the board.
func (h *BoardHandler) CreateColumn(c *gin.Context) {
	boardID := c.Param("id")
	if err := h.DB.First(&models.Board{}, "id = ?", boardID).Error; err != nil {
		respond.NotFound(c, "Board not found")
		return
	}

	var in struct {
		Name     string ~json:"name" binding:"required"~
		Color    string ~json:"color"~
		WIPLimit int    ~json:"wip_limit"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "A column needs a name")
		return
	}

	rank, err := h.Service.RankForColumn(boardID, "", "", "")
	if err != nil {
		respond.Internal(c, err)
		return
	}

	col := models.BoardColumn{
		BoardID:  boardID,
		Name:     in.Name,
		Color:    in.Color,
		Rank:     rank,
		WIPLimit: in.WIPLimit,
	}
	if err := h.DB.Create(&col).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	respond.Created(c, col, "Column added")
}

// UpdateColumn renames a column, recolours it, or sets its WIP limit.
func (h *BoardHandler) UpdateColumn(c *gin.Context) {
	var col models.BoardColumn
	if err := h.DB.First(&col, "id = ?", c.Param("id")).Error; err != nil {
		respond.NotFound(c, "Column not found")
		return
	}

	var in struct {
		Name     *string ~json:"name"~
		Color    *string ~json:"color"~
		WIPLimit *int    ~json:"wip_limit"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Could not read the update")
		return
	}

	updates := map[string]any{}
	if in.Name != nil {
		if *in.Name == "" {
			respond.BadRequest(c, "A column needs a name")
			return
		}
		updates["name"] = *in.Name
	}
	if in.Color != nil {
		updates["color"] = *in.Color
	}
	if in.WIPLimit != nil {
		if *in.WIPLimit < 0 {
			respond.BadRequest(c, "A work-in-progress limit cannot be negative")
			return
		}
		// Deliberately not checked against the column's current card count.
		// Setting a limit below what a column already holds is how a team
		// starts enforcing one: the column drains to the limit instead of
		// the setting being refused.
		updates["wip_limit"] = *in.WIPLimit
	}
	if len(updates) == 0 {
		respond.OK(c, col)
		return
	}

	if err := h.DB.Model(&col).Updates(updates).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, col, "Column updated")
}

// DeleteColumn removes a column and the cards in it.
func (h *BoardHandler) DeleteColumn(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.First(&models.BoardColumn{}, "id = ?", id).Error; err != nil {
		respond.NotFound(c, "Column not found")
		return
	}
	if err := h.Service.DeleteColumn(id); err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, gin.H{"id": id}, "Column deleted")
}

// MoveColumn reorders a column on its board.
func (h *BoardHandler) MoveColumn(c *gin.Context) {
	var col models.BoardColumn
	if err := h.DB.First(&col, "id = ?", c.Param("id")).Error; err != nil {
		respond.NotFound(c, "Column not found")
		return
	}

	var in struct {
		PrevID string ~json:"prev_id"~
		NextID string ~json:"next_id"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Could not read the move")
		return
	}

	rank, err := h.Service.RankForColumn(col.BoardID, in.PrevID, in.NextID, col.ID)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	if err := h.DB.Model(&col).Update("rank", rank).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, col, "Column moved")
}

// ----------------------------------------------------------------- cards ---

// CreateCard appends a card to the bottom of a column.
func (h *BoardHandler) CreateCard(c *gin.Context) {
	boardID := c.Param("id")

	var in struct {
		ColumnID    string     ~json:"column_id" binding:"required"~
		Title       string     ~json:"title" binding:"required"~
		Description string     ~json:"description"~
		AssigneeID  *string    ~json:"assignee_id"~
		DueDate     *time.Time ~json:"due_date"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "A card needs a title and a column")
		return
	}

	// The column has to be on this board. Without this check a crafted
	// column_id would file a card onto a board the caller was not looking at.
	var col models.BoardColumn
	if err := h.DB.First(&col, "id = ? AND board_id = ?", in.ColumnID, boardID).Error; err != nil {
		respond.NotFound(c, "Column not found on this board")
		return
	}

	if err := h.Service.CheckWIP(col.ID, ""); err != nil {
		if errors.Is(err, services.ErrWIPLimit) {
			respond.Conflict(c, err.Error())
			return
		}
		respond.Internal(c, err)
		return
	}

	rank, err := h.Service.RankForCard(col.ID, "", "", "")
	if err != nil {
		respond.Internal(c, err)
		return
	}

	card := models.BoardCard{
		BoardID:     boardID,
		ColumnID:    col.ID,
		Title:       in.Title,
		Description: in.Description,
		Rank:        rank,
		AssigneeID:  in.AssigneeID,
		DueDate:     in.DueDate,
		CreatedBy:   boardActorID(c),
	}
	if err := h.DB.Create(&card).Error; err != nil {
		respond.Internal(c, err)
		return
	}

	events.Emitted(c, "boards", "BoardCard", "card_created", card.ID, card.Title, "", nil, card)
	respond.Created(c, card, "Card added")
}

// UpdateCard edits a card in place. It does not move it: see MoveCard.
func (h *BoardHandler) UpdateCard(c *gin.Context) {
	var card models.BoardCard
	if err := h.DB.First(&card, "id = ?", c.Param("id")).Error; err != nil {
		respond.NotFound(c, "Card not found")
		return
	}

	var in struct {
		Title       *string    ~json:"title"~
		Description *string    ~json:"description"~
		AssigneeID  *string    ~json:"assignee_id"~
		DueDate     *time.Time ~json:"due_date"~
		ClearDue    bool       ~json:"clear_due_date"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Could not read the update")
		return
	}

	updates := map[string]any{}
	if in.Title != nil {
		if *in.Title == "" {
			respond.BadRequest(c, "A card needs a title")
			return
		}
		updates["title"] = *in.Title
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.AssigneeID != nil {
		if *in.AssigneeID == "" {
			updates["assignee_id"] = nil
		} else {
			updates["assignee_id"] = *in.AssigneeID
		}
	}
	// A null due_date in JSON is indistinguishable from an absent one once it
	// reaches a *time.Time, so clearing a date needs its own flag rather than
	// being silently impossible.
	if in.ClearDue {
		updates["due_date"] = nil
	} else if in.DueDate != nil {
		updates["due_date"] = *in.DueDate
	}

	if len(updates) == 0 {
		respond.OK(c, card)
		return
	}

	before := card
	if err := h.DB.Model(&card).Updates(updates).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	events.Emitted(c, "boards", "BoardCard", "card_updated", card.ID, card.Title, "", before, card)
	respond.OK(c, card, "Card updated")
}

// DeleteCard removes a card.
func (h *BoardHandler) DeleteCard(c *gin.Context) {
	var card models.BoardCard
	if err := h.DB.First(&card, "id = ?", c.Param("id")).Error; err != nil {
		respond.NotFound(c, "Card not found")
		return
	}
	if err := h.DB.Delete(&card).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	events.Emitted(c, "boards", "BoardCard", "card_deleted", card.ID, card.Title, "", card, nil)
	respond.OK(c, gin.H{"id": card.ID}, "Card deleted")
}

// MoveCard is the whole point of the plugin.
//
// The request names the destination column and the two cards the dragged card
// was dropped between. The response is one row written:
//
//	UPDATE board_cards SET column_id = ?, rank = ? WHERE id = ?
//
// No other card is touched, nothing is renumbered, and no transaction is held
// across rows somebody else might be dragging at the same moment.
func (h *BoardHandler) MoveCard(c *gin.Context) {
	var card models.BoardCard
	if err := h.DB.First(&card, "id = ?", c.Param("id")).Error; err != nil {
		respond.NotFound(c, "Card not found")
		return
	}

	var in struct {
		ColumnID string ~json:"column_id" binding:"required"~
		// The card immediately above the drop, and the one immediately below.
		// Either may be empty: no prev means the top of the column, no next
		// means the bottom.
		PrevID string ~json:"prev_id"~
		NextID string ~json:"next_id"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "A move needs a destination column")
		return
	}

	// A card moves between columns of one board, never between boards.
	var col models.BoardColumn
	if err := h.DB.First(&col, "id = ? AND board_id = ?", in.ColumnID, card.BoardID).Error; err != nil {
		respond.NotFound(c, "Column not found on this board")
		return
	}

	if col.ID != card.ColumnID {
		if err := h.Service.CheckWIP(col.ID, card.ID); err != nil {
			if errors.Is(err, services.ErrWIPLimit) {
				respond.Conflict(c, err.Error())
				return
			}
			respond.Internal(c, err)
			return
		}
	}

	rank, err := h.Service.RankForCard(col.ID, in.PrevID, in.NextID, card.ID)
	if err != nil {
		respond.Internal(c, err)
		return
	}

	before := card
	if err := h.DB.Model(&card).Updates(map[string]any{
		"column_id": col.ID,
		"rank":      rank,
	}).Error; err != nil {
		respond.Internal(c, err)
		return
	}

	events.Emitted(c, "boards", "BoardCard", "card_moved", card.ID, card.Title,
		"moved to "+col.Name, before, card)
	respond.OK(c, card, "Card moved")
}

// compile-time assertion that the rank package is wired in; it is used by the
// service rather than here, and an unused import would be the only sign that
// somebody removed it.
var _ = kanban.Digits
`)
}

// kanbanHandlerTestGo emits internal/handlers/board_test.go.
//
// The tests that matter are the ones about ordering, because that is what a
// board is. Each one states the behaviour a user would notice if it broke.
func kanbanHandlerTestGo(ctx Context) string {
	return goSource(ctx, `package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

func boardTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// One connection: every new connection to :memory: is its own empty
	// database, so a pool turns this into a flaky test.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.Board{}, &models.BoardColumn{}, &models.BoardCard{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Migrator().DropTable(&models.BoardCard{}, &models.BoardColumn{}, &models.Board{})
	})
	return db
}

// seedBoard makes a board with one column holding n cards titled "1".."n".
func seedBoard(t *testing.T, db *gorm.DB, n int) (models.Board, models.BoardColumn, []models.BoardCard) {
	t.Helper()
	svc := services.NewBoardService(db)

	board := models.Board{Name: "Test board"}
	if err := db.Create(&board).Error; err != nil {
		t.Fatalf("board: %v", err)
	}
	cols, err := svc.DefaultColumns(board.ID)
	if err != nil {
		t.Fatalf("columns: %v", err)
	}
	if err := db.Create(&cols).Error; err != nil {
		t.Fatalf("columns: %v", err)
	}

	col := cols[0]
	cards := make([]models.BoardCard, 0, n)
	for i := 0; i < n; i++ {
		rank, err := svc.RankForCard(col.ID, "", "", "")
		if err != nil {
			t.Fatalf("rank: %v", err)
		}
		card := models.BoardCard{
			BoardID:  board.ID,
			ColumnID: col.ID,
			Title:    string(rune('1' + i)),
			Rank:     rank,
		}
		if err := db.Create(&card).Error; err != nil {
			t.Fatalf("card: %v", err)
		}
		cards = append(cards, card)
	}
	return board, col, cards
}

func boardRouter(db *gorm.DB) (*gin.Engine, *BoardHandler) {
	gin.SetMode(gin.TestMode)
	h := NewBoardHandler(db)
	r := gin.New()
	r.POST("/board-cards/:id/move", h.MoveCard)
	r.GET("/boards/:id", h.GetBoard)
	return r, h
}

func moveCard(t *testing.T, r *gin.Engine, cardID string, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/board-cards/"+cardID+"/move", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// titles reads a column back the way the board does, and returns the card
// titles in the order they would appear on screen.
func titles(t *testing.T, db *gorm.DB, columnID string) string {
	t.Helper()
	var cards []models.BoardCard
	if err := db.Where("column_id = ?", columnID).Order("rank asc, id asc").Find(&cards).Error; err != nil {
		t.Fatalf("read column: %v", err)
	}
	out := ""
	for _, c := range cards {
		out += c.Title
	}
	return out
}

func TestMoveCardToTopOfItsColumn(t *testing.T) {
	db := boardTestDB(t)
	_, col, cards := seedBoard(t, db, 3)
	r, _ := boardRouter(db)

	if got := titles(t, db, col.ID); got != "123" {
		t.Fatalf("seeded order is %q, want %q", got, "123")
	}

	// Drag "3" above "1": no card above it, "1" below it.
	w := moveCard(t, r, cards[2].ID, map[string]any{
		"column_id": col.ID,
		"next_id":   cards[0].ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("move: %d %s", w.Code, w.Body.String())
	}
	if got := titles(t, db, col.ID); got != "312" {
		t.Fatalf("after moving 3 to the top the order is %q, want %q", got, "312")
	}
}

func TestMoveCardBetweenTwoOthers(t *testing.T) {
	db := boardTestDB(t)
	_, col, cards := seedBoard(t, db, 4)
	r, _ := boardRouter(db)

	// Drag "4" between "1" and "2".
	w := moveCard(t, r, cards[3].ID, map[string]any{
		"column_id": col.ID,
		"prev_id":   cards[0].ID,
		"next_id":   cards[1].ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("move: %d %s", w.Code, w.Body.String())
	}
	if got := titles(t, db, col.ID); got != "1423" {
		t.Fatalf("order is %q, want %q", got, "1423")
	}
}

// TestMoveCardWritesOneRow is the property the whole design exists for. If
// this ever fails, a move has started renumbering, and two people dragging at
// once will start losing each other's work.
func TestMoveCardWritesOneRow(t *testing.T) {
	db := boardTestDB(t)
	_, col, cards := seedBoard(t, db, 5)
	r, _ := boardRouter(db)

	var before []models.BoardCard
	if err := db.Where("column_id = ?", col.ID).Order("id asc").Find(&before).Error; err != nil {
		t.Fatal(err)
	}

	w := moveCard(t, r, cards[4].ID, map[string]any{
		"column_id": col.ID,
		"next_id":   cards[0].ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("move: %d %s", w.Code, w.Body.String())
	}

	var after []models.BoardCard
	if err := db.Where("column_id = ?", col.ID).Order("id asc").Find(&after).Error; err != nil {
		t.Fatal(err)
	}

	changed := 0
	for i := range before {
		if before[i].Rank != after[i].Rank || before[i].ColumnID != after[i].ColumnID {
			changed++
		}
	}
	if changed != 1 {
		t.Fatalf("moving one card changed %d rows; a move must write exactly one", changed)
	}
}

func TestMoveCardToAnotherColumn(t *testing.T) {
	db := boardTestDB(t)
	board, col, cards := seedBoard(t, db, 2)
	r, _ := boardRouter(db)

	var other models.BoardColumn
	if err := db.Where("board_id = ? AND id <> ?", board.ID, col.ID).First(&other).Error; err != nil {
		t.Fatal(err)
	}

	w := moveCard(t, r, cards[0].ID, map[string]any{"column_id": other.ID})
	if w.Code != http.StatusOK {
		t.Fatalf("move: %d %s", w.Code, w.Body.String())
	}
	if got := titles(t, db, col.ID); got != "2" {
		t.Fatalf("source column is %q, want %q", got, "2")
	}
	if got := titles(t, db, other.ID); got != "1" {
		t.Fatalf("destination column is %q, want %q", got, "1")
	}
}

// TestMoveCardRejectsAColumnOnAnotherBoard: without the board check, a
// crafted column_id would file somebody's card onto a board they were not
// looking at.
func TestMoveCardRejectsAColumnOnAnotherBoard(t *testing.T) {
	db := boardTestDB(t)
	_, _, cards := seedBoard(t, db, 1)
	_, otherCol, _ := seedBoard(t, db, 1)
	r, _ := boardRouter(db)

	w := moveCard(t, r, cards[0].ID, map[string]any{"column_id": otherCol.ID})
	if w.Code != http.StatusNotFound {
		t.Fatalf("moving onto another board's column returned %d, want 404", w.Code)
	}
}

// TestMoveCardSurvivesStaleNeighbours: the neighbours the browser saw may
// have moved before the request landed. The card has to end up somewhere real
// rather than at a rank computed from rows that no longer bound it.
func TestMoveCardSurvivesStaleNeighbours(t *testing.T) {
	db := boardTestDB(t)
	_, col, cards := seedBoard(t, db, 3)
	r, _ := boardRouter(db)

	// Name the neighbours in the wrong order, which is what a stale client
	// sends after somebody else reorders the column.
	w := moveCard(t, r, cards[0].ID, map[string]any{
		"column_id": col.ID,
		"prev_id":   cards[2].ID,
		"next_id":   cards[1].ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("move: %d %s", w.Code, w.Body.String())
	}
	got := titles(t, db, col.ID)
	if len(got) != 3 {
		t.Fatalf("the column has %d cards after the move, want 3 (%q)", len(got), got)
	}
}

func TestWIPLimitRefusesAnOverfullColumn(t *testing.T) {
	db := boardTestDB(t)
	board, col, cards := seedBoard(t, db, 2)
	r, _ := boardRouter(db)

	var other models.BoardColumn
	if err := db.Where("board_id = ? AND id <> ?", board.ID, col.ID).First(&other).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&other).Update("wip_limit", 1).Error; err != nil {
		t.Fatal(err)
	}

	if w := moveCard(t, r, cards[0].ID, map[string]any{"column_id": other.ID}); w.Code != http.StatusOK {
		t.Fatalf("first move: %d %s", w.Code, w.Body.String())
	}
	w := moveCard(t, r, cards[1].ID, map[string]any{"column_id": other.ID})
	if w.Code != http.StatusConflict {
		t.Fatalf("second move into a column limited to 1 returned %d, want 409", w.Code)
	}
}

// TestReorderWithinALimitedColumnIsAllowed: a limit is about how much work is
// in progress, not about whether it can be rearranged. Counting the moving
// card against its own column would make a full column unsortable.
func TestReorderWithinALimitedColumnIsAllowed(t *testing.T) {
	db := boardTestDB(t)
	_, col, cards := seedBoard(t, db, 3)
	if err := db.Model(&models.BoardColumn{}).Where("id = ?", col.ID).Update("wip_limit", 3).Error; err != nil {
		t.Fatal(err)
	}
	r, _ := boardRouter(db)

	w := moveCard(t, r, cards[2].ID, map[string]any{
		"column_id": col.ID,
		"next_id":   cards[0].ID,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("reordering inside a full column returned %d, want 200 (%s)", w.Code, w.Body.String())
	}
	if got := titles(t, db, col.ID); got != "312" {
		t.Fatalf("order is %q, want %q", got, "312")
	}
}
`)
}
