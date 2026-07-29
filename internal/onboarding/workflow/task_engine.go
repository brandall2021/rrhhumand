package workflow

import (
	"context"
	"time"

	"github.com/rrhhumand/api/internal/onboarding/domain"
	"github.com/rrhhumand/api/internal/onboarding/integration"
)

type TaskEngine struct {
	empSvc    integration.EmployeeService
	assetSvc  integration.AssetService
	trainSvc  integration.TrainingService
	accessSvc integration.AccessProvisioningService
	signSvc   integration.SignatureService
}

func NewTaskEngine(
	empSvc integration.EmployeeService,
	assetSvc integration.AssetService,
	trainSvc integration.TrainingService,
	accessSvc integration.AccessProvisioningService,
	signSvc integration.SignatureService,
) *TaskEngine {
	return &TaskEngine{
		empSvc:    empSvc,
		assetSvc:  assetSvc,
		trainSvc:  trainSvc,
		accessSvc: accessSvc,
		signSvc:   signSvc,
	}
}

type TemplateTaskDef struct {
	Title           string
	Description     string
	TaskType        domain.TaskType
	AssignedRole    domain.ResponsibleType
	Required        bool
	DaysOffset      int
}

var DefaultOnboardingTasks = []TemplateTaskDef{
	{Title: "Completar datos personales", Description: "Completar formulario de datos personales en el sistema", TaskType: domain.TaskInfo, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 0},
	{Title: "Firma de contrato", Description: "Firmar el contrato laboral", TaskType: domain.TaskSignature, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 0},
	{Title: "Crear email corporativo", Description: "Configurar cuenta de email corporativo", TaskType: domain.TaskAccount, AssignedRole: domain.ResponsibleIT, Required: true, DaysOffset: -2},
	{Title: "Crear usuario en sistemas", Description: "Crear usuario en los sistemas internos", TaskType: domain.TaskAccount, AssignedRole: domain.ResponsibleIT, Required: true, DaysOffset: -2},
	{Title: "Configurar VPN", Description: "Configurar acceso VPN", TaskType: domain.TaskSystem, AssignedRole: domain.ResponsibleIT, Required: false, DaysOffset: -1},
	{Title: "Asignar notebook", Description: "Asignar equipo de trabajo", TaskType: domain.TaskAsset, AssignedRole: domain.ResponsibleIT, Required: true, DaysOffset: -1},
	{Title: "Presentación al equipo", Description: "Presentar al nuevo empleado con el equipo", TaskType: domain.TaskMeeting, AssignedRole: domain.ResponsibleManager, Required: true, DaysOffset: 1},
	{Title: "Capacitación inicial", Description: "Completar capacitación de inducción", TaskType: domain.TaskTraining, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 7},
	{Title: "Capacitación de seguridad", Description: "Completar capacitación de seguridad informática", TaskType: domain.TaskTraining, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 7},
	{Title: "Revisar políticas internas", Description: "Leer y aceptar las políticas internas", TaskType: domain.TaskChecklist, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 3},
	{Title: "Feedback 30 días", Description: "Reunión de feedback con el manager", TaskType: domain.TaskMeeting, AssignedRole: domain.ResponsibleManager, Required: true, DaysOffset: 30},
	{Title: "Seguimiento 60 días", Description: "Reunión de seguimiento con el manager", TaskType: domain.TaskMeeting, AssignedRole: domain.ResponsibleManager, Required: false, DaysOffset: 60},
	{Title: "Evaluación período de prueba", Description: "Evaluación de desempeño inicial", TaskType: domain.TaskMeeting, AssignedRole: domain.ResponsibleManager, Required: true, DaysOffset: 90},
}

func (e *TaskEngine) GenerateOnboardingTasks(ctx context.Context, companyID, onboardingID, employeeID, startDate string, tasks []TemplateTaskDef) ([]domain.OnboardingTaskAssignment, error) {
	var assignments []domain.OnboardingTaskAssignment

	for i, t := range tasks {
		dueDate := calcDueDate(startDate, t.DaysOffset)

		assignment := domain.OnboardingTaskAssignment{
			OnboardingID: onboardingID,
			AssignedRole: (*string)(&t.AssignedRole),
			Status:       domain.TaskPending,
			DueDate:      &dueDate,
		}
		assignments = append(assignments, assignment)
		_ = i
	}

	return assignments, nil
}

func (e *TaskEngine) GenerateRemoteOnboardingTasks(employeeID, startDate string) []TemplateTaskDef {
	return []TemplateTaskDef{
		{Title: "Enviar notebook", Description: "Enviar equipo de trabajo al domicilio", TaskType: domain.TaskAsset, AssignedRole: domain.ResponsibleIT, Required: true, DaysOffset: -5},
		{Title: "Configurar VPN", Description: "Configurar acceso VPN remoto", TaskType: domain.TaskSystem, AssignedRole: domain.ResponsibleIT, Required: true, DaysOffset: -3},
		{Title: "Crear cuenta de videoconferencia", Description: "Configurar cuenta de Zoom/Meet corporativa", TaskType: domain.TaskAccount, AssignedRole: domain.ResponsibleIT, Required: true, DaysOffset: -2},
		{Title: "Firma digital de documentos", Description: "Firmar documentación digitalmente", TaskType: domain.TaskSignature, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 0},
		{Title: "Reunión virtual de bienvenida", Description: "Video llamada de presentación con el equipo", TaskType: domain.TaskMeeting, AssignedRole: domain.ResponsibleManager, Required: true, DaysOffset: 1},
		{Title: "Capacitación online", Description: "Completar capacitación inicial online", TaskType: domain.TaskTraining, AssignedRole: domain.ResponsibleEmployee, Required: true, DaysOffset: 3},
	}
}

func (e *TaskEngine) GenerateOffboardingTasks(employeeID, terminationType string) []domain.OffboardingTask {
	tasks := []domain.OffboardingTask{
		{Title: "Notificar baja", TaskType: string(domain.OffTaskNotificacion), Required: true},
		{Title: "Preparar documentación de salida", TaskType: string(domain.OffTaskDocumentacion), Required: true},
		{Title: "Realizar entrevista de salida", TaskType: string(domain.OffTaskEntrevista), Required: false},
		{Title: "Recuperar notebook", TaskType: string(domain.OffTaskActivos), Required: true},
		{Title: "Recuperar accesos y credenciales", TaskType: string(domain.OffTaskAccesos), Required: true},
		{Title: "Deshabilitar cuentas de sistemas", TaskType: string(domain.OffTaskAccesos), Required: true},
		{Title: "Transferir proyectos y responsabilidades", TaskType: string(domain.OffTaskTransferencia), Required: true},
		{Title: "Solicitar liquidación final", TaskType: string(domain.OffTaskLiquidacion), Required: true},
	}

	if terminationType == string(domain.TermResignation) {
		tasks = append(tasks, domain.OffboardingTask{
			Title: "Emitir certificado de trabajo", TaskType: string(domain.OffTaskCertificacion), Required: true,
		})
	}

	return tasks
}

func calcDueDate(startDate string, daysOffset int) string {
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Now().AddDate(0, 0, daysOffset).Format("2006-01-02")
	}
	return t.AddDate(0, 0, daysOffset).Format("2006-01-02")
}
