package main

import (
	"context"
	"fbOnboarding/ent"
	"fbOnboarding/ent/consumer"
	"fbOnboarding/ent/uploadedfile"
	"fbOnboarding/enum"
	s3base "fbOnboarding/s3"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Application interface {
	GetToken(ctx context.Context, input GetTokenReq) (*GetTokenResp, error)
	CheckSubmission(ctx context.Context, input CheckSubmissionReq) (*CheckSubmissionResp, error)
	UploadFile(ctx context.Context, input UploadFileReq) (*UploadFileResp, error)
	GetApprovalStatus(ctx context.Context, input GetApprovalStatusReq) (*GetApprovalStatusResp, error)
	ApproveApplicationStaus(ctx context.Context, input ApproveFormReq) (*ApproveFormResp, error)
}

type FbApp struct {
	dbClient *ent.Client
	s3Client *s3base.S3Client
}

func NewApplication(dbClient *ent.Client, s3Client *s3base.S3Client) (Application, error) {
	return &FbApp{
		dbClient: dbClient,
		s3Client: s3Client,
	}, nil
}

type CustomClaims struct {
	CustomerID int64 `json:"customer_id"`
	jwt.RegisteredClaims
}

func (fb *FbApp) GetToken(ctx context.Context, input GetTokenReq) (*GetTokenResp, error) {
	claims := CustomClaims{
		CustomerID: input.CustomerID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpiresInDuration)),
			Issuer:    "fb-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}

	var customer *ent.Consumer
	customer, err = fb.dbClient.Consumer.Query().Where(consumer.ID(input.CustomerID)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if customer == nil {
		customer, err = fb.dbClient.Consumer.Create().SetID(input.CustomerID).SetSubmissionStatus(false).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &GetTokenResp{
		Token:            signedToken,
		SubmissionStatus: customer.SubmissionStatus,
	}, nil
}

func (fb *FbApp) CheckSubmission(ctx context.Context, input CheckSubmissionReq) (*CheckSubmissionResp, error) {
	var customer *ent.Consumer

	customer, err := fb.dbClient.Consumer.Query().Where(consumer.ID(input.CustomerID)).Only(ctx)
	if err != nil {
		return nil, err
	}

	return &CheckSubmissionResp{
		SubmissionStatus: customer.SubmissionStatus,
	}, nil
}

func (fb *FbApp) UploadFile(ctx context.Context, input UploadFileReq) (*UploadFileResp, error) {
	key := fmt.Sprintf("fb/uploads/%v/%v", input.CustomerID, input.FileName)
	bucket := "fbbucket"

	out, err := fb.dbClient.UploadedFile.Query().Where(uploadedfile.HasConsumerWith(consumer.ID(input.CustomerID))).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		logger.Error().Msgf("failed to get uploaded exist: %v", err)
		return nil, err
	}

	if out != nil {
		out, err := fb.dbClient.UploadedFile.UpdateOneID(out.ID).SetKey(key).Save(ctx)
		if err != nil {
			logger.Error().Msgf("failed to get update uploaded file: %v", err)
			return nil, err
		}
		return &UploadFileResp{
			FormStatus: out.FormStatus,
			CustomerID: input.CustomerID,
			FormID:     out.ID,
		}, nil
	}

	if err := fb.s3Client.UploadFile(context.Background(), bucket, key, input.FileBytes); err != nil {
		logger.Error().Msgf("failed to upload to s3: %v", err)
		return nil, err
	}

	upload, err := fb.dbClient.UploadedFile.Create().
		SetID(input.CustomerID).
		SetFormStatus(enum.FormStatusPENDING).
		SetKey(key).
		Save(ctx)
	if err != nil {
		logger.Error().Msgf("failed to get ent uploadfile : %v", err)
		return nil, err
	}

	if err := fb.dbClient.Consumer.UpdateOneID(input.CustomerID).
		SetUploadedfileID(upload.ID).
		SetSubmissionStatus(true).
		Exec(ctx); err != nil {
		logger.Error().Msgf("failed to update consumer submission status: %v", err)
		return nil, err
	}

	return &UploadFileResp{
		FormStatus: upload.FormStatus,
		CustomerID: input.CustomerID,
		FormID:     upload.ID,
	}, nil
}

func (fb *FbApp) GetApprovalStatus(ctx context.Context, input GetApprovalStatusReq) (*GetApprovalStatusResp, error) {
	out, err := fb.dbClient.UploadedFile.Query().Where(uploadedfile.HasConsumerWith(consumer.ID(input.CustomerID))).Only(ctx)
	if err != nil {
		return nil, err
	}

	return &GetApprovalStatusResp{
		FormStatus: out.FormStatus,
		CustomerID: out.ID,
	}, nil
}

func (fb *FbApp) ApproveApplicationStaus(ctx context.Context, input ApproveFormReq) (*ApproveFormResp, error) {
	out, err := fb.dbClient.UploadedFile.UpdateOneID(input.FormID).SetFormStatus(input.FormStatus).Save(ctx)
	if err != nil {
		return nil, err
	}

	consumer, err := out.QueryConsumer().Only(ctx)
	if err != nil {
		return nil, err
	}

	return &ApproveFormResp{
		CustomerID: consumer.ID,
		FormID:     out.ID,
		FormStatus: out.FormStatus,
	}, err

}
