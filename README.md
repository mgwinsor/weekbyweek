```
@startuml

package main {

    struct AuthAdapter <<adapter>> {
        hasher auth.PasswordHasher
        HashPassword(password string) (string, error)
    }

}

package app {

    interface "AuthService" as AppAuthService <<port>> {
        HashPassword(password string) (string, error)
    }

    struct PageData <<component>> {
        PageTitle string
        User database.User
        FormValues map[string]string
        Errors map[string]string
    }

    struct Server <<component>> {
        db database.Querier
        authService AuthService
        templates *template.Template
        SetupRoutes(mux *http.ServeMux)
        registerGet(w http.ResponseWriter, r *http.Request)
        registerPost(w http.ResponseWriter, r *http.Request)
    }
}

package auth {

    interface "PasswordHasher" as AuthPwdHasher <<port>> {
        GenerateFromPassword(password []byte, cost int) ([]byte, error)
    }

    struct "BcryptHasher" as BcryptImpl <<component>> {
        GenerateFromPassword(password []byte, cost int) ([]byte, error)
    }
}

package database {

    interface "Querier" as DBQuerier <<port>> {
        CreateUser()
        DeleteUser()
        GetUserByEmail()
        GetUserByUsername()
        UpdateUserDateOfBirth()
        UpdateUserEmail()
        UpdateUserPassword()
    }
}

AppAuthService <|.. AuthAdapter
AuthPwdHasher <|.. BcryptImpl

Server o-- DBQuerier
Server o-- AppAuthService
AuthAdapter o-- AuthPwdHasher

@enduml
```
