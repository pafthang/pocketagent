package spaceapis

import "log/slog"

func logInviteLink(log *slog.Logger, publicBaseURL, token string) {
	if log == nil {
		return
	}
	log.Info("invite link (configure SMTP for production)",
		"accept_url", publicBaseURL+"/invites/accept",
		"preview_url", publicBaseURL+"/invites/"+token,
		"token", token,
	)
}

func logVerificationLink(log *slog.Logger, publicBaseURL, token string) {
	if log == nil {
		return
	}
	log.Info("email verification link (configure SMTP for production)",
		"verify_url", publicBaseURL+"/auth/verify-email",
		"token", token,
	)
}