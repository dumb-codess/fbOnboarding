package main

import (
	"encoding/json"
	"errors"
	"fbOnboarding/db"
	s3base "fbOnboarding/s3"
	"io"
	"net/http"
	"strconv"
)

type Controller struct {
	app Application
}

func NewController() (*Controller, error) {
	db := db.GetDBlient()
	s3, err := s3base.NewS3Client(cfg.AwsEnpoint)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	fbapp, err := NewApplication(db, s3)
	if err != nil {
		logger.Fatal().Msgf("Failed create Application: %v", err.Error())
		return nil, err
	}

	return &Controller{
		app: fbapp,
	}, nil
}

func (c *Controller) GetToken(req *http.Request) (*Response, error) {
	ctx := req.Context()
	var input GetTokenReq

	data := req.URL.Query().Get("customer_id")
	if data == "" {
		logger.Error().Msgf("customer_id not found in url")
		return nil, errors.New("customer_id not found ")
	}

	cID, _ := strconv.Atoi(data)
	logger.Info().Msgf("request from Customer id %v", cID)
	input.CustomerID = int64(cID)

	token, err := c.app.GetToken(ctx, input)
	if err != nil {
		logger.Error().Msgf("failed to get token: %v", err)
		return nil, err
	}

	return &Response{
		Data:       token,
		StatusCode: http.StatusOK,
	}, nil
}

func (c *Controller) CheckSubmission(req *http.Request) (*Response, error) {
	ctx := req.Context()
	var input CheckSubmissionReq

	customerID, ok := ctx.Value("customer_id").(int64)
	if !ok {
		logger.Error().Msg("customer_id type assertion error")
		return nil, errors.New("customer_id type error")
	}

	input.CustomerID = customerID

	res, err := c.app.CheckSubmission(ctx, input)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	logger.Info().Msgf("customer with id %v with submission status %v", customerID, res.SubmissionStatus)
	return &Response{
		Data:       res,
		StatusCode: http.StatusOK,
	}, nil

}

func (c *Controller) Uploadfile(req *http.Request) (*Response, error) {
	ctx := req.Context()
	var input UploadFileReq
	customerID, ok := ctx.Value("customer_id").(int64)
	if !ok {
		logger.Error().Msg("customer_id type assertion error")
		return nil, errors.New("customer_id type error")
	}

	logger.Info().Msgf("customer %v request to upload file %v", customerID, input.FileName)

	input.CustomerID = customerID
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		logger.Error().Msgf("failed to get file from request body: %v", err)
		return nil, err
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		logger.Error().Msgf("failed to get file format: %v", err)
		return nil, err
	}
	defer file.Close()

	fBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error().Msgf("failed to read file: %v", err)
		return nil, err
	}

	input.FileBytes = fBytes
	input.FileName = handler.Filename

	res, err := c.app.UploadFile(ctx, input)
	if err != nil {
		logger.Error().Msgf("failed to upload file: %v", err)
		return nil, err
	}

	return &Response{
		Data:       res,
		StatusCode: http.StatusOK,
	}, nil

}

func (c *Controller) GetApprovalStatus(req *http.Request) (*Response, error) {
	ctx := req.Context()
	var input GetApprovalStatusReq

	customerID, ok := ctx.Value("customer_id").(int64)
	if !ok {
		logger.Error().Msg("custoner_id type assertion error")
		return nil, errors.New("customer_id type error")
	}

	input.CustomerID = customerID

	res, err := c.app.GetApprovalStatus(ctx, input)
	if err != nil {
		logger.Error().Msgf("failed to get approval Status: %v", err)
		return nil, err
	}

	logger.Info().Msgf("customer %v with form approval status %v", customerID, res.FormStatus)
	return &Response{
		Data:       res,
		StatusCode: http.StatusOK,
	}, nil
}

func (c *Controller) ApproveApplicationStaus(req *http.Request) (*Response, error) {
	ctx := req.Context()
	var input ApproveFormReq

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	res, err := c.app.ApproveApplicationStaus(ctx, input)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	logger.Info().Msgf("Application with id %v has approval status %v", input.FormID, input.FormStatus)
	return &Response{
		Data:       res,
		StatusCode: http.StatusOK,
	}, nil
}
