package handler

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
	"workpulse/internal/repo"
	"workpulse/internal/service"
)

type ReportHandler struct{ svc *service.ReportService }

func NewReportHandler(svc *service.ReportService) *ReportHandler { return &ReportHandler{svc: svc} }

type reportLinkPayload struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   string `json:"target_id" binding:"required"`
	Relation   string `json:"relation"`
}

type upsertReportReq struct {
	Type        string              `json:"type" binding:"required"`
	PeriodStart string              `json:"period_start" binding:"required"`
	PeriodEnd   string              `json:"period_end" binding:"required"`
	TeamID      *string             `json:"team_id"`
	Title       *string             `json:"title"`
	Summary     *string             `json:"summary"`
	Content     map[string]any      `json:"content"`
	Payload     map[string]any      `json:"payload"`
	Links       []reportLinkPayload `json:"links"`
}

func (h *ReportHandler) List(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	filter := repo.ReportFilter{}
	if v := c.Query("type"); v != "" {
		filter.Type = &[]string{v}[0]
	}
	if v := c.Query("status"); v != "" {
		filter.Status = &[]string{v}[0]
	}
	if v := c.Query("author_user_id"); v != "" {
		id := uuid.MustParse(v)
		filter.AuthorID = &id
	}
	if v := c.Query("from"); v != "" {
		if t, err := parseDate(v); err == nil {
			filter.PeriodStart = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := parseDate(v); err == nil {
			filter.PeriodEnd = &t
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			filter.Offset = p * filter.Limit
		}
	}

	reports, total, err := h.svc.List(c.Request.Context(), orgID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": reports, "total": total}})
}

func (h *ReportHandler) Get(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	report, links, err := h.svc.Get(c.Request.Context(), orgID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"report": report, "links": links}})
}

func (h *ReportHandler) Create(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	actor := uuid.MustParse(c.GetString("user_id"))
	var req upsertReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ps, pe, err := parsePeriod(req.PeriodStart, req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contentBytes, _ := json.Marshal(req.Content)
	payloadBytes, _ := json.Marshal(req.Payload)

	report := &models.Report{
		OrgID:         orgID,
		AuthorUserID:  actor,
		Type:          req.Type,
		PeriodStart:   ps,
		PeriodEnd:     pe,
		Title:         req.Title,
		Summary:       req.Summary,
		Content:       datatypes.JSON(contentBytes),
		Payload:       datatypes.JSON(payloadBytes),
		SchemaVersion: 1,
	}
	if req.TeamID != nil && *req.TeamID != "" {
		team := uuid.MustParse(*req.TeamID)
		report.TeamID = &team
	}
	links, err := toLinks(req.Links)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Create(c.Request.Context(), actor, report, links); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"report": report, "links": links}})
}

func (h *ReportHandler) Update(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	actor := uuid.MustParse(c.GetString("user_id"))
	id := uuid.MustParse(c.Param("id"))
	var req upsertReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ps, pe, err := parsePeriod(req.PeriodStart, req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contentBytes, _ := json.Marshal(req.Content)
	payloadBytes, _ := json.Marshal(req.Payload)

	report := &models.Report{
		ID:            id,
		OrgID:         orgID,
		AuthorUserID:  actor,
		Type:          req.Type,
		PeriodStart:   ps,
		PeriodEnd:     pe,
		Title:         req.Title,
		Summary:       req.Summary,
		Content:       datatypes.JSON(contentBytes),
		Payload:       datatypes.JSON(payloadBytes),
		SchemaVersion: 1,
	}
	if req.TeamID != nil && *req.TeamID != "" {
		team := uuid.MustParse(*req.TeamID)
		report.TeamID = &team
	}
	links, err := toLinks(req.Links)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Update(c.Request.Context(), actor, report, links); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"report": report, "links": links}})
}

func (h *ReportHandler) Submit(c *gin.Context)  { h.changeStatus(c, "submit") }
func (h *ReportHandler) Approve(c *gin.Context) { h.changeStatus(c, "approve") }
func (h *ReportHandler) Reject(c *gin.Context)  { h.changeStatus(c, "reject") }

func (h *ReportHandler) changeStatus(c *gin.Context, action string) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	actor := uuid.MustParse(c.GetString("user_id"))
	id := uuid.MustParse(c.Param("id"))
	var req struct {
		Comment *string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&req)
	report, err := h.svc.ChangeStatus(c.Request.Context(), orgID, id, actor, action, req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

func (h *ReportHandler) Export(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	filter := repo.ReportFilter{Limit: 500}
	if v := c.Query("type"); v != "" {
		filter.Type = &[]string{v}[0]
	}
	if v := c.Query("status"); v != "" {
		filter.Status = &[]string{v}[0]
	}
	if v := c.Query("author_user_id"); v != "" {
		id := uuid.MustParse(v)
		filter.AuthorID = &id
	}
	if v := c.Query("from"); v != "" {
		if t, err := parseDate(v); err == nil {
			filter.PeriodStart = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := parseDate(v); err == nil {
			filter.PeriodEnd = &t
		}
	}

	reports, _, err := h.svc.List(c.Request.Context(), orgID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	_ = writer.Write([]string{"id", "type", "period_start", "period_end", "status", "title", "author_user_id"})
	for _, r := range reports {
		_ = writer.Write([]string{
			r.ID.String(),
			r.Type,
			r.PeriodStart.Format("2006-01-02"),
			r.PeriodEnd.Format("2006-01-02"),
			r.Status,
			valueOrEmpty(r.Title),
			r.AuthorUserID.String(),
		})
	}
	writer.Flush()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=reports.csv")
	c.String(http.StatusOK, buf.String())
}

func toLinks(payloads []reportLinkPayload) ([]models.ReportLink, error) {
	links := make([]models.ReportLink, 0, len(payloads))
	for _, p := range payloads {
		tid, err := uuid.Parse(p.TargetID)
		if err != nil {
			return nil, err
		}
		rel := p.Relation
		if rel == "" {
			rel = "related"
		}
		links = append(links, models.ReportLink{
			TargetType: p.TargetType,
			TargetID:   tid,
			Relation:   rel,
		})
	}
	return links, nil
}

func parsePeriod(start, end string) (time.Time, time.Time, error) {
	ps, err := parseDate(start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	pe, err := parseDate(end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return ps, pe, nil
}

func parseDate(v string) (time.Time, error) {
	return time.Parse("2006-01-02", v)
}

func valueOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
