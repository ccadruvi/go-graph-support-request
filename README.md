# go-graph-support-request

This repository reproduces an assumed bug with the Azure Graph API a small go program.

## How to use

Run `az login --allow-no-subscriptions` on an Azure tenant with an active Azure AD Premium P1 license and log in with a user that has the permissions to create groups and administrative units in Azure AD.

Set an environment variable with the tenant ID:

```sh
export AZURE_TENANT_ID="..."
```

And run the go program:

```sh
go run main.go
```

## What the program does

1. An administrative unit named `Test AU` is created.
2. A group named `TestGroupOutsideOfAU2` is created at tenant level, without the assigmnet of any administrative unit.
3. A group named `TestGroupInAU` is created inside the administrative unit `Test AU`.

## The assumed bug

Creating a group inside an administrative unit and then performing a `GET` request on the group with the `$select=allowExternalSenders,autoSubscribeNewMembers,hideFromAddressLists,hideFromOutlookClients` query parameter, it takes a lot longer than when creating the group outside of an administrative unit.

Example:

```txt
Time taken to successfully GET group created outside of an administrative unit:  32.156557895s
...
Time taken to successfully GET group when creating it inside of an administrative unit: 11m24.011496993s
```
