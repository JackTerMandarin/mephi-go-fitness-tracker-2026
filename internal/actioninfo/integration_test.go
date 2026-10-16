package actioninfo_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/actioninfo"
	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/trainings"
)

var (
	_ actioninfo.DataParser = (*trainings.Training)(nil)
	_ actioninfo.DataParser = (*daysteps.DaySteps)(nil)
)

func TestParseKeepsPreviousStateOnError(t *testing.T) {
	person := personaldata.Personal{Name: "Иван", Weight: 75, Height: 1.75}
	training := trainings.Training{Steps: 1000, TrainingType: "Бег", Duration: time.Hour, Personal: person}
	walk := daysteps.DaySteps{Steps: 1000, Duration: time.Hour, Personal: person}
	for _, input := range []string{"2000,Ходьба,invalid", "2000,Ходьба,0m", "99999999999999999999999,Бег,1h"} {
		before := training
		if err := training.Parse(input); err == nil {
			t.Fatalf("Parse(%q) should fail", input)
		}
		if training != before {
			t.Fatalf("failed parsing changed training: %+v", training)
		}
	}
	for _, input := range []string{"2000,invalid", "2000,0m", "99999999999999999999999,1h"} {
		before := walk
		if err := walk.Parse(input); err == nil {
			t.Fatalf("Parse(%q) should fail", input)
		}
		if walk != before {
			t.Fatalf("failed parsing changed walk: %+v", walk)
		}
	}
	if err := training.Parse("2000,Ходьба,30m"); err != nil || training.Personal != person {
		t.Fatalf("successful parsing must preserve personal data: %+v, %v", training, err)
	}
	if err := walk.Parse("2000,30m"); err != nil || walk.Personal != person {
		t.Fatalf("successful parsing must preserve personal data: %+v, %v", walk, err)
	}
}

func TestInfoContinuesAfterParseAndActionErrors(t *testing.T) {
	var logs bytes.Buffer
	previousLog := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(previousLog) })

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	previousStdout := os.Stdout
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = previousStdout
		reader.Close()
		writer.Close()
	})

	training := &trainings.Training{Personal: personaldata.Personal{Weight: 75, Height: 1.75}}
	actioninfo.Info([]string{"6000,Бег,1h", "bad data", "1000,Плавание,1h", "3000,Ходьба,30m"}, training)
	writer.Close()
	os.Stdout = previousStdout
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	text := string(output)
	if strings.Count(text, "Тип тренировки:") != 2 || !strings.Contains(text, "354.38") || !strings.Contains(text, "88.59") {
		t.Fatalf("valid reports missing or stale report repeated: %q", text)
	}
	if !strings.Contains(logs.String(), "bad data") || !strings.Contains(logs.String(), "неизвестный тип тренировки") {
		t.Fatalf("errors not logged: %q", logs.String())
	}
}

func TestTrainingActionInfoRejectsInvalidData(t *testing.T) {
	for _, kind := range []string{"Бег", "Ходьба"} {
		training := trainings.Training{Steps: 1000, TrainingType: kind, Duration: time.Hour}
		if report, err := training.ActionInfo(); err == nil || report != "" {
			t.Fatalf("invalid personal data: got (%q, %v)", report, err)
		}
	}
	training := trainings.Training{TrainingType: "Плавание"}
	if report, err := training.ActionInfo(); report != "" || err == nil || err.Error() != "неизвестный тип тренировки" {
		t.Fatalf("unknown training: got (%q, %v)", report, err)
	}
}
