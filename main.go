package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

var auClient *msgraph.AdministrativeUnitsClient
var groupsClient *msgraph.GroupsClient
var au *msgraph.AdministrativeUnit

func init() {
	tenantId := os.Getenv("AZURE_TENANT_ID")
	var status int
	var err error
	ctx := context.Background()
	authConfig := &auth.Config{
		Environment:         environments.Global,
		TenantID:            tenantId,
		EnableAzureCliToken: true,
	}

	authorizer, err := authConfig.NewAuthorizer(ctx, environments.MsGraphGlobal)
	if err != nil {
		log.Fatal(err)
	}

	auClient = msgraph.NewAdministrativeUnitsClient(tenantId)
	groupsClient = msgraph.NewGroupsClient(tenantId)

	auClient.BaseClient.Authorizer = authorizer
	groupsClient.BaseClient.Authorizer = authorizer

	administrativeUnitName := "Test AU"
	administrativeUnit := msgraph.AdministrativeUnit{DisplayName: &administrativeUnitName}
	au, status, err = auClient.Create(context.Background(), administrativeUnit)
	if err != nil {
		log.Fatalf("error creating administrative unit: %s", err)
	}
	if status != 201 {
		log.Fatalf("error creating administrative unit, got status code %d", status)
	}
}

func main() {
	startTime := time.Now()
	groupName := "TestGroupInAu"
	truePtr := true

	groupName = "TestGroupOutsideOfAU2"
	group, status, err := groupsClient.Create(context.Background(), msgraph.Group{
		DisplayName:     &groupName,
		MailNickname:    &groupName,
		MailEnabled:     &truePtr,
		SecurityEnabled: &truePtr,
		GroupTypes:      &[]string{"Unified"},
	})
	if err != nil {
		log.Fatalf("error creating group: %s", err)
	}
	if status != 201 {
		log.Fatalf("error creating group, got status code %d", status)
	}
	endTime := time.Now()
	log.Println("Time taken to create group outside of AU: ", endTime.Sub(startTime))
	log.Printf("Group ID: %s, Group name: %s", *group.ID(), *group.DisplayName)
	startTime = time.Now()
	for {
		status := getGroup(group)
		if status != 404 {
			break
		}
		log.Printf("Group not found yet. Time elapsed: %s", time.Since(startTime))
		time.Sleep(5 * time.Second)
	}
	endTime = time.Now()
	log.Println("Time taken to successfully GET group created outside of an administrative unit: ", endTime.Sub(startTime))

	startTime = time.Now()
	groupName = "TestGroupInAU"
	group, status, err = auClient.CreateGroup(context.Background(), *au.ID, &msgraph.Group{
		DisplayName:     &groupName,
		MailNickname:    &groupName,
		MailEnabled:     &truePtr,
		SecurityEnabled: &truePtr,
		GroupTypes:      &[]string{"Unified"},
	})
	if err != nil {
		log.Fatalf("error creating group: %s", err)
	}
	if status != 201 {
		log.Fatalf("error creating group, got status code %d", status)
	}
	endTime = time.Now()
	log.Println("Time taken to create group in AU: ", endTime.Sub(startTime))
	log.Printf("Group ID: %s, Group name: %s", *group.ID(), *group.DisplayName)
	startTime = time.Now()
	for {
		status := getGroup(group)
		if status != 404 {
			break
		}
		log.Printf("Group not found yet. Time elapsed: %s", time.Since(startTime))
		time.Sleep(5 * time.Second)
	}
	endTime = time.Now()
	log.Println("Time taken to successfully GET group when creating it inside of an administrative unit: ", endTime.Sub(startTime))
}

func getGroup(group *msgraph.Group) int {
	_, status, err := groupsClient.Get(context.Background(), *group.ID(), odata.Query{Select: []string{"allowExternalSenders", "autoSubscribeNewMembers", "hideFromAddressLists", "hideFromOutlookClients"}})
	if err != nil {
		log.Println(err)
	}
	return status
}
