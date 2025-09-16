create database terraloom;

CREATE SEQUENCE public.account_id_sequence
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE public.category_id_sequence
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE public.product_id_sequence
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE public.user_id_sequence
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

CREATE TABLE public.accounts (
	id int8 DEFAULT nextval('account_id_sequence'::regclass) NOT NULL,
	display_name varchar(100) NOT NULL,
	email varchar(200) NOT NULL,
	login_password varchar(200) NOT NULL,
	registered_address varchar(200) NULL,
	is_active bool NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	username varchar(100) NULL,
	CONSTRAINT accounts_pkey PRIMARY KEY (id)
);

CREATE TABLE public.categories (
	id int8 DEFAULT nextval('category_id_sequence'::regclass) NOT NULL,
	"name" varchar(100) NOT NULL,
	description text NOT NULL,
	is_active bool NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	CONSTRAINT categories_pkey PRIMARY KEY (id)
);

CREATE TABLE public.orders (
	order_reference varchar(255) NOT NULL,
	order_date timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	status varchar(200) DEFAULT 'PENDING'::character varying NULL,
	total int8 DEFAULT 0 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	account_username varchar(200) NULL,
	delivery_address text NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (order_reference)
);

CREATE TABLE public.order_items (
	order_item_reference varchar(255) NOT NULL,
	order_reference varchar(255) NULL,
	product_id int8 NULL,
	price_snapshot int8 DEFAULT 0 NULL,
	quantity int8 DEFAULT 0 NULL,
	total int8 DEFAULT 0 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	product_name_snapshot varchar(200) NULL,
	product_image_url_snapshot varchar(200) NULL,
	CONSTRAINT order_items_pkey PRIMARY KEY (order_item_reference)
);

CREATE TABLE public.payments (
	payment_reference varchar(255) NOT NULL,
	order_reference varchar(255) NULL,
	total int8 DEFAULT 0 NULL,
	status varchar(200) DEFAULT 'PENDING'::character varying NULL,
	payment_date timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	card_number varchar(200) NULL,
	card_holder_name varchar(200) NULL,
	CONSTRAINT payments_pkey PRIMARY KEY (payment_reference)
);

CREATE TABLE public.products (
	id int8 DEFAULT nextval('product_id_sequence'::regclass) NOT NULL,
	category_id int8 NOT NULL,
	"name" varchar(100) NOT NULL,
	description text NOT NULL,
	stock int8 NOT NULL,
	price int8 NOT NULL,
	image_url text NOT NULL,
	is_active bool NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id)
);

CREATE TABLE public.users (
	id int8 DEFAULT nextval('user_id_sequence'::regclass) NOT NULL,
	"role" varchar(100) NOT NULL,
	username varchar(100) NOT NULL,
	display_name varchar(100) NOT NULL,
	email varchar(200) NOT NULL,
	login_password varchar(200) NOT NULL,
	is_active bool NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(100) NULL,
	updated_by varchar(100) NULL,
	deleted_at timestamp NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);
