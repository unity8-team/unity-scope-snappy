/* Copyright (C) 2015 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/scope"
)

// main is the entry point of the scope.
//
// Supported environment variables:
// - WEBDM_URL: address[:port] on which WebDM is listening
func main() {

	// TODO: HACK HACK HACK: Work around for bug #1630370
	// ("runtime error: cgo argument has Go pointer to Go pointer"")
	if os.Getenv("GODEBUG") != "cgocheck=0" {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		env := os.Environ()
		env = append(env, "GODEBUG=cgocheck=0")
		cmd.Env = env
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		fmt.Println(err)
		os.Exit(0)
	}

	scope, err := scope.New()
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to create scope: %s", err)
		return
	}

	err = scopes.Run(scope)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to run scope: %s", err)
	}
}
