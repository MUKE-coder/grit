package models

import (
	"time"

	"gorm.io/gorm"
	"saas/apps/api/internal/ids"

	"saas/apps/api/internal/crypto"
)

// SAMLKeypair is this application's service-provider identity.
//
// SAML requires the SP to have an X.509 keypair: it signs the authentication
// requests it sends, and publishes the certificate in its metadata so the IdP
// can verify them. One keypair serves every connection — it identifies the
// application, not the customer.
//
// It is generated on first use rather than asked for during setup, because
// "generate a self-signed certificate" is the step where SAML integrations
// stall. The private key is encrypted at rest with the same AES-256-GCM field
// encryption used for other secrets, and never leaves the server.
type SAMLKeypair struct {
	ID string `gorm:"primarykey;size:36" json:"id"`

	// CertPEM is public — it goes in the SP metadata the customer uploads to
	// their IdP.
	CertPEM string `gorm:"type:text;not null" json:"cert_pem"`
	// KeyPEM is the private half. Encrypted at rest, never serialized.
	KeyPEM crypto.EncryptedString `gorm:"type:text;not null" json:"-"`

	NotAfter  time.Time `json:"not_after"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (k *SAMLKeypair) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = ids.New()
	}
	return nil
}
