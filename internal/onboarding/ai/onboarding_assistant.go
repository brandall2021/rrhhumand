package ai

import "fmt"

type OnboardingAssistant struct{}

func NewOnboardingAssistant() *OnboardingAssistant {
	return &OnboardingAssistant{}
}

type TaskProposal struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	ResponsibleType string `json:"responsible_type"`
	DaysOffset      int    `json:"days_offset"`
}

type OnboardingProposal struct {
	Tasks []TaskProposal `json:"tasks"`
}

func (a *OnboardingAssistant) GenerateChecklist(position, department, workMode, contractType, location, level string) *OnboardingProposal {
	proposal := &OnboardingProposal{}

	proposal.Tasks = append(proposal.Tasks, TaskProposal{
		Title: "Completar documentación personal", Description: "Subir DNI y datos personales",
		Category: "DOCUMENTACION", ResponsibleType: "EMPLOYEE", DaysOffset: 0,
	})
	proposal.Tasks = append(proposal.Tasks, TaskProposal{
		Title: "Firmar contrato laboral", Description: "Firmar contrato y anexos",
		Category: "LEGAL", ResponsibleType: "HR", DaysOffset: 0,
	})

	if workMode == "REMOTE" || workMode == "HIBRIDO" {
		proposal.Tasks = append(proposal.Tasks,
			TaskProposal{Title: "Configurar VPN", Description: "Acceso remoto seguro", Category: "IT", ResponsibleType: "IT", DaysOffset: -3},
			TaskProposal{Title: "Enviar equipo", Description: "Notebook y periféricos", Category: "EQUIPMENT", ResponsibleType: "IT", DaysOffset: -5},
		)
	}

	if department == "SISTEMAS" || department == "IT" || department == "TECNOLOGIA" {
		proposal.Tasks = append(proposal.Tasks,
			TaskProposal{Title: "Crear acceso a Git", Description: "Repositorios corporativos", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: -1},
			TaskProposal{Title: "Configurar entorno de desarrollo", Description: "Herramientas de desarrollo", Category: "IT", ResponsibleType: "IT", DaysOffset: 0},
		)
	}

	if level == "SENIOR" || level == "LEAD" || level == "MANAGER" {
		proposal.Tasks = append(proposal.Tasks,
			TaskProposal{Title: "Revisión de objetivos estratégicos", Description: "Alinear objetivos con la dirección", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 5},
		)
	}

	if level == "DIRECTOR" {
		proposal.Tasks = append(proposal.Tasks,
			TaskProposal{Title: "Presentación con dirección", Description: "Reunión con el equipo directivo", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 1},
		)
	}

	days := []int{30, 60, 90}
	descriptions := []string{"Reunión de feedback inicial", "Seguimiento de adaptación", "Evaluación de período de prueba"}
	for i, d := range days {
		proposal.Tasks = append(proposal.Tasks, TaskProposal{
			Title: descriptions[i], Description: "Seguimiento con el manager",
			Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: d,
		})
	}

	return proposal
}

func (a *OnboardingAssistant) AnswerQuestion(question string, context interface{}) string {
	answers := map[string]string{
		"¿Qué tareas faltan?":              "Revisa la sección de tareas pendientes en el dashboard de onboarding.",
		"¿Qué documentación falta?":         "Los documentos pendientes están marcados en la sección de documentación.",
		"¿Qué debe hacer IT?":              "IT debe crear cuentas, asignar equipos y configurar accesos.",
		"¿Qué debe hacer el manager?":       "El manager debe presentar al equipo, asignar tareas y realizar seguimiento.",
		"¿Cuándo termina el onboarding?":    "La fecha estimada de finalización está en el dashboard del proceso.",
		"¿Qué empleados tienen onboarding atrasado?": "Consulta el dashboard de RRHH para ver onboardings atrasados.",
	}

	if answer, ok := answers[question]; ok {
		return answer
	}
	return "No tengo información específica para esa consulta. Por favor, contacta a RRHH."
}

type ExitInterviewAnalysis struct {
	Summary         string   `json:"summary"`
	RecurringTopics []string `json:"recurring_topics"`
	Trends          []string `json:"trends"`
	Recommendations []string `json:"recommendations"`
}

func (a *OnboardingAssistant) AnalyzeExitInterview(feedback string) *ExitInterviewAnalysis {
	return &ExitInterviewAnalysis{
		Summary:         "Análisis preliminar de entrevista de salida.",
		RecurringTopics: []string{"Comunicación", "Carga laboral", "Compensación"},
		Trends:          []string{"Revisar tendencias por departamento"},
		Recommendations: []string{"Mejorar canales de comunicación interna"},
	}
}

func (a *OnboardingAssistant) GenerateOnboardingSummary(tasksTotal, tasksCompleted int, duration string) string {
	return fmt.Sprintf("Onboarding completado en %s. %d de %d tareas completadas.", duration, tasksCompleted, tasksTotal)
}


