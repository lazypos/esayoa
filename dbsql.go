package esayoa

var GSQLArray = []string{
	//建基础表
	`CREATE TABLE public.userbase
	(
	    id SERIAL,
	    uid character varying(50) COLLATE pg_catalog."default",
	    password character varying(50) COLLATE pg_catalog."default",
	    auth integer,
	    regtime date,
	    CONSTRAINT userbase_pkey PRIMARY KEY (id)
	)`,
	//插入超级用户
	`insert into userbase (uid,password,auth) values('adminesayoa','','99')`,
	//配置表
	`CREATE TABLE public.oaconfig
	(
	    cfgname character varying(200) NOT NULL,
	    cfgvalue character varying(4096),
	    "desc" character varying(2048),
	    PRIMARY KEY (cfgname)
	)`,
	//创建消息通知表
	`CREATE TABLE public.notifys
	(
	    notifyid SERIAL,
	    title character varying(200) NOT NULL,
	    content text,
	    sendto character varying(4096),
	    updatetime character varying(50) NOT NULL,
	    hash character varying(200),
	    PRIMARY KEY (notifyid)
	)`,
	//创建文件已读表
	`CREATE TABLE public.notifyread
	(
	    uid character varying(50) NOT NULL,
	    nfyreadids text,
	    PRIMARY KEY (uid)
	)`,
	//请假
	`CREATE TABLE public.leave
	(
	    uid character varying(50) NOT NULL,
	    createtime character varying(20) NOT NULL,
	    begintime character varying(20) NOT NULL,
	    endtime character varying(20) NOT NULL,
	    leavetype integer NOT NULL,
	    explain character varying(4096),
	    state integer
	)`,
	//用户信息
	`CREATE TABLE public.userinfo
	(
	    uid character varying(50) NOT NULL,
	    usernum character varying(50),
	    nickname character varying(100),
	    phone character varying(20),
	    email character varying(50),
	    ntyids text,
	    PRIMARY KEY (uid)
	)`,
	//创建函数
	`CREATE OR REPLACE FUNCTION CreateNotify(
	    IN titlename character varying(4096), IN article text, IN sendto character varying(4096),
	    IN fjname character varying(4096), IN fjsize character varying(200))
	RETURNS bigint
	AS $$
	    insert into notifys(title,content,updatetime,sendto, fujian, fjszie)values($1,$2,now(),$3,$4, $5);
	   	SELECT currval(pg_get_serial_sequence('notifys','notifyid'));
	$$ LANGUAGE SQL;`,
}
