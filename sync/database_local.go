package sync

func (database *Database) localMysqldumpCmdBuilder(additionalArgs []string, useFilter bool) []interface{} {
	var args []string

	if database.Local.User != "" {
		args = append(args, "-u" + database.Local.User)
	}

	if database.Local.Password != "" {
		args = append(args, "-p" + database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h" + database.Local.Hostname)
	}

	if database.Local.Port != "" {
		args = append(args, "-P" + database.Local.Port)
	}

	if len(args) > 0 {
		args = append(args, additionalArgs...)
	}

	// exclude
	excludeArgs, includeArgs := database.mysqlTableFilter(&database.Local.Connection, "local");
	if useFilter && len(excludeArgs) > 0 {
		args = append(args, excludeArgs...)
	}

	// schema
	args = append(args, database.Local.Schema)

	// include
	if useFilter && len(includeArgs) > 0 {
		args = append(args, includeArgs...)
	}

	return database.Local.Connection.CommandBuilder("mysqldump", args...)
}

func (database *Database) localMysqlCmdBuilder(args ...string) []interface{} {
	args = append(args, "-BN")

	if database.Local.User != "" {
		args = append(args, "-u" + database.Local.User)
	}

	if database.Local.Password != "" {
		args = append(args, "-p" + database.Local.Password)
	}

	if database.Local.Hostname != "" {
		args = append(args, "-h" + database.Local.Hostname)
	}

	if database.Local.Port != "" {
		args = append(args, "-P" + database.Local.Port)
	}

	if database.Local.Schema != "" {
		args = append(args, database.Local.Schema)
	}

	return database.Local.Connection.CommandBuilder("mysql", args...)
}

