package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"oss.kftd.co.id/v2/submission-management/dto/request"
	"oss.kftd.co.id/v2/submission-management/dto/response"
	"oss.kftd.co.id/v2/submission-management/model"
	"oss.kftd.co.id/v2/submission-management/repository"
)

const (
	SubmissionStatusDraft = "DRAFT"
)

type SubmissionService struct {
	db             *gorm.DB
	submissionRepo *repository.SubmissionRepository
}

func NewSubmissionService(
	db *gorm.DB,
	submissionRepo *repository.SubmissionRepository,
) *SubmissionService {
	return &SubmissionService{
		db:             db,
		submissionRepo: submissionRepo,
	}
}

// CreateSubmission membuat submission header,
// customer, material, dan status history
// dalam satu transaction.
func (s *SubmissionService) CreateSubmission(
	ctx context.Context,
	req request.CreateSubmissionRequest,
	userID uint64,
) (*response.SubmissionResponse, error) {

	// =========================================
	// BEGIN TRANSACTION
	// =========================================

	tx := s.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Jika terjadi panic, rollback transaction.
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Repository menggunakan transaction yang sama.
	txRepo := s.submissionRepo.WithTx(tx)

	// =========================================
	// CREATE SUBMISSION
	// =========================================

	submission := &model.Submission{
		UUID:      uuid.NewString(),
		SofficeID: req.SofficeID,
		VariantID: req.VariantID,
		ValueType: req.ValueType,
		Status:    SubmissionStatusDraft,
		CreatedBy: userID,
	}

	if req.Notes != "" {
		submission.Notes = &req.Notes
	}

	err := txRepo.Create(
		ctx,
		submission,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// =========================================
	// CREATE CUSTOMERS
	// =========================================

	customers := make(
		[]model.SubmissionCustomer,
		0,
		len(req.Customers),
	)

	for _, item := range req.Customers {
		customers = append(
			customers,
			model.SubmissionCustomer{
				UUID:         uuid.NewString(),
				SubmissionID: submission.ID,
				CustomerID:   item.CustomerID,
				CreatedBy:    userID,
			},
		)
	}

	err = txRepo.CreateCustomers(
		ctx,
		customers,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// =========================================
	// CREATE MATERIALS
	// =========================================

	materials := make(
		[]model.SubmissionMaterial,
		0,
		len(req.Materials),
	)

	for _, item := range req.Materials {
		materials = append(
			materials,
			model.SubmissionMaterial{
				UUID:         uuid.NewString(),
				SubmissionID: submission.ID,
				MaterialID:   item.MaterialID,
				QtyJual:      item.QtyJual,
				SalesUOM:     item.SalesUOM,
				Value:        item.Value,
				CreatedBy:    userID,
			},
		)
	}

	err = txRepo.CreateMaterials(
		ctx,
		materials,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// =========================================
	// CREATE STATUS HISTORY
	// =========================================

	statusHistory := &model.SubmissionStatus{
		UUID:         uuid.NewString(),
		SubmissionID: submission.ID,
		FromStatus:   "",
		ToStatus:     SubmissionStatusDraft,
		ActorID:      userID,
	}

	err = txRepo.CreateStatus(
		ctx,
		statusHistory,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// =========================================
	// COMMIT
	// =========================================

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// =========================================
	// RESPONSE
	// =========================================

	return s.buildResponse(
		submission,
		customers,
		materials,
	), nil
}

// GetSubmission mengambil submission berdasarkan UUID
// beserta customer dan material.
func (s *SubmissionService) GetSubmission(
	ctx context.Context,
	submissionUUID string,
) (*response.SubmissionResponse, error) {

	submission, err := s.submissionRepo.FindSubmissionByUUID(
		ctx,
		submissionUUID,
	)

	if err != nil {
		return nil, err
	}

	customers, err := s.submissionRepo.FindCustomersBySubmissionID(
		ctx,
		submission.ID,
	)

	if err != nil {
		return nil, err
	}

	materials, err := s.submissionRepo.FindMaterialsBySubmissionID(
		ctx,
		submission.ID,
	)

	if err != nil {
		return nil, err
	}

	return s.buildResponse(
		submission,
		customers,
		materials,
	), nil
}

// buildResponse melakukan mapping model ke response DTO.
func (s *SubmissionService) buildResponse(
	submission *model.Submission,
	customers []model.SubmissionCustomer,
	materials []model.SubmissionMaterial,
) *response.SubmissionResponse {

	result := &response.SubmissionResponse{
		UUID:        submission.UUID,
		BusinessID:  submission.BusinessID,
		SofficeID:   submission.SofficeID,
		VariantID:   submission.VariantID,
		ValueType:   submission.ValueType,
		CurrentRole: submission.CurrentRole,
		Status:      submission.Status,
		CreatedBy:   submission.CreatedBy,

		Customers: make(
			[]response.SubmissionCustomerResponse,
			0,
			len(customers),
		),

		Materials: make(
			[]response.SubmissionMaterialResponse,
			0,
			len(materials),
		),
	}

	if submission.Notes != nil {
		result.Notes = *submission.Notes
	}

	for _, customer := range customers {
		result.Customers = append(
			result.Customers,
			response.SubmissionCustomerResponse{
				UUID:       customer.UUID,
				CustomerID: customer.CustomerID,
			},
		)
	}

	for _, material := range materials {
		result.Materials = append(
			result.Materials,
			response.SubmissionMaterialResponse{
				UUID:       material.UUID,
				MaterialID: material.MaterialID,
				QtyJual:    material.QtyJual,
				SalesUOM:   material.SalesUOM,
				Value:      material.Value,
			},
		)
	}

	return result
}
