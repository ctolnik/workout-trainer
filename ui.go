package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type UI struct {
	progressArea *pterm.AreaPrinter
	statsArea    *pterm.AreaPrinter
}

func NewUI() *UI {
	return &UI{}
}

//	func (ui *UI) Clear() {
//		pterm.Clear()
//	}
func (ui *UI) Clear() {
	// Универсальная очистка экрана для всех ОС
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (ui *UI) ShowWelcome() {
	ui.Clear()

	// Анимированный заголовок
	title, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("WORKOUT", pterm.NewStyle(pterm.FgLightBlue)),
		putils.LettersFromStringWithStyle("TRAINER", pterm.NewStyle(pterm.FgLightMagenta)),
	).Srender()

	pterm.DefaultCenter.Print(title)

	// Красивый бокс с информацией
	pterm.DefaultBox.WithTitle("🏋️‍♀️ Персональный Тренер").WithTitleTopCenter().WithBoxStyle(pterm.NewStyle(pterm.FgLightCyan)).Print(
		"💪 Интерактивные тренировки\n" +
			"⏱️  Таймеры и счетчики\n" +
			"📊 Визуализация прогресса\n" +
			"🎯 Персонализированный подход")

	pterm.Println()
}

func (ui *UI) ShowWorkoutInfo(workout WorkoutPlan) {
	ui.Clear()

	// Главный заголовок
	pterm.DefaultHeader.WithFullWidth().WithMargin(2).WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).Print("📋 ИНФОРМАЦИЯ О ТРЕНИРОВКЕ")

	// Панель с информацией - исправлено
	panels := pterm.Panels{
		{
			{Data: pterm.Sprintf("%s\n\n%s",
				pterm.LightMagenta("📖 Описание:"),
				workout.Description)},
			{Data: pterm.Sprintf("%s\n%s",
				pterm.LightCyan("⏱️ Длительность:"),
				workout.Duration)},
		},
	}

	pterm.DefaultPanel.WithPanels(panels).WithPadding(2).Render()

	// Список недель
	weekList := []pterm.BulletListItem{}
	for _, week := range workout.Weeks {
		dayItems := []pterm.BulletListItem{}
		for _, day := range week.Days {
			emoji := "🏃‍♂️"
			if day.IsRestDay {
				emoji = "😌"
			}
			dayItems = append(dayItems, pterm.BulletListItem{
				Level: 1,
				Text:  fmt.Sprintf("%s %s - %s", emoji, day.Name, day.Title),
			})
		}
		weekList = append(weekList, pterm.BulletListItem{
			Level:       0,
			Text:        fmt.Sprintf("📅 Неделя %d", week.Number),
			TextStyle:   pterm.NewStyle(pterm.FgLightBlue, pterm.Bold),
			BulletStyle: pterm.NewStyle(pterm.FgLightBlue),
		})
		weekList = append(weekList, dayItems...)
	}

	pterm.DefaultBulletList.WithItems(weekList).Render()
}

func (ui *UI) ShowDayStart(day Day) {
	ui.Clear()

	if day.IsRestDay {
		// День отдыха
		pterm.DefaultBox.WithTitle("😌 ДЕНЬ ОТДЫХА").WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgLightGreen)).Print(
			fmt.Sprintf("%s - %s\n\n🛌 %s", day.Name, day.Title, day.RestMessage))
		return
	}

	// Обычный тренировочный день
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgRed)).Print(
		fmt.Sprintf("🔥 %s - %s", day.Name, day.Title))

	// Показываем структуру дня
	sectionList := []pterm.BulletListItem{}
	for i, section := range day.Sections {
		sectionList = append(sectionList, pterm.BulletListItem{
			Level: 0,
			Text:  fmt.Sprintf("%d. %s (%s)", i+1, section.Name, section.Duration),
		})

		for _, exercise := range section.Exercises {
			emoji := ui.getExerciseEmoji(exercise.Type)
			details := ui.getExerciseDetails(exercise)
			sectionList = append(sectionList, pterm.BulletListItem{
				Level: 1,
				Text:  fmt.Sprintf("%s %s %s", emoji, exercise.Name, details),
			})
		}
	}

	pterm.DefaultBulletList.WithItems(sectionList).Render()
}

func (ui *UI) ShowSectionStart(section Section) {
	ui.Clear()

	pterm.DefaultBox.WithTitle(fmt.Sprintf("📝 %s", section.Name)).WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgLightYellow)).Print(
		fmt.Sprintf("⏱️ Примерная длительность: %s\n\n🏃‍♂️ Готовимся к выполнению упражнений!", section.Duration))
}

