package dml

const AddUser string = "INSERT INTO zbx1.users(u_name, u_email, u_username, u_password, u_created) VALUES ($1, $2, $3, $4, $5) returning id"

const UpdateUser string = "UPDATE zbx1.users SET u_name = $1, u_email = $2, u_username = $3, u_created = $4 WHERE id = $5"

const FindAllUsers string = "SELECT id, u_name, u_email, u_username, u_created FROM zbx1.users"

const FindUserByID string = "SELECT id, u_name, u_email, u_username, u_password, u_created FROM zbx1.users WHERE id = $1"

const FindUserByUsername string = "SELECT id, u_name, u_email, u_username, u_password, u_created FROM zbx1.users WHERE u_username = $1"

const DeleteUser string = "DELETE FROM zbx1.users WHERE id = $1"

const UpdateUserPassword string = "UPDATE zbx1.users SET u_password= $1, u_created = $2 WHERE id = $3"

const AddUserFile string = "INSERT INTO zbx1.users_files(user_id, file_id) VALUES ($1, $2)"
