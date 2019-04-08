package dml

const AddUser string = "INSERT INTO zbx1.users(u_name, u_email, u_username, u_password, u_created) VALUES ($1, $2, $3, $4, $5)"

const UpdateUser string = "UPDATE zbx1.users SET u_name = $1, u_email = $2, u_username = $3, u_password = $4, u_created = $5 WHERE id = $6"

const FindAllUsers string = "SELECT id, u_name, u_email, u_username, u_created FROM zbx1.users"

const FindUserByID string = "SELECT id, u_name, u_email, u_username, u_created FROM zbx1.users WHERE id = $1"

const FindUserByUsername string = "SELECT id, u_name, u_email, u_username, u_created FROM zbx1.users WHERE u_username = $1"

const DeleteUser string = "DELETE FROM zbx1.users WHERE id = $1"

// SQL FOR ROLES, DEPENDENCIES OF USERS
const AddUserRole string = "INSERT INTO zbx1.users_roles(user_id, role_id) VALUES ($1, $1)"

const UpdateUserRole string = "UPDATE zbx1.users_roles SET role_id = $1 WHERE user_id = $2"

// SQL FOR FILES, DEPENDENCIES OF USERS
const AddUserFile string = "INSERT INTO zbx1.users_files(user_id, file_id) VALUES ($1, $2)"
