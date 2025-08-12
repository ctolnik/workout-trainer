package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Trainer struct {
	workout WorkoutPlan
	ui      *UI
}

func NewTrainer(ui *UI) *Trainer {
	return &Trainer{
		ui: ui,
	}
}

func (t *Trainer) LoadWorkout(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &t.workout)
}

func (t *Trainer) StartWorkout() {
	t.ui.ShowWorkoutInfo(t.workout)
	t.ui.WaitForUser("Готовы начать тренировку?")

	for _, week := range t.workout.Weeks {
		t.startWeek(week)
	}

	t.ui.ShowWorkoutComplete()
}

func (t *Trainer) startWeek(week Week) {
	t.ui.ShowInfo(fmt.Sprintf("📅 Начинаем неделю %d", week.Number))
	t.ui.WaitForUser("Готовы начать неделю?")

	for _, day := range week.Days {
		t.startDay(day)
	}

	t.ui.ShowSuccess(fmt.Sprintf("🎉 Неделя %d завершена!", week.Number))
}

func (t *Trainer) startDay(day Day) {
	t.ui.ShowDayStart(day)

	if day.IsRestDay {
		t.ui.WaitForUser("Нажмите Enter для продолжения...")
		return
	}

	t.ui.WaitForUser("Готовы начать тренировку?")

	for _, section := range day.Sections {
		t.startSection(section)
	}

	t.ui.ShowSuccess(fmt.Sprintf("✅ Тренировка '%s' завершена!", day.Title))
}

func (t *Trainer) startSection(section Section) {
	t.ui.ShowSectionStart(section)
	t.ui.WaitForUser("Переходим к упражнениям")

	for _, exercise := range section.Exercises {
		t.performExercise(exercise)
	}
}

func (t *Trainer) performExercise(exercise Exercise) {
	t.ui.ShowExerciseStart(exercise)

	switch exercise.Type {
	case "info":
		t.performInfo(exercise)
	case "timed":
		t.performTimedExercise(exercise)
	case "reps":
		t.performRepsExercise(exercise)
	case "rest":
		t.performRest(exercise)
	}
}

func (t *Trainer) performInfo(exercise Exercise) {
	if exercise.ReadTime > 0 {
		t.ui.ShowTimer(exercise.ReadTime, "Изучаем технику", exercise)
	}
	t.ui.WaitForUser("Понятно? Переходим дальше")
}

func (t *Trainer) performTimedExercise(exercise Exercise) {
	sets := exercise.Sets
	if sets == 0 {
		sets = 1
	}

	for set := 1; set <= sets; set++ {
		t.ui.ShowTimedExercise(exercise, set, sets)
		t.ui.WaitForUser("Готовы начать?")

		t.ui.ShowTimer(exercise.Duration, exercise.Name, exercise)

		if set < sets {
			t.ui.ShowSuccess("✅ Подход завершен!")
			t.ui.ShowRestTimer(60, "Отдых между подходами")
		}
	}

	t.ui.ShowSuccess("🎉 Упражнение завершено!")
}

func (t *Trainer) performRepsExercise(exercise Exercise) {
	sets := exercise.Sets
	if sets == 0 {
		sets = 1
	}

	for set := 1; set <= sets; set++ {
		t.ui.ShowRepsExercise(exercise, set, sets)
		t.ui.WaitForUser("Выполните упражнение и нажмите Enter")

		if set < sets {
			t.ui.ShowSuccess("✅ Подход завершен!")
			t.ui.ShowRestTimer(90, "Отдых между подходами")
		}
	}

	t.ui.ShowSuccess("🎉 Упражнение завершено!")
}

func (t *Trainer) performRest(exercise Exercise) {
	t.ui.ShowRestTimer(exercise.Duration, exercise.Name)
}
