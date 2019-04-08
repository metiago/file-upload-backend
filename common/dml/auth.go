package dml

const AuthUser string = "SELECT  u_username, u_password FROM zbx1.users WHERE u_username = $1"
