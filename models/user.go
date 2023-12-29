package models

import "time"

// TODO: maybe save ratings / skill level with User in the future
// instead of them making a new form every time.
// (saved configs)

type User struct {
	ID            string `gorm:"not null; unique" json:"id"`
	Username      string `gorm:"not null" json:"username"`
	Discriminator string
	GlobalName    string `gorm:"not null; unique" json:"global_name"`
	Avatar        string
	Email         string
	Locale        string
	CreatedAt     time.Time `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"default:current_timestamp"`
}

// "{\"id\":\"313823868450111498\",
// \"username\":\"wolffy._.\",
// \"avatar\":\"472c69fdac76ec079d7606c6861ce310\",
// \"discriminator\":\"0\",
// \"public_flags\":0,
// \"premium_type\":0,
// \"flags\":0,
// \"banner\":null,
// \"accent_color\":null,
// \"global_name\":\"djwolff\",
// \"avatar_decoration_data\":null,
// \"banner_color\":null,
// \"mfa_enabled\":true,
// \"locale\":\"en-US\"}\n"
