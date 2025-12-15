package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
	"workpulse/internal/service"
)

type OKRHandler struct{ svc *service.OKRService }

func NewOKRHandler(svc *service.OKRService) *OKRHandler { return &OKRHandler{svc: svc} }

type createObjectiveReq struct {
	CycleID     string   `json:"cycle_id" binding:"required"`
	TeamID      *string  `json:"team_id"`
	OwnerUserID *string  `json:"owner_user_id"`
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
}

type updateObjectiveReq struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
}

type createKRReq struct {
	ObjectiveID string   `json:"objective_id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	MetricType  string   `json:"metric_type"`
	TargetValue *float64 `json:"target_value"`
	Unit        *string  `json:"unit"`
}

type updateKRReq struct {
	Title        string   `json:"title"`
	MetricType   string   `json:"metric_type"`
	TargetValue  *float64 `json:"target_value"`
	CurrentValue *float64 `json:"current_value"`
	Confidence   *int     `json:"confidence"`
	Status       string   `json:"status"`
}

type linkReq struct {
	ObjectiveID *string `json:"objective_id"`
	KeyResultID *string `json:"key_result_id"`
	EntityType  string  `json:"entity_type" binding:"required"`
	EntityID    string  `json:"entity_id" binding:"required"`
	Relation    string  `json:"relation"`
}

func (h *OKRHandler) ListObjectives(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var cycleID, teamID, ownerID *uuid.UUID
	var status *string
	if v := c.Query("cycle_id"); v != "" {
		id := uuid.MustParse(v)
		cycleID = &id
	}
	if v := c.Query("team_id"); v != "" {
		id := uuid.MustParse(v)
		teamID = &id
	}
	if v := c.Query("owner_user_id"); v != "" {
		id := uuid.MustParse(v)
		ownerID = &id
	}
	if v := c.Query("status"); v != "" {
		status = &[]string{v}[0]
	}
	objs, krs, err := h.svc.ListObjectives(c.Request.Context(), orgID, cycleID, teamID, ownerID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"objectives": objs, "key_results": krs}})
}

func (h *OKRHandler) CreateObjective(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	userID := c.GetString("user_id")
	var req createObjectiveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := &models.OKRObjective{
		ID:            uuid.New(),
		OrgID:         orgID,
		CycleID:       uuid.MustParse(req.CycleID),
		Title:         req.Title,
		Description:   req.Description,
		Status:        defaultString(req.Status, "active"),
		SchemaVersion: 1,
		Tags:          datatypes.JSON([]byte("[]")),
		Payload:       datatypes.JSON([]byte("{}")),
	}
	if req.OwnerUserID != nil {
		id := uuid.MustParse(*req.OwnerUserID)
		obj.OwnerUserID = id
	} else if userID != "" {
		obj.OwnerUserID = uuid.MustParse(userID)
	}
	if req.TeamID != nil && *req.TeamID != "" {
		id := uuid.MustParse(*req.TeamID)
		obj.TeamID = &id
	}
	if len(req.Tags) > 0 {
		obj.Tags = datatypes.JSON([]byte(toJSONList(req.Tags)))
	}
	if err := h.svc.CreateObjective(c.Request.Context(), obj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": obj})
}

func (h *OKRHandler) UpdateObjective(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req updateObjectiveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := &models.OKRObjective{ID: id, OrgID: orgID, SchemaVersion: 1}
	if req.Title != "" {
		obj.Title = req.Title
	}
	obj.Description = req.Description
	if req.Status != "" {
		obj.Status = req.Status
	}
	if len(req.Tags) > 0 {
		obj.Tags = datatypes.JSON([]byte(toJSONList(req.Tags)))
	}
	if err := h.svc.UpdateObjective(c.Request.Context(), obj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": obj})
}

func (h *OKRHandler) ArchiveObjective(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	if err := h.svc.ArchiveObjective(c.Request.Context(), orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "archived"})
}

func (h *OKRHandler) CreateKeyResult(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req createKRReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kr := &models.OKRKeyResult{
		ID:            uuid.New(),
		OrgID:         orgID,
		ObjectiveID:   uuid.MustParse(req.ObjectiveID),
		Title:         req.Title,
		MetricType:    defaultString(req.MetricType, "number"),
		TargetValue:   req.TargetValue,
		Unit:          req.Unit,
		SchemaVersion: 1,
		MetricPayload: datatypes.JSON([]byte("{}")),
	}
	if err := h.svc.CreateKeyResult(c.Request.Context(), kr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kr})
}

func (h *OKRHandler) UpdateKeyResult(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req updateKRReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kr := &models.OKRKeyResult{ID: id, OrgID: orgID, SchemaVersion: 1}
	if req.Title != "" {
		kr.Title = req.Title
	}
	if req.MetricType != "" {
		kr.MetricType = req.MetricType
	}
	if req.TargetValue != nil {
		kr.TargetValue = req.TargetValue
	}
	if req.CurrentValue != nil {
		kr.CurrentValue = req.CurrentValue
	}
	if req.Confidence != nil {
		kr.Confidence = *req.Confidence
	}
	if req.Status != "" {
		kr.Status = req.Status
	}
	if err := h.svc.UpdateKeyResult(c.Request.Context(), kr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kr})
}

func (h *OKRHandler) ArchiveKeyResult(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	if err := h.svc.ArchiveKeyResult(c.Request.Context(), orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "archived"})
}

func (h *OKRHandler) UpdateProgress(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req struct {
		Current    float64 `json:"current"`
		Confidence int     `json:"confidence"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateProgress(c.Request.Context(), orgID, id, req.Current, req.Confidence); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "updated"})
}

func (h *OKRHandler) AddLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req linkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ObjectiveID == nil && req.KeyResultID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "objective_id or key_result_id is required"})
		return
	}
	link := &models.OKRLink{
		ID:            uuid.New(),
		OrgID:         orgID,
		EntityType:    req.EntityType,
		EntityID:      uuid.MustParse(req.EntityID),
		Relation:      defaultString(req.Relation, "related"),
		SchemaVersion: 1,
	}
	if req.ObjectiveID != nil {
		id := uuid.MustParse(*req.ObjectiveID)
		link.ObjectiveID = &id
	}
	if req.KeyResultID != nil {
		id := uuid.MustParse(*req.KeyResultID)
		link.KeyResultID = &id
	}
	if err := h.svc.AddLink(c.Request.Context(), link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": link})
}

func (h *OKRHandler) RemoveLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	if err := h.svc.RemoveLink(c.Request.Context(), orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "unlinked"})
}

func (h *OKRHandler) Metrics(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var cycleID, teamID, ownerID *uuid.UUID
	if v := c.Query("cycle_id"); v != "" {
		id := uuid.MustParse(v)
		cycleID = &id
	}
	if v := c.Query("team_id"); v != "" {
		id := uuid.MustParse(v)
		teamID = &id
	}
	if v := c.Query("owner_user_id"); v != "" {
		id := uuid.MustParse(v)
		ownerID = &id
	}
	metrics, err := h.svc.ComputeMetrics(c.Request.Context(), orgID, cycleID, teamID, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": metrics})
}

func defaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func toJSONList(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	encoded := "["
	for i, v := range values {
		encoded += "\"" + v + "\""
		if i < len(values)-1 {
			encoded += ","
		}
	}
	encoded += "]"
	return encoded
}
