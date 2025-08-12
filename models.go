package main

type WorkoutPlan struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Duration    string `yaml:"duration"`
	Weeks       []Week `yaml:"weeks"`
}

type Week struct {
	Number int   `yaml:"number"`
	Days   []Day `yaml:"days"`
}

type Day struct {
	Name        string    `yaml:"name"`
	Title       string    `yaml:"title"`
	Sections    []Section `yaml:"sections"`
	IsRestDay   bool      `yaml:"is_rest_day"`
	RestMessage string    `yaml:"rest_message"`
}

type Section struct {
	Name      string     `yaml:"name"`
	Duration  string     `yaml:"duration"`
	Exercises []Exercise `yaml:"exercises"`
}

type Exercise struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`                  // "timed", "reps", "rest", "info"
	Duration    int    `yaml:"duration,omitempty"`    // секунды для timed/rest
	Reps        int    `yaml:"reps,omitempty"`        // повторения для reps
	Sets        int    `yaml:"sets,omitempty"`        // подходы
	Description string `yaml:"description,omitempty"` // описание для info
	ReadTime    int    `yaml:"read_time,omitempty"`   // время на прочтение для info
}
