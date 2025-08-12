package main

import (
	"os"

	"github.com/pterm/pterm"
)

func main() {
	// Настройка pterm
	pterm.EnableStyling()
	pterm.DefaultCenter.WithCenterEachLineSeparately()

	ui := NewUI()
	ui.ShowWelcome()

	if len(os.Args) < 2 {
		ui.ShowError("Использование: go run . <путь-к-файлу-тренировки>")
		ui.ShowInfo("Пример: go run . workouts/beginner_2weeks.yaml")
		os.Exit(1)
	}

	workoutFile := os.Args[1]

	trainer := NewTrainer(ui)
	if err := trainer.LoadWorkout(workoutFile); err != nil {
		ui.ShowError("Ошибка загрузки тренировки: " + err.Error())
		os.Exit(1)
	}

	trainer.StartWorkout()
}
