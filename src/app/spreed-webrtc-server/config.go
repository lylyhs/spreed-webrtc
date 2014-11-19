/*
 * Spreed WebRTC.
 * Copyright (C) 2013-2014 struktur AG
 *
 * This file is part of Spreed WebRTC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/strukturag/phoenix"
)

type Config struct {
	Title                  string   // Title
	ver                    string   // Version (not exported to Javascript)
	S                      string   // Static URL prefix with version
	B                      string   // Base URL
	Token                  string   // Server token
	StunURIs               []string // STUN server URIs
	TurnURIs               []string // TURN server URIs
	Tokens                 bool     // True when we got a tokens file
	Version                string   // Server version number
	UsersEnabled           bool     // Flag if users are enabled
	UsersAllowRegistration bool     // Flag if users can register
	UsersMode              string   // Users mode string
	DefaultRoomEnabled     bool     // Flag if default room ("") is enabled
	Plugin                 string   // Plugin to load
	globalRoomID           string   // Id of the global room (not exported to Javascript)
	authorizeRoomCreation  bool     // Whether a user account is required to create rooms (not exported to Javascript)
}

func NewConfig(container phoenix.Container, tokens bool) *Config {
	ver := container.GetStringDefault("app", "ver", "")

	version := container.Version()
	if version != "unreleased" {
		ver = fmt.Sprintf("%s%s", ver, strings.Replace(version, ".", "", -1))
	} else {
		ts := fmt.Sprintf("%d", time.Now().Unix())
		if ver == "" {
			ver = ts
		}
		version = fmt.Sprintf("unreleased.%s", ts)
	}

	// Read base path from config and make sure it ends with a slash.
	basePath := container.GetStringDefault("http", "basePath", "/")
	if !strings.HasSuffix(basePath, "/") {
		basePath = fmt.Sprintf("%s/", basePath)
	}
	if basePath != "/" {
		container.Printf("Using '%s' base base path.", basePath)
	}

	//TODO(longsleep): When we have a database, generate this once from random source and store it.
	serverToken := container.GetStringDefault("app", "serverToken", "i-did-not-change-the-public-token-boo")

	stunURIsString := container.GetStringDefault("app", "stunURIs", "")
	stunURIs := strings.Split(stunURIsString, " ")
	trimAndRemoveDuplicates(&stunURIs)

	turnURIsString := container.GetStringDefault("app", "turnURIs", "")
	turnURIs := strings.Split(turnURIsString, " ")
	trimAndRemoveDuplicates(&turnURIs)

	return &Config{
		Title:                  container.GetStringDefault("app", "title", "Spreed WebRTC"),
		ver:                    ver,
		S:                      fmt.Sprintf("static/ver=%s", ver),
		B:                      basePath,
		Token:                  serverToken,
		StunURIs:               stunURIs,
		TurnURIs:               turnURIs,
		Tokens:                 tokens,
		Version:                version,
		UsersEnabled:           container.GetBoolDefault("users", "enabled", false),
		UsersAllowRegistration: container.GetBoolDefault("users", "allowRegistration", false),
		UsersMode:              container.GetStringDefault("users", "mode", ""),
		DefaultRoomEnabled:     container.GetBoolDefault("app", "defaultRoomEnabled", true),
		Plugin:                 container.GetStringDefault("app", "plugin", ""),
		globalRoomID:           container.GetStringDefault("app", "globalRoom", ""),
		authorizeRoomCreation:  container.GetBoolDefault("app", "authorizeRoomCreation", false),
	}
}

func (config *Config) Get(request *http.Request) (int, interface{}, http.Header) {
	return 200, config, http.Header{"Content-Type": {"application/json; charset=utf-8"}}
}

// Helper function to clean up string arrays.
func trimAndRemoveDuplicates(data *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *data {
		x = strings.TrimSpace(x)
		if len(x) > 0 && !found[x] {
			found[x] = true
			(*data)[j] = (*data)[i]
			j++
		}
	}
	*data = (*data)[:j]
}
