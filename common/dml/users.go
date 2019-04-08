package dml

const AddUser string = "INSERT INTO users(u_name, u_email, u_username, u_password, u_created) VALUES (?,?,?,?, ?)"

const UpdateUser string = "UPDATE users SET u_name = ?, u_email = ?, u_username = ?, u_password = ?, u_created = ? WHERE id = ?"

const FindAllUsers string = "SELECT id, u_name, u_email, u_username, u_created FROM users"

const FindUserByID string = "SELECT id, u_name, u_email, u_username, u_created FROM users WHERE id = ?"

const FindUserByUsername string = "SELECT id, u_name, u_email, u_username, u_created FROM users WHERE u_username = ?"

const DeleteUser string = "DELETE FROM users WHERE id = ?"

// SQL FOR ROLES, DEPENDENCIES OF USERS
const AddUserRole string = "INSERT INTO users_roles(user_id, role_id) VALUES (?, ?)"

const UpdateUserRole string = "UPDATE users_roles SET role_id = ? WHERE user_id = ?"

// SQL FOR FILES, DEPENDENCIES OF USERS
const AddUserFile string = "INSERT INTO users_files(user_id, file_id) VALUES (?, ?)"
