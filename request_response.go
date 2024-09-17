package main

import "fbOnboarding/enum"

type (
	GetTokenReq struct {
		CustomerID int64 `json:"customer_id"`
	}

	GetTokenResp struct {
		Token            string `json:"token"`
		SubmissionStatus bool   `json:"submission_status"`
	}

	CheckSubmissionReq struct {
		CustomerID int64 `json:"customer_id"`
	}

	CheckSubmissionResp struct {
		SubmissionStatus bool   `json:"submission_status"`
		Message          string `json:"message"`
	}

	UploadFileReq struct {
		FileBytes  []byte `json:"-"`
		FileName   string `json:"filename"`
		CustomerID int64  `json:"customer_id"`
	}

	UploadFileResp struct {
		FormStatus enum.FormStatus `json:"form_status"`
		CustomerID int64           `json:"customer_id"`
		FormID     int64           `json:"form_id"`
	}

	GetApprovalStatusReq struct {
		CustomerID int64 `json:"customer_id"`
	}

	GetApprovalStatusResp struct {
		FormStatus enum.FormStatus `json:"form_status"`
		CustomerID int64           `json:"customer_id"`
	}

	ApproveFormReq struct {
		FormID     int64           `json:"form_id"`
		FormStatus enum.FormStatus `json:"form_status"`
	}

	ApproveFormResp struct {
		CustomerID int64           `json:"customer_id"`
		FormID     int64           `json:"form_id"`
		FormStatus enum.FormStatus `json:"form_status"`
	}
)
