# IAM

* Permission
  * CREATE
  * UPDATE
  * DELETE
  * READ
* Roles
  * Administrator
    * Permission
      * CREATE
      * UPDATE
      * DELETE
      * READ
  * Reader
    * Permission
      * READ
* Organization

```
docker run --name mysql -v $(pwd)/Documents/github/mysql:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=admin -p 3306:3306 -d mysql:latest
```

## Flow

* User root signup flow
  * Email
  * Password
  * Company
* Normal user signin
  * Check if it's a first login
* Add user to the Administrator group
* Invite user
  * Email
  * Default role will be reader
  * Password
* Add user to group

## Root signup flow - route /signup

* User will register the account
  * Email
  * Password
  * Company
* Use case:
  * user signup with different root credential and different company
    * Create company and user
  * user signup with same root credential and same company
    * Error - Company already exist
  * user signup with same root credential and different company
    * Create company and user
  * user signup with different root credential and same company
    * Error - Company already exist

## Signin user flow - /signin

* Signin require
  * Email
  * Password
* Return JWT

## Invite user flow - /invite/users

* Signin
* Get JWT token
  * Check permission
  * User should send list of user and role in the request payload
