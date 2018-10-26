// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found at: https://github.com/gin-gonic/gin/blob/master/LICENSE

// Modified to remove the WWW-Authenticate header for uses in TempGopher

package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authPair struct {
	value string
	user  string
}

type authPairs []authPair

func (a authPairs) searchCredential(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if pair.value == authValue {
			return pair.user, true
		}
	}
	return "", false
}

// BasicAuth returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the user name and the value is the password. This does not set a www-authenticate header.
func BasicAuth(accounts gin.Accounts) gin.HandlerFunc {
	pairs := processAccounts(accounts)
	return func(c *gin.Context) {
		// Search user in the slice of allowed credentials
		user, found := pairs.searchCredential(c.GetHeader("Authorization"))
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(gin.AuthUserKey, user)
	}
}

func processAccounts(accounts gin.Accounts) authPairs {
	if len(accounts) == 0 {
		log.Panic("Empty list of authorized credentials")
	}
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts {
		if user == "" {
			log.Panic("User can not be empty")
		}
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			value: value,
			user:  user,
		})
	}
	return pairs
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}
