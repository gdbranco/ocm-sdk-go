/*
Copyright (c) 2019 Red Hat, Inc.

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

package main

// This example shows how to delete a subscription.

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift-online/uhc-sdk-go"
)

func main() {
	// Create a context:
	ctx := context.Background()

	// Create a logger that has the debug level enabled:
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(true).
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't build logger: %v\n", err)
		os.Exit(1)
	}

	// Create the connection, and remember to close it:
	token := os.Getenv("UHC_TOKEN")
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		BuildContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create connection: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	// Get the client for the resource that manages the collection of subscriptions:
	collection := connection.AccountsMgmt().V1().Subscriptions()

	// Get the client for the resource that manages the subscription that we want to delete.
	// Note that this will not send any request to the server yet, so it will succeed even if
	// the subscription doesn't exist.
	resource := collection.Subscription("1BDFg66jv2kDfBh6bBog3IsZWVH")

	// Send the request to delete the subscription:
	_, err = resource.Delete().SendContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't delete subscription: %v\n", err)
		os.Exit(1)
	}
}
