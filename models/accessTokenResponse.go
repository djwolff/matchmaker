package models

type DiscordUser struct {
	Id            string
	Username      string
	Discriminator string
	GlobalName    string
	Avatar        string
	Email         string
	Locale        string
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
