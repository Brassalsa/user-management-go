# User Management System

## Introduction

User Management system is a RESTful api, made with go lang, mongoDB and more, that allows user to make an account with their email, username, name and password, also upload their profile image when account is created.

## Technolgies used

### net/http

Go lang's net/http is a standard liberary package that is used to make a server.
[Docs: net/http](https://pkg.go.dev/net/http)

### Mongo DB and mongoDB driver for go

[Mongo DB](https://www.mongodb.com/) for storing data.

### Authentication JWT

[JWT](https://jwt.io/) for token generation.

## About

### User

Users can change their name, password and profile image, enter their email.

### Admin

Admin users can view, modify and delete users. They can also change their and other user's email and username also they can promote normal users to admin. Admin users cannot change other admin details.

## API's

### Health

- Health check (get): /healthz
- Error (get): /err

### User: /api/v1/users

- Register user (post): /register
- Login (post): /login
- Account Details (get): /
- Account update (put): /
- Change Avatar (post): /users/avatar
- Change Password (post): /password
- Delete Account (delete): /delete

### Admin: /api/v1/admin

- Create admin account (post): /register
- Delete user (delete): /user/{userID}
