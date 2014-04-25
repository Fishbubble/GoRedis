### DATA STRUCTURE

	string
		+[name]string = "latermoon"
	hash
		+[info]hash = ""
		_h[info]name = "latermoon"
		_h[info]age = "27"
		_h[info]sex = "M"
	list
		+[list]list = "0,3"
		_l[list]#0 = "a"
		_l[list]#1 = "b"
		_l[list]#2 = "c"
		_l[list]#3 = "d"
	zset
		+[user_rank]zset = "3"
		_z[user_rank]m#100422 = "-2"
		_z[user_rank]m#100423 = "1"
		_z[user_rank]m#300000 = "2"
		_z[user_rank]s#-2#100422 = ""
		_z[user_rank]s#1#100423 = ""
		_z[user_rank]s#2#300000 = ""

	table
		+[users]table = "idx:[momoid,regtime]"
		_t[users]rows = "100"
		_t[users]id#00001 = {momoid:"100422", age:12, sex:"M", regtime:1397801227}
		_t[users]idx#momoid#100422 = "00001"
		_t[users]idx#regtime#1397801227 = "00001"

### TABLE

	table_exec users "select * from users where momoid='100422'"