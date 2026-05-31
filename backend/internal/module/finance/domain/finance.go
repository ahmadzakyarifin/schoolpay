package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type BillType struct {
	bun.BaseModel `bun:"table:bill_types,alias:bt"`

	ID            uint       `bun:"id,pk,autoincrement" json:"id"`
	Name          string     `bun:"name" json:"name" binding:"required"`
	Description   *string    `bun:"description" json:"description"`
	Type          string     `bun:"type" json:"type" binding:"required"` // one_time, recurring
	DefaultAmount float64    `bun:"default_amount" json:"default_amount"`
	IsActive      bool       `bun:"is_active,default:true" json:"is_active"`
	CreatedAt     time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
	RuleCount     int        `bun:"rule_count,scanonly" json:"rule_count"`
}

type BillingRule struct {
	bun.BaseModel `bun:"table:billing_rules,alias:br"`

	ID               uint       `bun:"id,pk,autoincrement" json:"id"`
	BillTypeID       uint       `bun:"bill_type_id" json:"bill_type_id" binding:"required"`
	ClassID          *uint      `bun:"class_id" json:"class_id,omitempty"`
	TargetType       string     `bun:"target_type" json:"target_type"`
	TargetID         uint       `bun:"target_id" json:"target_id"`
	Amount           float64    `bun:"amount" json:"amount" binding:"required"`
	PeriodType       *string    `bun:"period_type" json:"period_type"` // bulanan, tahunan
	AllowInstallment bool       `bun:"allow_installment" json:"allow_installment"`
	MaxInstallment   *int       `bun:"max_installment" json:"max_installment,omitempty"`
	DueDay           int        `bun:"due_day,default:10" json:"due_day"`
	StartDate        *time.Time `bun:"start_date" json:"start_date,omitempty"`
	EndDate          *time.Time `bun:"end_date" json:"end_date,omitempty"`
	IsActive         bool       `bun:"is_active,default:true" json:"is_active"`
	CreatedAt        time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt        time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt        *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`

	// Join
	BillTypeName string `bun:"bill_type_name,scanonly" json:"bill_type_name,omitempty"`
	ClassName    string `bun:"class_name,scanonly" json:"class_name,omitempty"`
	BillCount    int    `bun:"bill_count,scanonly" json:"bill_count"`
}

type StudentBill struct {
	bun.BaseModel `bun:"table:student_bills,alias:sb"`

	ID              uint       `bun:"id,pk,autoincrement" json:"id"`
	StudentID       uint       `bun:"student_id" json:"student_id" binding:"required"`
	BillTypeID      uint       `bun:"bill_type_id" json:"bill_type_id" binding:"required"`
	BillingRuleID   *uint      `bun:"billing_rule_id" json:"billing_rule_id,omitempty"`
	Name            *string    `bun:"name" json:"name,omitempty"`
	AcademicYear    string     `bun:"academic_year" json:"academic_year"`
	Period          *string    `bun:"period" json:"period"` // e.g. "2024-07"
	PeriodMonth     *int       `bun:"period_month" json:"period_month,omitempty"`
	PeriodYear      *int       `bun:"period_year" json:"period_year,omitempty"`
	PeriodStartDate *time.Time `bun:"period_start_date" json:"period_start_date,omitempty"`
	PeriodEndDate   *time.Time `bun:"period_end_date" json:"period_end_date,omitempty"`
	Amount          float64    `bun:"amount" json:"amount" binding:"required"`
	TotalPaid       float64    `bun:"total_paid,default:0" json:"total_paid"`
	Status          string     `bun:"status,default:unpaid" json:"status"` // unpaid, partial, paid, overdue
	DueDate         time.Time  `bun:"due_date" json:"due_date"`
	EndDate         *time.Time `bun:"end_date" json:"end_date,omitempty"`
	Description     string     `bun:"description" json:"description,omitempty"`
	CreatedAt       time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt       *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
	VoidedAt        *time.Time `bun:"voided_at" json:"voided_at,omitempty"`
	VoidReason      *string    `bun:"void_reason" json:"void_reason,omitempty"`

	LastNotifiedAt *time.Time `bun:"last_notified_at" json:"last_notified_at,omitempty"`
	NextNotifiedAt *time.Time `bun:"next_notified_at" json:"next_notified_at,omitempty"`
	UQKeyPeriod    *string    `bun:"uq_key_period,unique" json:"uq_key_period,omitempty"`

	// Join
	StudentName      string     `bun:"student_name,scanonly" json:"student_name,omitempty"`
	DepositBalance   float64    `bun:"deposit_balance,scanonly" json:"deposit_balance"`
	BillTypeName     string     `bun:"bill_type_name,scanonly" json:"bill_type_name,omitempty"`
	AllowInstallment bool       `bun:"allow_installment,scanonly" json:"allow_installment"`
	MaxInstallment   int        `bun:"max_installment,scanonly" json:"max_installment"`
	RuleStartDate    *time.Time `bun:"rule_start_date,scanonly" json:"rule_start_date,omitempty"`
	RuleEndDate      *time.Time `bun:"rule_end_date,scanonly" json:"rule_end_date,omitempty"`
}

type StudentBillSummary struct {
	ID             uint       `bun:"id" json:"id"`
	StudentID      uint       `bun:"student_id" json:"student_id"`
	StudentName    string     `bun:"student_name" json:"student_name"`
	Amount         float64    `bun:"amount" json:"amount"`
	TotalPaid      float64    `bun:"total_paid" json:"total_paid"`
	Outstanding    float64    `bun:"outstanding" json:"outstanding"`
	BillCount      int        `bun:"bill_count" json:"bill_count"`
	PaidCount      int        `bun:"paid_count" json:"paid_count"`
	PartialCount   int        `bun:"partial_count" json:"partial_count"`
	OverdueCount   int        `bun:"overdue_count" json:"overdue_count"`
	UnpaidCount    int        `bun:"unpaid_count" json:"unpaid_count"`
	DepositBalance float64    `bun:"deposit_balance" json:"deposit_balance"`
	Status         string     `bun:"status" json:"status"`
	NearestDueDate *time.Time `bun:"nearest_due_date" json:"nearest_due_date,omitempty"`
	LastBillAt     *time.Time `bun:"last_bill_at" json:"last_bill_at,omitempty"`
}

type Payment struct {
	bun.BaseModel `bun:"table:payments,alias:p"`

	ID                    uint       `bun:"id,pk,autoincrement" json:"id"`
	StudentID             uint       `bun:"student_id" json:"student_id" binding:"required"`
	Amount                float64    `bun:"amount" json:"amount" binding:"required"`
	DepositApplied        float64    `bun:"deposit_applied,default:0" json:"deposit_applied"`
	Channel               string     `bun:"channel" json:"channel" binding:"required"`
	Method                string     `bun:"method" json:"method" binding:"required"`
	TransactionRef        string     `bun:"transaction_ref,unique" json:"transaction_ref"`
	ExternalOrderID       *string    `bun:"external_order_id,unique" json:"external_order_id,omitempty"`
	ExternalTransactionID *string    `bun:"external_transaction_id,unique" json:"external_transaction_id,omitempty"`
	IdempotencyKey        *string    `bun:"idempotency_key,unique" json:"idempotency_key,omitempty"`
	GatewayProvider       string     `bun:"gateway_provider" json:"gateway_provider,omitempty"`
	GatewayID             string     `bun:"gateway_id" json:"gateway_id,omitempty"`
	GatewayRawResponse    *string    `bun:"gateway_raw_response" json:"gateway_raw_response,omitempty"`
	Status                string     `bun:"status,default:pending" json:"status"`
	PaidAt                *time.Time `bun:"paid_at" json:"paid_at,omitempty"`
	IsBypassRule          bool       `bun:"is_bypass_rule" json:"is_bypass_rule"`
	CreatedBy             string     `bun:"created_by,default:SYSTEM" json:"created_by"`
	BypassReason          *string    `bun:"bypass_reason" json:"bypass_reason,omitempty"`
	ProofAttachment       *string    `bun:"proof_attachment" json:"proof_attachment,omitempty"`
	Note                  *string    `bun:"note" json:"note,omitempty"`
	IntentBillIDs         *string    `bun:"intent_bill_ids" json:"intent_bill_ids,omitempty"`
	ReconcileAttempts     int        `bun:"reconcile_attempts,default:0" json:"reconcile_attempts"`
	LastReconcileError    *string    `bun:"last_reconcile_error" json:"last_reconcile_error,omitempty"`
	LastCheckedAt         *time.Time `bun:"last_checked_at" json:"last_checked_at,omitempty"`
	ReconciledAt          *time.Time `bun:"reconciled_at" json:"reconciled_at,omitempty"`
	CreatedAt             time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt             time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt             *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
	VoidedAt              *time.Time `bun:"voided_at" json:"voided_at,omitempty"`
	ReversalOfPaymentID   *uint      `bun:"reversal_of_payment_id" json:"reversal_of_payment_id,omitempty"`

	// Details
	Details []PaymentDetail `bun:"rel:has-many,join:id=payment_id" json:"details,omitempty"`

	// Non-DB fields
	PaymentURL  string `bun:"-" json:"payment_url,omitempty"`
	RedirectURL string `bun:"-" json:"redirect_url,omitempty"`
	SnapToken   string `bun:"-" json:"snap_token,omitempty"`
	StudentName string `bun:"-" json:"student_name,omitempty"`
}

type PaymentDetail struct {
	bun.BaseModel `bun:"table:payment_details,alias:pd"`

	ID            uint      `bun:"id,pk,autoincrement" json:"id"`
	PaymentID     uint      `bun:"payment_id" json:"payment_id"`
	StudentBillID uint      `bun:"student_bill_id" json:"student_bill_id"`
	Amount        float64   `bun:"amount" json:"amount"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`

	// Join
	BillTypeName string `bun:"bill_type_name,scanonly" json:"bill_type_name,omitempty"`
	Period       string `bun:"period,scanonly" json:"period,omitempty"`
	AcademicYear string `bun:"academic_year,scanonly" json:"academic_year,omitempty"`
}

type Receipt struct {
	PaymentID     uint          `json:"payment_id"`
	ReceiptNumber string        `json:"receipt_number"`
	StudentName   string        `json:"student_name"`
	NIS           string        `json:"nis"`
	ClassName     string        `json:"class_name"`
	Amount        float64       `json:"amount"`
	PaymentMethod string        `json:"payment_method"`
	PaidAt        time.Time     `json:"paid_at"`
	Items         []ReceiptItem `json:"items"`
}

type ReceiptItem struct {
	BillName string  `json:"bill_name"`
	Period   string  `json:"period,omitempty"`
	Amount   float64 `json:"amount"`
}

type ArrearRecord struct {
	ID          uint       `bun:"id" db:"id" json:"id"`
	StudentName string     `bun:"student_name" db:"student_name" json:"student_name"`
	ClassName   string     `bun:"class_name" db:"class_name" json:"class_name"`
	BillName    string     `bun:"bill_name" db:"bill_name" json:"bill_name"`
	Period      string     `bun:"period" db:"period" json:"period"`
	Amount      float64    `bun:"amount" db:"amount" json:"amount"`
	TotalPaid   float64    `bun:"total_paid" db:"total_paid" json:"total_paid"`
	Status      string     `bun:"status" db:"status" json:"status"`
	DueDate     time.Time  `bun:"due_date" db:"due_date" json:"due_date"`
	StartDate   *time.Time `bun:"start_date" json:"start_date,omitempty"`
	EndDate     *time.Time `bun:"end_date" json:"end_date,omitempty"`
}

type CriticalBillRecord struct {
	ID           uint       `bun:"id" json:"id"`
	StudentID    uint       `bun:"student_id" json:"student_id"`
	BillTypeID   uint       `bun:"bill_type_id" json:"bill_type_id"`
	Amount       float64    `bun:"amount" json:"amount"`
	TotalPaid    float64    `bun:"total_paid" json:"total_paid"`
	DueDate      time.Time  `bun:"due_date" json:"due_date"`
	StudentName  string     `bun:"student_name" json:"student_name"`
	BillTypeName string     `bun:"bill_type_name" json:"bill_type_name"`
	ParentName   string     `bun:"parent_name" json:"parent_name"`
	ParentPhone  string     `bun:"parent_phone" json:"parent_phone"`
	StartDate    *time.Time `bun:"start_date" json:"start_date,omitempty"`
	EndDate      *time.Time `bun:"end_date" json:"end_date,omitempty"`
}

type DepositMovement struct {
	bun.BaseModel `bun:"table:deposit_movements,alias:dm"`

	ID          uint      `bun:"id,pk,autoincrement" json:"id"`
	StudentID   uint      `bun:"student_id" json:"student_id" binding:"required"`
	Type        string    `bun:"type" json:"type" binding:"required"` // IN, OUT
	Amount      float64   `bun:"amount" json:"amount" binding:"required"`
	Reason      string    `bun:"reason" json:"reason" binding:"required"` // OVERPAYMENT, BILL_VOIDED, BILL_AMOUNT_REDUCED, MANUAL_DEPOSIT, PAY_BILL, WITHDRAWAL, GATEWAY_ALLOCATION_FAILED
	ReferenceID *string   `bun:"reference_id" json:"reference_id,omitempty"`
	CreatedBy   string    `bun:"created_by" json:"created_by"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`

	// Join
	StudentName string `bun:"-" json:"student_name,omitempty"`
}
