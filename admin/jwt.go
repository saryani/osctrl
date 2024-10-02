package main

import (
	"crypto/rsa"
	"crypto/tls"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmpsec/osctrl/settings"
	"github.com/jmpsec/osctrl/types"
	"github.com/spf13/viper"
	"github.com/rs/zerolog/log"
)

// Function to load the configuration file
func loadJWTConfiguration(file string) (types.JSONConfigurationJWT, error) {
	var cfg types.JSONConfigurationJWT
	log.Info().Msgf("Loading %s", file)
	// Load file and read config
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	// JWT values
	jwtRaw := viper.Sub(settings.AuthJWT)
	if err := jwtRaw.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	// No errors!
	return cfg, nil
}

// Helper to parse JWT tokens because the SAML library is total garbage
func parseJWTFromCookie(keypair tls.Certificate, cookie string) (JWTData, error) {
	type TokenClaims struct {
		jwt.StandardClaims
		Attributes map[string][]string `json:"attr"`
	}
	tokenClaims := TokenClaims{}
	token, err := jwt.ParseWithClaims(cookie, &tokenClaims, func(t *jwt.Token) (interface{}, error) {
		return keypair.PrivateKey.(*rsa.PrivateKey).Public(), nil
	})

	if err != nil || !token.Valid {
		return JWTData{}, err
	}
	return JWTData{
		Subject:  tokenClaims.Subject,
		Email:    tokenClaims.Subject,
		Username: tokenClaims.Subject,
	}, nil
}
