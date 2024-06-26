/*
Copyright (c) 2022 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This example shows how to retrieve the collection of follow-ups
// from the web-rca project.

package main

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	web "github.com/openshift-online/ocm-sdk-go/webrca/v1"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("OCM_TOKEN")

	connection, err := sdk.NewConnectionBuilder().URL("http://localhost:8000").Tokens(token).BuildContext(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build connection: %v\n", err)
		os.Exit(1)
	}

	defer connection.Close()

	incident_id := "bc0760a1-7f40-4b68-86b1-3f8f2a5b0f59" // Adjust as needed
	collection := connection.WebRCA().V1().Incidents().Incident(incident_id).FollowUps()

	size := 10
	page := 1

	for {
		response, err := collection.List().
			//OrderBy("created_at").
			Size(size).
			Page(page).
			SendContext(ctx)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't retrieve page %d: %s\n", page, err)
			os.Exit(1)
		}

		response.Items().Each(func(fu *web.FollowUp) bool {
			fmt.Printf("%s - %s - %s\n", fu.Incident().ID(), fu.Url(), fu.Status())
			return true
		})

		// Break the loop if the size of the page is less than requested, otherwise go to
		// the next page:
		if response.Size() < size {
			break
		}

		page++
	}
}
