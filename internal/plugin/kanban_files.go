package plugin

import "strings"

// goSource renders an emitted Go template: ~ becomes a backtick, so struct
// tags can be written inline, and {{MODULE}} becomes the installing project's
// module path.
func goSource(ctx Context, src string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// kanbanModelsGo emits internal/models/board.go.
func kanbanModelsGo(ctx Context) string {
	return goSource(ctx, `package models

import (
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// Board is one Kanban board: a named set of columns holding cards.
//
// Boards are visible to every signed-in user. CreatedBy is for display, not
// access control. If boards need to be scoped per tenant, install the
// multitenant plugin, which scopes every model including these.
type Board struct {
	ID          string ~gorm:"primarykey;size:36" json:"id"~
	Name        string ~gorm:"size:160;not null" json:"name" binding:"required"~
	Description string ~gorm:"type:text" json:"description"~

	CreatedBy string ~gorm:"size:36;index" json:"created_by"~

	Columns []BoardColumn ~gorm:"foreignKey:BoardID" json:"columns,omitempty"~

	CreatedAt time.Time      ~json:"created_at"~
	UpdatedAt time.Time      ~json:"updated_at"~
	DeletedAt gorm.DeletedAt ~gorm:"index" json:"-"~
}

func (b *Board) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = ids.New()
	}
	return nil
}

// BoardColumn is one vertical lane: "To do", "In progress", "Done".
type BoardColumn struct {
	ID      string ~gorm:"primarykey;size:36" json:"id"~
	BoardID string ~gorm:"size:36;index;not null" json:"board_id"~
	Name    string ~gorm:"size:120;not null" json:"name" binding:"required"~

	// Color names a token in the admin palette rather than holding a hex
	// value, so a board follows the theme instead of fighting it. The UI maps
	// unknown names to the default.
	Color string ~gorm:"size:32" json:"color"~

	// Rank orders the columns on the board. A string, not an integer: see
	// internal/kanban/rank.go for why.
	Rank string ~gorm:"size:255;not null;index:idx_board_columns_board_rank,priority:2" json:"rank"~

	// WIPLimit caps how many cards this column may hold. Zero means no limit.
	// A limit is a signal that the team has agreed on, so exceeding it is
	// refused rather than warned about.
	//
	// The column name is pinned. GORM's naming strategy turns WIPLimit into
	// w_ip_limit, because it splits a run of capitals before the last one.
	// Reading the struct still works, so this is invisible until something
	// writes by column name: Updates(map[string]any{"wip_limit": ...}) then
	// fails against a column that does not exist, in the handler rather than
	// at migration time.
	WIPLimit int ~gorm:"column:wip_limit;default:0" json:"wip_limit"~

	// The index the board query actually uses: look up one column's cards in
	// order. Declared on BoardID with priority 1 so it pairs with Rank above.
	Cards []BoardCard ~gorm:"foreignKey:ColumnID" json:"cards,omitempty"~

	CreatedAt time.Time      ~json:"created_at"~
	UpdatedAt time.Time      ~json:"updated_at"~
	DeletedAt gorm.DeletedAt ~gorm:"index" json:"-"~
}

func (c *BoardColumn) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = ids.New()
	}
	return nil
}

// BoardCard is one card.
//
// BoardID is denormalised from the column deliberately. Every authorisation
// check and every board-wide query asks "which board is this card on", and
// without it that is a join on every request to answer a question that never
// changes: a card moves between columns of one board, never between boards.
type BoardCard struct {
	ID       string ~gorm:"primarykey;size:36" json:"id"~
	BoardID  string ~gorm:"size:36;index;not null" json:"board_id"~
	ColumnID string ~gorm:"size:36;not null;index:idx_board_cards_column_rank,priority:1" json:"column_id"~

	Title       string ~gorm:"size:240;not null" json:"title" binding:"required"~
	Description string ~gorm:"type:text" json:"description"~

	// Rank orders cards within their column. The read is always
	// "WHERE column_id = ? ORDER BY rank, id", which is exactly the composite
	// index these two tags declare.
	Rank string ~gorm:"size:255;not null;index:idx_board_cards_column_rank,priority:2" json:"rank"~

	AssigneeID *string    ~gorm:"size:36;index" json:"assignee_id,omitempty"~
	DueDate    *time.Time ~json:"due_date,omitempty"~

	CreatedBy string ~gorm:"size:36;index" json:"created_by"~

	CreatedAt time.Time      ~json:"created_at"~
	UpdatedAt time.Time      ~json:"updated_at"~
	DeletedAt gorm.DeletedAt ~gorm:"index" json:"-"~
}

func (c *BoardCard) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = ids.New()
	}
	return nil
}
`)
}

// kanbanServiceGo emits internal/services/board.go.
func kanbanServiceGo(ctx Context) string {
	return goSource(ctx, `package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"{{MODULE}}/internal/kanban"
	"{{MODULE}}/internal/models"
)

// ErrWIPLimit is returned when a move would push a column past the
// work-in-progress limit its team agreed on.
var ErrWIPLimit = errors.New("column is at its work-in-progress limit")

// BoardService holds the board logic the handlers call. The interesting part
// is RankForCard: everything else is bookkeeping.
type BoardService struct {
	DB *gorm.DB
}

func NewBoardService(db *gorm.DB) *BoardService {
	return &BoardService{DB: db}
}

// RankForCard returns the rank a card takes when dropped into columnID
// between prevID (the card above the drop) and nextID (the card below it).
//
// The client sends neighbour ids rather than an index. An index is stale the
// instant somebody else moves a card, and a stale index silently puts the
// card in the wrong place. Neighbour ids survive that: if one of them has
// moved away, this notices and falls back to a position that exists rather
// than computing a rank from a row that is no longer there.
//
// Either id may be empty: no prev means the top of the column, no next means
// the bottom, neither means the column is empty or the caller is appending.
func (s *BoardService) RankForCard(columnID, prevID, nextID, excludeID string) (string, error) {
	prev := s.cardRank(columnID, prevID, excludeID)
	next := s.cardRank(columnID, nextID, excludeID)

	if prev != "" && next != "" && prev >= next {
		// The neighbours the client named are no longer adjacent, or no
		// longer in that order. Appending is not where the user aimed, but it
		// is a position that exists, and the board they see after the refetch
		// is the truth. Guessing an interior position from stale data would
		// put the card somewhere nobody chose.
		prev, next = s.lastCardRank(columnID, excludeID), ""
	}

	if prev == "" && next == "" {
		// No neighbours named: append to the end of the column.
		prev = s.lastCardRank(columnID, excludeID)
	}

	rank, err := kanban.Between(prev, next)
	if err != nil {
		return "", fmt.Errorf("computing card rank: %w", err)
	}
	return rank, nil
}

// RankForColumn is the same for columns, which are ordered across a board.
func (s *BoardService) RankForColumn(boardID, prevID, nextID, excludeID string) (string, error) {
	prev := s.columnRank(boardID, prevID, excludeID)
	next := s.columnRank(boardID, nextID, excludeID)

	if prev != "" && next != "" && prev >= next {
		prev, next = s.lastColumnRank(boardID, excludeID), ""
	}
	if prev == "" && next == "" {
		prev = s.lastColumnRank(boardID, excludeID)
	}

	rank, err := kanban.Between(prev, next)
	if err != nil {
		return "", fmt.Errorf("computing column rank: %w", err)
	}
	return rank, nil
}

// CheckWIP reports whether moving a card into columnID would exceed the
// column's limit. The moving card is excluded from the count, so dropping a
// card back into the column it came from is never refused.
func (s *BoardService) CheckWIP(columnID, movingCardID string) error {
	var col models.BoardColumn
	if err := s.DB.First(&col, "id = ?", columnID).Error; err != nil {
		return err
	}
	if col.WIPLimit <= 0 {
		return nil
	}

	q := s.DB.Model(&models.BoardCard{}).Where("column_id = ?", columnID)
	if movingCardID != "" {
		q = q.Where("id <> ?", movingCardID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return err
	}
	if int(n) >= col.WIPLimit {
		return fmt.Errorf("%w: %s holds %d of %d", ErrWIPLimit, col.Name, n, col.WIPLimit)
	}
	return nil
}

// DefaultColumns is what a new board starts with, so that creating one lands
// on a usable board rather than an empty screen with a single "add a column"
// button.
func (s *BoardService) DefaultColumns(boardID string) ([]models.BoardColumn, error) {
	names := []string{"To do", "In progress", "Done"}
	colors := []string{"info", "warning", "success"}

	ranks, err := kanban.NAfter("", "", len(names))
	if err != nil {
		return nil, err
	}

	cols := make([]models.BoardColumn, 0, len(names))
	for i, name := range names {
		cols = append(cols, models.BoardColumn{
			BoardID: boardID,
			Name:    name,
			Color:   colors[i],
			Rank:    ranks[i],
		})
	}
	return cols, nil
}

// LoadBoard returns a board with its columns and their cards, each in rank
// order.
//
// Two queries, not one per column: a board with twelve columns would
// otherwise be thirteen round trips, which is the N+1 that makes a board feel
// slow exactly when it has enough on it to be useful.
func (s *BoardService) LoadBoard(id string) (*models.Board, error) {
	var board models.Board
	if err := s.DB.First(&board, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := s.DB.
		Where("board_id = ?", id).
		Order("rank asc, id asc").
		Find(&board.Columns).Error; err != nil {
		return nil, err
	}
	if len(board.Columns) == 0 {
		return &board, nil
	}

	colIDs := make([]string, 0, len(board.Columns))
	for _, c := range board.Columns {
		colIDs = append(colIDs, c.ID)
	}

	var cards []models.BoardCard
	if err := s.DB.
		Where("column_id IN ?", colIDs).
		// Ties are possible: two people dropping into the same slot at the
		// same moment compute the same rank. Breaking the tie by id makes the
		// order stable and identical for both of them.
		Order("rank asc, id asc").
		Find(&cards).Error; err != nil {
		return nil, err
	}

	byColumn := map[string][]models.BoardCard{}
	for _, card := range cards {
		byColumn[card.ColumnID] = append(byColumn[card.ColumnID], card)
	}
	for i := range board.Columns {
		board.Columns[i].Cards = byColumn[board.Columns[i].ID]
	}
	return &board, nil
}

// DeleteBoard removes a board and everything on it.
//
// GORM's soft delete does not cascade, and a database-level ON DELETE CASCADE
// never fires for an UPDATE that only sets deleted_at. Without this, deleting
// a board leaves its columns and cards behind: invisible, still counted by
// every aggregate, and restored into nothing if the board is undeleted.
func (s *BoardService) DeleteBoard(id string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("board_id = ?", id).Delete(&models.BoardCard{}).Error; err != nil {
			return err
		}
		if err := tx.Where("board_id = ?", id).Delete(&models.BoardColumn{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Board{}, "id = ?", id).Error
	})
}

// DeleteColumn removes a column and its cards.
func (s *BoardService) DeleteColumn(id string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("column_id = ?", id).Delete(&models.BoardCard{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.BoardColumn{}, "id = ?", id).Error
	})
}

func (s *BoardService) cardRank(columnID, cardID, excludeID string) string {
	if cardID == "" || cardID == excludeID {
		return ""
	}
	var card models.BoardCard
	if err := s.DB.Select("rank").
		Where("id = ? AND column_id = ?", cardID, columnID).
		First(&card).Error; err != nil {
		// Not an error: the neighbour has been moved or deleted since the
		// drag began. The caller treats a missing neighbour as "no bound".
		return ""
	}
	return card.Rank
}

func (s *BoardService) lastCardRank(columnID, excludeID string) string {
	var card models.BoardCard
	q := s.DB.Select("rank").Where("column_id = ?", columnID)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Order("rank desc, id desc").First(&card).Error; err != nil {
		return ""
	}
	return card.Rank
}

func (s *BoardService) columnRank(boardID, columnID, excludeID string) string {
	if columnID == "" || columnID == excludeID {
		return ""
	}
	var col models.BoardColumn
	if err := s.DB.Select("rank").
		Where("id = ? AND board_id = ?", columnID, boardID).
		First(&col).Error; err != nil {
		return ""
	}
	return col.Rank
}

func (s *BoardService) lastColumnRank(boardID, excludeID string) string {
	var col models.BoardColumn
	q := s.DB.Select("rank").Where("board_id = ?", boardID)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Order("rank desc, id desc").First(&col).Error; err != nil {
		return ""
	}
	return col.Rank
}
`)
}
