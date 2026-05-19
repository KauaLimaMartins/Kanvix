package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kanvix/backend/internal/models"
)

func EnsureSeeded(ctx context.Context, db *gorm.DB) error {
	var n int64
	if err := db.WithContext(ctx).Model(&models.Workspace{}).Count(&n).Error; err != nil {
		return fmt.Errorf("count workspaces: %w", err)
	}
	if n > 0 {
		return nil
	}
	return SeedDemo(ctx, db)
}

func ResetDemo(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM task_labels;").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM tasks;").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM columns;").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM projects;").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM labels;").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM workspaces;").Error; err != nil {
			return err
		}
		return SeedDemo(ctx, tx)
	})
}

func SeedDemo(ctx context.Context, db *gorm.DB) error {
	now := time.Now().UTC()
	date := func(daysFromNow int) *string {
		s := now.AddDate(0, 0, daysFromNow).Format("2006-01-02")
		return &s
	}

	demoUsers := []models.User{
		{ID: "u1", Email: "you@demo.local", Name: "You", AvatarColor: "#6366f1", CreatedAt: now, UpdatedAt: now},
		{ID: "u2", Email: "alex@demo.local", Name: "Alex Rivera", AvatarColor: "#ec4899", CreatedAt: now, UpdatedAt: now},
		{ID: "u3", Email: "sam@demo.local", Name: "Sam Patel", AvatarColor: "#14b8a6", CreatedAt: now, UpdatedAt: now},
		{ID: "u4", Email: "jordan@demo.local", Name: "Jordan Kim", AvatarColor: "#f59e0b", CreatedAt: now, UpdatedAt: now},
	}

	for _, u := range demoUsers {
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&u).Error; err != nil {
			return fmt.Errorf("seed users: %w", err)
		}
	}

	wsCompany := models.Workspace{ID: "ws-company", OwnerID: "u1", Name: "Company", Icon: "Building2", Color: "#6366f1", CreatedAt: now, UpdatedAt: now}
	wsFreelance := models.Workspace{ID: "ws-freelance", OwnerID: "u1", Name: "Freelances", Icon: "Briefcase", Color: "#14b8a6", CreatedAt: now, UpdatedAt: now}
	wsPersonal := models.Workspace{ID: "ws-personal", OwnerID: "u1", Name: "Personal Projects", Icon: "Sparkles", Color: "#ec4899", CreatedAt: now, UpdatedAt: now}
	workspaces := []models.Workspace{wsCompany, wsFreelance, wsPersonal}
	if err := db.WithContext(ctx).Create(&workspaces).Error; err != nil {
		return fmt.Errorf("seed workspaces: %w", err)
	}

	pWeb := models.Project{ID: "p-web", WorkspaceID: wsCompany.ID, Name: "Website redesign", Description: "Marketing site v2 launch.", CreatedAt: now, UpdatedAt: now}
	pMobile := models.Project{ID: "p-mobile", WorkspaceID: wsCompany.ID, Name: "Mobile app", Description: "iOS + Android MVP.", CreatedAt: now, UpdatedAt: now}
	pHome := models.Project{ID: "p-home", WorkspaceID: wsPersonal.ID, Name: "Home projects", Description: "", CreatedAt: now, UpdatedAt: now}
	pClient := models.Project{ID: "p-client", WorkspaceID: wsFreelance.ID, Name: "Acme landing page", Description: "One-off freelance gig.", CreatedAt: now, UpdatedAt: now}
	projects := []models.Project{pWeb, pMobile, pHome, pClient}
	if err := db.WithContext(ctx).Create(&projects).Error; err != nil {
		return fmt.Errorf("seed projects: %w", err)
	}

	mkCols := func(projectID string) []models.Column {
		return []models.Column{
			{ID: projectID + "-todo", ProjectID: projectID, Name: "To do", Order: 0, CreatedAt: now, UpdatedAt: now},
			{ID: projectID + "-doing", ProjectID: projectID, Name: "In progress", Order: 1, CreatedAt: now, UpdatedAt: now},
			{ID: projectID + "-review", ProjectID: projectID, Name: "Review", Order: 2, CreatedAt: now, UpdatedAt: now},
			{ID: projectID + "-done", ProjectID: projectID, Name: "Done", Order: 3, CreatedAt: now, UpdatedAt: now},
		}
	}

	columns := append(append(append(mkCols(pWeb.ID), mkCols(pMobile.ID)...), mkCols(pHome.ID)...), mkCols(pClient.ID)...)
	if err := db.WithContext(ctx).Create(&columns).Error; err != nil {
		return fmt.Errorf("seed columns: %w", err)
	}

	mkLabels := func(workspaceID string, defs [][2]string) []models.Label {
		out := make([]models.Label, 0, len(defs))
		for _, d := range defs {
			out = append(out, models.Label{
				ID:          uuid.NewString(),
				WorkspaceID: workspaceID,
				Name:        d[0],
				Color:       d[1],
				CreatedAt:   now,
				UpdatedAt:   now,
			})
		}
		return out
	}

	labelsCompany := mkLabels(wsCompany.ID, [][2]string{
		{"Design", "#ec4899"},
		{"Engineering", "#6366f1"},
		{"Research", "#f59e0b"},
		{"Copy", "#14b8a6"},
	})
	labelsFreelance := mkLabels(wsFreelance.ID, [][2]string{
		{"Client A", "#6366f1"},
		{"Client B", "#f97316"},
		{"Invoice", "#10b981"},
	})
	labelsPersonal := mkLabels(wsPersonal.ID, [][2]string{
		{"Home", "#f59e0b"},
		{"Health", "#10b981"},
		{"Ideas", "#a855f7"},
	})

	labels := append(append(labelsCompany, labelsFreelance...), labelsPersonal...)
	if err := db.WithContext(ctx).Create(&labels).Error; err != nil {
		return fmt.Errorf("seed labels: %w", err)
	}

	labelID := func(labels []models.Label, name string) string {
		for _, l := range labels {
			if l.Name == name {
				return l.ID
			}
		}
		return ""
	}

	mkTask := func(projectID, columnID, title string, order int, extras func(*models.Task), labelIDs []string) (models.Task, []models.TaskLabel) {
		t := models.Task{
			ID:          uuid.NewString(),
			ProjectID:   projectID,
			ColumnID:    columnID,
			Title:       title,
			Description: "",
			Order:       order,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if extras != nil {
			extras(&t)
		}
		joins := make([]models.TaskLabel, 0, len(labelIDs))
		for _, lid := range labelIDs {
			if lid == "" {
				continue
			}
			joins = append(joins, models.TaskLabel{TaskID: t.ID, LabelID: lid, CreatedAt: now})
		}
		return t, joins
	}

	var tasks []models.Task
	var taskLabels []models.TaskLabel

	add := func(t models.Task, joins []models.TaskLabel) {
		tasks = append(tasks, t)
		taskLabels = append(taskLabels, joins...)
	}

	t1, j1 := mkTask(pWeb.ID, pWeb.ID+"-todo", "Audit current homepage", 0, func(t *models.Task) {
		a := "u2"
		t.AssigneeID = &a
		t.DueDate = date(4)
	}, []string{labelID(labelsCompany, "Research")})
	add(t1, j1)

	t2, j2 := mkTask(pWeb.ID, pWeb.ID+"-todo", "Collect competitor references", 1, nil, []string{labelID(labelsCompany, "Research")})
	add(t2, j2)

	t3, j3 := mkTask(pWeb.ID, pWeb.ID+"-doing", "Design new hero section", 0, func(t *models.Task) {
		a := "u1"
		t.AssigneeID = &a
		t.DueDate = date(7)
		t.Description = "<p>Explore <strong>three directions</strong>: bold, editorial, minimal.</p><ul><li>Bold typography</li><li>Editorial grid</li><li>Minimal whitespace</li></ul>"
	}, []string{labelID(labelsCompany, "Design")})
	add(t3, j3)

	t4, j4 := mkTask(pWeb.ID, pWeb.ID+"-review", "Pricing page copy", 0, func(t *models.Task) {
		a := "u3"
		t.AssigneeID = &a
		t.DueDate = date(2)
	}, []string{labelID(labelsCompany, "Copy")})
	add(t4, j4)

	t5, j5 := mkTask(pWeb.ID, pWeb.ID+"-done", "Set up analytics", 0, func(t *models.Task) {
		a := "u4"
		t.AssigneeID = &a
		t.DueDate = date(-10)
	}, []string{labelID(labelsCompany, "Engineering")})
	add(t5, j5)

	t13, j13 := mkTask(pWeb.ID, pWeb.ID+"-doing", "Finalize navigation structure", 1, func(t *models.Task) {
		a := "u2"
		t.AssigneeID = &a
		t.DueDate = date(5)
	}, []string{labelID(labelsCompany, "Design"), labelID(labelsCompany, "Research")})
	add(t13, j13)

	t14, j14 := mkTask(pWeb.ID, pWeb.ID+"-todo", "Landing page performance budget", 2, func(t *models.Task) {
		a := "u4"
		t.AssigneeID = &a
		t.DueDate = date(9)
	}, []string{labelID(labelsCompany, "Engineering")})
	add(t14, j14)

	t6, j6 := mkTask(pMobile.ID, pMobile.ID+"-todo", "User onboarding flow", 0, nil, []string{labelID(labelsCompany, "Design")})
	add(t6, j6)

	t7, j7 := mkTask(pMobile.ID, pMobile.ID+"-doing", "Push notifications", 0, func(t *models.Task) {
		a := "u4"
		t.AssigneeID = &a
		t.DueDate = date(6)
	}, []string{labelID(labelsCompany, "Engineering")})
	add(t7, j7)

	t8, j8 := mkTask(pMobile.ID, pMobile.ID+"-done", "App icon", 0, nil, []string{labelID(labelsCompany, "Design")})
	add(t8, j8)

	t15, j15 := mkTask(pMobile.ID, pMobile.ID+"-review", "App Store listing copy", 0, func(t *models.Task) {
		a := "u3"
		t.AssigneeID = &a
		t.DueDate = date(3)
	}, []string{labelID(labelsCompany, "Copy")})
	add(t15, j15)

	t16, j16 := mkTask(pMobile.ID, pMobile.ID+"-todo", "Crash reporting setup", 1, func(t *models.Task) {
		a := "u4"
		t.AssigneeID = &a
		t.DueDate = date(12)
	}, []string{labelID(labelsCompany, "Engineering")})
	add(t16, j16)

	t9, j9 := mkTask(pHome.ID, pHome.ID+"-todo", "Repaint living room", 0, func(t *models.Task) {
		t.DueDate = date(21)
	}, []string{labelID(labelsPersonal, "Home")})
	add(t9, j9)

	t10, j10 := mkTask(pHome.ID, pHome.ID+"-doing", "Plan kitchen remodel", 0, func(t *models.Task) {
		t.DueDate = date(30)
	}, []string{labelID(labelsPersonal, "Home")})
	add(t10, j10)

	t17, j17 := mkTask(pHome.ID, pHome.ID+"-review", "Book annual checkup", 0, func(t *models.Task) {
		t.DueDate = date(14)
	}, []string{labelID(labelsPersonal, "Health")})
	add(t17, j17)

	t11, j11 := mkTask(pClient.ID, pClient.ID+"-todo", "Wireframe homepage", 0, func(t *models.Task) {
		a := "u1"
		t.AssigneeID = &a
		t.DueDate = date(1)
	}, []string{labelID(labelsFreelance, "Client A")})
	add(t11, j11)

	t12, j12 := mkTask(pClient.ID, pClient.ID+"-doing", "Send first invoice", 0, func(t *models.Task) {
		a := "u1"
		t.AssigneeID = &a
		t.DueDate = date(0)
	}, []string{labelID(labelsFreelance, "Invoice")})
	add(t12, j12)

	t18, j18 := mkTask(pClient.ID, pClient.ID+"-done", "Kickoff call notes", 0, func(t *models.Task) {
		t.DueDate = date(-3)
		t.Description = "<p>Capture scope, deliverables, and deadlines.</p>"
	}, []string{labelID(labelsFreelance, "Client A")})
	add(t18, j18)

	if err := db.WithContext(ctx).Create(&tasks).Error; err != nil {
		return fmt.Errorf("seed tasks: %w", err)
	}
	if len(taskLabels) > 0 {
		if err := db.WithContext(ctx).Create(&taskLabels).Error; err != nil {
			return fmt.Errorf("seed task_labels: %w", err)
		}
	}

	return nil
}
