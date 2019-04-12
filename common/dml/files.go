package dml

const AddFile string = "INSERT INTO zbx1.files(f_name, f_ext, f_data, f_created) VALUES ($1, $2, $3, $4) returning id"

const FindAllFiles string = "SELECT id, f_name, f_ext, f_created, f_data zbx1.FROM files"

const FindFileByID string = "SELECT id, f_name, f_ext, f_data FROM zbx1.files WHERE id = ($1)"

const FindAllFilesByUsername string = "SELECT f.id, f.f_name, f.f_ext, f.f_created, f.f_data, u.u_name FROM zbx1.files f INNER JOIN zbx1.users_files uf ON f.id = uf.file_id INNER JOIN zbx1.users u ON u.id = uf.user_id WHERE u.u_username = ($1)"
