package models

import "time"

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex;size:50;not null"`
}

func (Role) TableName() string {
	return "roles"
}

type User struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Email        string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string     `json:"-" gorm:"size:255;not null"`
	FirstName    string     `json:"first_name" gorm:"size:255;not null"`
	LastName     string     `json:"last_name" gorm:"size:255;not null"`
	Phone        string     `json:"phone" gorm:"size:50;not null"`
	RoleID       uint       `json:"role_id" gorm:"not null"`
	Role         *Role      `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	AccountType  string     `json:"account_type" gorm:"size:20;not null"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type Company struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	Email         string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CompanyName   string    `json:"company_name" gorm:"size:255;not null"`
	LegalName     string    `json:"legal_name" gorm:"size:255;not null"`
	BIN           string    `json:"bin" gorm:"uniqueIndex;size:12;not null"`
	ContactPerson string    `json:"contact_person" gorm:"size:255"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}

type Subscription struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	UserID       uint       `json:"user_id" gorm:"not null;index"`
	User         *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Status       string     `json:"status" gorm:"size:50;default:'pending';not null;index"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Amount       float64    `json:"amount" gorm:"type:decimal(10,2);default:12990.00"`
	AdminComment string     `json:"admin_comment" gorm:"type:text"`
	RequestedAt  time.Time  `json:"requested_at"`
	ApprovedAt   *time.Time `json:"approved_at"`
	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedBy   *User      `json:"approved_by,omitempty" gorm:"foreignKey:ApprovedByID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

const (
	SubscriptionStatusPending     = "pending"
	SubscriptionStatusInProgress  = "in_progress"
	SubscriptionStatusInvoiceSent = "invoice_sent"
	SubscriptionStatusActive      = "active"
	SubscriptionStatusExpired     = "expired"
	SubscriptionStatusRejected    = "rejected"
)

type Document struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	UserID       uint       `json:"user_id" gorm:"not null;index"`
	User         *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OriginalName string     `json:"original_name" gorm:"size:255;not null"`
	Status       string     `json:"status" gorm:"size:50;default:'uploaded';not null;index"`
	ErrorMessage string     `json:"error_message,omitempty" gorm:"type:text"`
	FileSize     int64      `json:"file_size"`
	ProcessedAt  *time.Time `json:"processed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Document) TableName() string {
	return "documents"
}

const (
	DocumentStatusUploaded  = "uploaded"
	DocumentStatusCompleted = "completed"
	DocumentStatusError     = "error"
)

// DocumentFields stores AI-extracted header fields for a document.
// Column names match exactly the Excel export fields.
// One row per document (uniqueIndex on document_id).
type DocumentFields struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	DocumentID        uint      `json:"document_id" gorm:"not null;uniqueIndex"`
	DocumentType      string    `json:"document_type" gorm:"size:100"`
	DeclarationNumber string    `json:"declaration_number" gorm:"size:255"`
	Date              string    `json:"date" gorm:"size:50"`
	Sender            string    `json:"sender" gorm:"size:500"`
	Receiver          string    `json:"receiver" gorm:"size:500"`
	CountryOrigin     string    `json:"country_origin" gorm:"size:100"`
	CountryDest       string    `json:"country_dest" gorm:"size:100"`
	Currency          string    `json:"currency" gorm:"size:10"`
	TotalValue        float64   `json:"total_value" gorm:"type:decimal(15,2)"`
	CustomsValue      float64   `json:"customs_value" gorm:"type:decimal(15,2)"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (DocumentFields) TableName() string {
	return "document_fields"
}

// DocumentItem stores one goods line for a document.
// Column names match exactly the Excel export columns.
// Multiple rows per document.
type DocumentItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DocumentID  uint      `json:"document_id" gorm:"not null;index"`
	Number      int       `json:"number"`
	HSCode      string    `json:"hs_code" gorm:"column:hs_code;size:20"`
	Description string    `json:"description" gorm:"type:text"`
	Quantity    float64   `json:"quantity" gorm:"type:decimal(15,3)"`
	Unit        string    `json:"unit" gorm:"size:50"`
	WeightNet   float64   `json:"weight_net" gorm:"type:decimal(15,3)"`
	WeightGross float64   `json:"weight_gross" gorm:"type:decimal(15,3)"`
	Value       float64   `json:"value" gorm:"type:decimal(15,2)"`
	DutyRate    string    `json:"duty_rate" gorm:"size:20"`
	VATRate     string    `json:"vat_rate" gorm:"column:vat_rate;size:20"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DocumentItem) TableName() string {
	return "document_items"
}

type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:500;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