func (ui *UI) ShowExerciseStart(exercise Exercise) {
	ui.Clear()

	emoji := ui.getExerciseEmoji(exercise.Type)
	color := ui.getExerciseColor(exercise.Type)

	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(color)).Print(
		fmt.Sprintf("%s %s", emoji, exercise.Name))

	if exercise.Description != "" {
		pterm.DefaultBox.WithTitle("📖 Описание упражнения").WithTitleTopLeft().Print(exercise.Description)
	}

	details := ui.getExerciseDetails(exercise)
	if details != "" {
		pterm.Info.Print("🎯 " + details)
	}
}

func (ui *UI) ShowTimedExercise(exercise Exercise, currentSet, totalSets int) {
	ui.Clear()

	emoji := ui.getExerciseEmoji(exercise.Type)
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).Print(
		fmt.Sprintf("%s %s", emoji, exercise.Name))

	if totalSets > 1 {
		pterm.Info.Printf("🔄 Подход %d из %d\n", currentSet, totalSets)
	}

	if exercise.Description != "" {
		pterm.DefaultBox.WithTitle("📋 Техника выполнения").WithTitleTopLeft().Print(exercise.Description)
	}
}

func (ui *UI) ShowTimer(seconds int, activity string, exercise Exercise) {
	totalSeconds := seconds

	for i := seconds; i > 0; i-- {
		ui.Clear()

		// Заголовок
		emoji := ui.getExerciseEmoji(exercise.Type)
		pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).Print(
			fmt.Sprintf("%s %s", emoji, activity))

		// Время
		minutes := i / 60
		secs := i % 60
		var timeStr string
		if minutes > 0 {
			timeStr = fmt.Sprintf("%02d:%02d", minutes, secs)
		} else {
			timeStr = fmt.Sprintf("%02d сек", secs)
		}

		// Большие цифры времени
		bigTime, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle(timeStr, pterm.NewStyle(pterm.FgLightGreen)),
		).Srender()

		pterm.DefaultCenter.Print(bigTime)

		// Прогресс-бар
		progress := float64(totalSeconds-i) / float64(totalSeconds) * 100
		progressBar, _ := pterm.DefaultProgressbar.WithTotal(100).WithCurrent(int(progress)).Start()
		progressBar.Stop()

		// Статистика
		ui.showTimeStats(i, totalSeconds, activity)

		time.Sleep(1 * time.Second)
	}

	ui.ShowSuccess("✅ " + activity + " завершен!")
}

func (ui *UI) ShowRepsExercise(exercise Exercise, currentSet, totalSets int) {
	ui.Clear()

	emoji := ui.getExerciseEmoji(exercise.Type)
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightMagenta)).Print(
		fmt.Sprintf("%s %s", emoji, exercise.Name))

	if totalSets > 1 {
		// Прогресс подходов
		setProgress, _ := pterm.DefaultProgressbar.WithTotal(totalSets).WithCurrent(currentSet - 1).Start()
		setProgress.UpdateTitle(fmt.Sprintf("Подход %d из %d", currentSet, totalSets))
		setProgress.Stop()
	}

	// Большие цифры для повторений
	repsText := fmt.Sprintf("%d", exercise.Reps)
	bigReps, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle(repsText, pterm.NewStyle(pterm.FgLightYellow)),
	).Srender()

	pterm.DefaultCenter.Print("🎯 ВЫПОЛНИТЕ ПОВТОРЕНИЙ:")
	pterm.DefaultCenter.Print(bigReps)

	if exercise.Description != "" {
		pterm.DefaultBox.WithTitle("📋 Техника выполнения").WithTitleTopLeft().Print(exercise.Description)
	}
}

func (ui *UI) ShowRestTimer(seconds int, restType string) {
	totalSeconds := seconds

	for i := seconds; i > 0; i-- {
		ui.Clear()

		pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightGreen)).Print(
			fmt.Sprintf("😌 %s", restType))

		// Время отдыха
		minutes := i / 60
		secs := i % 60
		timeStr := fmt.Sprintf("%02d:%02d", minutes, secs)

		bigTime, _ := pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle(timeStr, pterm.NewStyle(pterm.FgLightCyan)),
		).Srender()

		pterm.DefaultCenter.Print(bigTime)

		// Круговой прогресс для отдыха
		progress := float64(totalSeconds-i) / float64(totalSeconds) * 100

		// Создаем красивый прогресс-бар отдыха
		spinner, _ := pterm.DefaultSpinner.WithText(fmt.Sprintf("💤 Восстанавливаемся... %.1f%%", progress)).Start()
		time.Sleep(500 * time.Millisecond)
		spinner.Stop()

		// Мотивационные сообщения
		ui.showRestMotivation(i, totalSeconds)

		time.Sleep(500 * time.Millisecond)
	}

	ui.ShowSuccess("⚡ Отдохнули! Готовы продолжать!")
}

