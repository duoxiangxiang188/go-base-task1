package models

//MigrateAll 迁移所有数据表

func MigrateAll() error {
	//按顺序迁移（外键依赖：posts依赖users，comments依赖users和posts）
	if err := MigrateUsers(); err != nil {
		return err
	}
	if err := MigratePosts(); err != nil {
		return err
	}
	if err := MigrateComments(); err != nil {
		return err
	}
	return nil
}
