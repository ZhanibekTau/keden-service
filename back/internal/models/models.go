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
	SubscriptionStatusPending  = "pending"
	SubscriptionStatusActive   = "active"
	SubscriptionStatusExpired  = "expired"
	SubscriptionStatusRejected = "rejected"
)

type Document struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	User           *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OriginalName   string     `json:"original_name" gorm:"size:255;not null"`
	ExcelFilePath  string     `json:"excel_file_path" gorm:"size:500"`
	Status         string     `json:"status" gorm:"size:50;default:'uploaded';not null;index"`
	ErrorMessage   string     `json:"error_message,omitempty" gorm:"type:text"`
	AIResponseJSON string     `json:"-" gorm:"type:jsonb"`
	FileSize       int64      `json:"file_size"`
	QueuedAt       *time.Time `json:"queued_at"`
	ProcessedAt    *time.Time `json:"processed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (Document) TableName() string {
	return "documents"
}

const (
	DocumentStatusUploaded   = "uploaded"
	DocumentStatusQueued     = "queued"
	DocumentStatusProcessing = "processing"
	DocumentStatusCompleted  = "completed"
	DocumentStatusError      = "error"
)

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
