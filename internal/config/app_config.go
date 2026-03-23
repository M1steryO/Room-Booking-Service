package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	jwtSecretEnvName         = "JWT_SECRET"
	jwtTTLEnvName            = "JWT_TTL"
	slotHorizonDaysEnvName   = "SLOT_HORIZON_DAYS"
	conferenceBaseURLEnvName = "CONFERENCE_BASE_URL"
	conferenceTimeoutEnvName = "CONFERENCE_TIMEOUT"
)

type AppConfig interface {
	JWTSecret() string
	JWTTTL() time.Duration
	SlotHorizonDays() int
	ConferenceBaseURL() string
	ConferenceTimeout() time.Duration
}

type appConfig struct {
	jwtSecret         string
	jwtTTL            time.Duration
	slotHorizonDays   int
	conferenceBaseURL string
	conferenceTimeout time.Duration
}

func NewAppConfig() (AppConfig, error) {
	jwtSecret := os.Getenv(jwtSecretEnvName)
	if len(jwtSecret) == 0 {
		return nil, errors.New("jwt secret not found")
	}

	jwtTTLRaw := os.Getenv(jwtTTLEnvName)
	if len(jwtTTLRaw) == 0 {
		return nil, errors.New("jwt ttl not found")
	}
	jwtTTL, err := time.ParseDuration(jwtTTLRaw)
	if err != nil {
		return nil, errors.New("invalid jwt ttl")
	}

	slotHorizonDaysRaw := os.Getenv(slotHorizonDaysEnvName)
	if len(slotHorizonDaysRaw) == 0 {
		return nil, errors.New("slot horizon days not found")
	}
	slotHorizonDays, err := strconv.Atoi(slotHorizonDaysRaw)
	if err != nil {
		return nil, errors.New("invalid slot horizon days")
	}

	conferenceBaseURL := os.Getenv(conferenceBaseURLEnvName)
	if len(conferenceBaseURL) == 0 {
		return nil, errors.New("conference base url not found")
	}

	conferenceTimeoutRaw := os.Getenv(conferenceTimeoutEnvName)
	if len(conferenceTimeoutRaw) == 0 {
		return nil, errors.New("conference timeout not found")
	}
	conferenceTimeout, err := time.ParseDuration(conferenceTimeoutRaw)
	if err != nil {
		return nil, errors.New("invalid conference timeout")
	}

	return &appConfig{
		jwtSecret:         jwtSecret,
		jwtTTL:            jwtTTL,
		slotHorizonDays:   slotHorizonDays,
		conferenceBaseURL: conferenceBaseURL,
		conferenceTimeout: conferenceTimeout,
	}, nil
}

func (c *appConfig) JWTSecret() string {
	return c.jwtSecret
}

func (c *appConfig) JWTTTL() time.Duration {
	return c.jwtTTL
}

func (c *appConfig) SlotHorizonDays() int {
	return c.slotHorizonDays
}

func (c *appConfig) ConferenceBaseURL() string {
	return c.conferenceBaseURL
}

func (c *appConfig) ConferenceTimeout() time.Duration {
	return c.conferenceTimeout
}
