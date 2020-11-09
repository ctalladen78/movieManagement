CREATE TABLE public.moviestbl (
	title varchar(80) NULL,
	releasedYear varchar(80) NULL,
	rating  varchar(80) NULL,
	createddate timestamp NULL,
	lastmodifieddate timestamp NULL,
	genres text NULL,
	"Id" serial NOT NULL,
	sfid varchar(200) NULL,
	CONSTRAINT pk_title PRIMARY KEY ("Id")
);