func (ui *UI) ShowWorkoutComplete() {
	ui.Clear()

	// Анимированное поздравление
	congratsText, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("ОТЛИЧНО", pterm.NewStyle(pterm.FgLightGreen)),
	).Srender()

	pterm.DefaultCenter.Print(congratsText)

	// Коробка с поздравлениями
	pterm.DefaultBox.WithTitle("🎉 ПОЗДРАВЛЯЕМ!").WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgLightGreen)).Print(
		"🏆 Вы завершили программу!\n" +
			"💪 Ваше тело стало сильнее!\n" +
			"🧠 Ваш дух стал крепче!\n" +
			"❤️  Ваше здоровье улучшилось!\n\n" +
			"🚀 Продолжайте в том же духе!")

	// Конфетти эффект
	for i := 0; i < 3; i++ {
		pterm.Print("🎊 ")
		time.Sleep(200 * time.Millisecond)
	}
	pterm.Println()
}

func (ui *UI) ShowSuccess(message string) {
	pterm.Success.Println(message)
	time.Sleep(2 * time.Second)
}

func (ui *UI) ShowError(message string) {
	pterm.Error.Println(message)
}

func (ui *UI) ShowInfo(message string) {
	pterm.Info.Println(message)
}

func (ui *UI) WaitForUser(message string) {
	pterm.DefaultInteractiveConfirm.WithDefaultText(message).WithDefaultValue(true).Show()
}

func (ui *UI) getExerciseEmoji(exerciseType string) string {
	switch exerciseType {
	case "timed":
		return "⏱️"
	case "reps":
		return "🔢"
	case "rest":
		return "😌"
	case "info":
		return "📖"
	default:
		return "🏃‍♂️"
	}
}

func (ui *UI) getExerciseColor(exerciseType string) pterm.Color {
	switch exerciseType {
	case "timed":
		return pterm.BgLightBlue
	case "reps":
		return pterm.BgLightMagenta
	case "rest":
		return pterm.BgLightGreen
	case "info":
		return pterm.BgLightYellow
	default:
		return pterm.BgLightCyan
	}
}

func (ui *UI) getExerciseDetails(exercise Exercise) string {
	switch exercise.Type {
	case "timed":
		sets := exercise.Sets
		if sets <= 1 {
			return fmt.Sprintf("(%d сек)", exercise.Duration)
		}
		return fmt.Sprintf("(%d подходов по %d сек)", sets, exercise.Duration)
	case "reps":
		sets := exercise.Sets
		if sets <= 1 {
			return fmt.Sprintf("(%d повторений)", exercise.Reps)
		}
		return fmt.Sprintf("(%d подходов по %d повторений)", sets, exercise.Reps)
	case "rest":
		return fmt.Sprintf("(%d сек отдыха)", exercise.Duration)
	case "info":
		if exercise.ReadTime > 0 {
			return fmt.Sprintf("(%d сек на изучение)", exercise.ReadTime)
		}
		return "(информация)"
	default:
		return ""
	}
}

func (ui *UI) showTimeStats(remaining, total int, activity string) {
	elapsed := total - remaining
	progress := float64(elapsed) / float64(total) * 100

	table := pterm.TableData{
		{"📊 Статистика", "Значение"},
		{"⏱️ Прошло времени", fmt.Sprintf("%d сек", elapsed)},
		{"⏳ Осталось времени", fmt.Sprintf("%d сек", remaining)},
		{"📈 Прогресс", fmt.Sprintf("%.1f%%", progress)},
		{"🎯 Упражнение", activity},
	}

	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}

func (ui *UI) showRestMotivation(remaining, total int) {
	messages := []string{
		"💪 Мышцы восстанавливаются!",
		"🫁 Дышите глубоко и ровно",
		"💧 Не забудьте пить воду",
		"🧘‍♀️ Расслабьтесь и отдохните",
		"⚡ Готовимся к следующему упражнению",
		"🎯 Вы отлично справляетесь!",
	}

	// Показываем разные сообщения в зависимости от оставшегося времени
	messageIndex := (total - remaining) % len(messages)
	pterm.Info.Println(messages[messageIndex])
}
