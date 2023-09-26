/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/access/http"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
	"github.com/onflow/flow-go-sdk/test"
)

func main() {
	TransactionArgumentsDemo()
}

func TransactionArgumentsDemo() {
	ctx := context.Background()
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	message := test.GreetingGenerator().Random()
	greeting := cadence.String(message)

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	tx := flow.NewTransaction().
		SetScript(test.GreetingScript).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(serviceAcctAddr)

	err = tx.AddArgument(greeting)
	examples.Handle(err)

	fmt.Println("Sending transaction:")
	fmt.Println()
	fmt.Println("----------------")
	fmt.Println("Script:")
	fmt.Println(string(tx.Script))
	fmt.Println("Arguments:")
	fmt.Printf("greeting: %s\n", greeting)
	fmt.Println("----------------")
	fmt.Println()

	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	_ = examples.WaitForSeal(ctx, flowClient, tx.ID())
}
