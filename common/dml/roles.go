package dml

const AddRole string = "INSERT INTO roles(r_name, r_created) VALUES (?, ?)"

const UpdateRole string = "UPDATE roles SET r_name = ?, r_created = ? WHERE id = ?"

const DeleteRole string = "DELETE from roles WHERE id = ?"

const FindRoleByID string = "SELECT id, r_name, r_created FROM roles WHERE id = (?)"

const FindAllRoles string = "SELECT id, r_name, r_created FROM roles"

const FindRoleByUserID string = "select distinct r.id, r.r_name, r.r_created from roles r inner join users_roles ur on ur.role_id = r.id where ur.user_id = (?)"
