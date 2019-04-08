package dml

const AddFile string = "INSERT INTO files(f_name, f_ext, f_data, f_created) VALUES (?,?,?,?)"

const FindAllFiles string = "SELECT id, f_name, f_ext, f_created, f_data FROM files"

const FindFileByID string = "SELECT id, f_name, f_ext, f_data FROM files WHERE id = (?)"

const FindAllFilesByUsername string = "SELECT f.id, f.f_name, f.f_ext, f.f_created, f.f_data FROM files f INNER JOIN users_files uf ON f.id = uf.file_id INNER JOIN users u ON u.id = uf.user_id WHERE u.u_username = (?)"
