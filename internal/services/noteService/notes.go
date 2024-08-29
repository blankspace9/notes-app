package noteservice

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
)

const (
	emptyValue = 0
)

type NoteService struct {
	log          *slog.Logger
	notesManager NotesManager
	spellChecker SpellChecker
}

type NotesManager interface {
	SaveNote(ctx context.Context, note models.Note) (noteID int64, err error)
	GetNotesByUserId(ctx context.Context, userID int64) ([]models.Note, error)
	GetNotesPageByUserId(ctx context.Context, userID int64, page, limit int) ([]models.Note, error)
}

type SpellChecker interface {
	CheckSpelling(text string) ([]models.SpellError, error)
}

func New(log *slog.Logger, notesManager NotesManager, spellChecker SpellChecker) *NoteService {
	return &NoteService{
		log:          log,
		notesManager: notesManager,
		spellChecker: spellChecker,
	}
}

func (ns *NoteService) CreateNote(ctx context.Context, note string, userID int64) (int64, []models.SpellError, error) {
	const op = "services.NoteService.CreateNote"

	log := ns.log.With(slog.String("op", op))

	log.Info("attempting to create note")

	spellingErrors, err := ns.spellChecker.CheckSpelling(note)
	if err != nil {
		log.Error("failed to check spelling errors", sl.Err(err))

		return 0, nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := ns.notesManager.SaveNote(ctx, models.Note{
		Note:      note,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Error("failed to save note", sl.Err(err))

		return 0, nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, spellingErrors, nil
}

func (ns *NoteService) GetNotes(ctx context.Context, userID int64, page, limit int) ([]models.Note, error) {
	const op = "services.NoteService.GetNotes"

	log := ns.log.With(slog.String("op", op))

	log.Info("attempting to get notes")

	var notes []models.Note
	var err error
	if page == emptyValue || limit == emptyValue {
		notes, err = ns.notesManager.GetNotesByUserId(ctx, userID)
	} else {
		notes, err = ns.notesManager.GetNotesPageByUserId(ctx, userID, page, limit)
	}
	if err != nil {
		log.Error("failed to get notes", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}
