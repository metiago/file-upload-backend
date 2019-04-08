package dml

const AddRole string = "INSERT INTO zbx1.roles(r_name, r_created) VALUES ($1, $2)"

const UpdateRole string = "UPDATE zbx1.roles SET r_name = $1, r_created = $2 WHERE id = $3"

const DeleteRole string = "DELETE from zbx1.roles WHERE id = $1"

const FindRoleByID string = "SELECT id, r_name, r_created FROM zbx1.oles WHERE id = ($1)"

const FindAllRoles string = "SELECT id, r_name, r_created FROM zbx1.roles"

const FindRoleByUserID string = "select distinct r.id, r.r_name, r.r_created from zbx1.roles r inner join zbx1.users_roles ur on ur.role_id = r.id where ur.user_id = ($1)"
