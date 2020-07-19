<h1 align="center">Terra<i>pukka</i><br>
<img src="https://img.shields.io/github/languages/code-size/svasandani/terrapukka" />
<img src="https://img.shields.io/github/license/svasandani/terrapukka" />
<img src="https://img.shields.io/github/last-commit/svasandani/terrapukka" />
<img src="https://img.shields.io/github/go-mod/go-version/svasandani/terrapukka" />
<br>
</h1>
<br>
A Go OAuth provider for TerraLing. Currently in development.

## Dependencies
- [go](https://golang.org)
- [go-mysql-driver](https://github.com/go-sql-driver/mysql)
- [mysql](https://mysql.com)

## Installation
- Download and install MySQL
- Make sure Git is installed
- Clone the repo

  `$ git clone https://github.com/svasandani/terrapukka`

- That's literally it

## Testing
No tests have currently been written. See [#3](https://github.com/svasandani/terrapukka/issues/3).

## Endpoints
The service currently has endpoints for registering `Clients` and `Users`, authorizing `Users`, and granting `User` data access (in this case, their names and emails) to `Clients`. The endpoints are:

### Client

- `api/register`

  Register a new `user`. Takes in the following structure, with fields required as marked:
  ```
  {

    "name": user's name (e.g. John Doe), required,

    "email": user's email (e.g. jd@example.com), required,

    "password": user's password (at least 8 characters), required

  }
  ```

  Returns the following:
  ```
  {

    "redirect_uri": URI to redirect the user to,

    "auth_code": temporary authorization code,

    "state": state given by client at redirect time

  }
  ```

- `api/client/register`

  Register a new `client`. Takes in the following structure, with fields required as marked:
  ```
  {

    "name": client's name (e.g. Terraling), required,

    "redirect_uri": URI to redirect the user to after successful authentication, required

  }
  ```

  Returns the following:
  ```
  {

    "name": client's registered name,

    "id": client's id, generated by the application,

    "secret": client's secret, generated by the application,

    "redirect_uri": client's registered redirect_uri

  }
  ```

- `api/auth`

  Authorize a `user`. Takes in the following structure, with fields required as marked:
  ```
  {

    "response_type": type of authorization request, usually "code", required,

    "client_id": client's ID, returned at registration, required,

    "redirect_uri": URI to redirect the user to after successful authentication, must match registered URI, required,

    "user": user model containing email and password, required {

      "email": user's email, required,

      "password": user's password, required

    }

    "state": random token generated by client, expect to match response state, optional

  }
  ```

  Returns the following:
  ```
  {

    "redirect_uri": URI to redirect the user to,

    "auth_code": temporary authorization code,

    "state": state given by client at redirect time

  }
  ```

- `api/client/auth`

  Authenticate a `client` attempting to access `user` data. Takes in the following structure, with fields required as marked:
  ```
  {

    "grant_type": type of data request, usually "identity", required,

    "auth_code": user's temporary authorization code, returned from user authorization, required,

    "client": client model containing id and secret, required {

      "id": client's id, required,

      "secret": client's secret, required

      "redirect_uri": client's redirect_uri, must match registered URI, required

    }

  }
  ```

  Returns the following:
  ```
  {

    "user": requested user data {

      "name": user's name,

      "email": user's email

    }

  }
  ```

## Contributing
Look through the issues and read through the code to see what needs help. Some tags:
- `@TODO` - problems that are attached to issues.
- `@QOL` - problems that aren't major and so may not be attached to issues.
