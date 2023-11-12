CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE public.grupo_prd (
	id uuid DEFAULT uuid_generate_v4 (),
	nome varchar(20) NOT NULL,
	
	CONSTRAINT grupo_prd_pkey PRIMARY KEY (id)
);

CREATE TABLE public.sbu_grupo_prd (
	id uuid DEFAULT uuid_generate_v4 (),
	nome_sub_grupo varchar(40) NOT NULL,
	id_grupo_prd uuid,
	
	CONSTRAINT sub_grupo_prd_pkey PRIMARY KEY (id)
);



CREATE TABLE usuarios_user(
	id uuid DEFAULT uuid_generate_v4 (),
	user_name varchar(40) not null,
	email varchar(255) not null, 
	passwrod varchar(255) not null,
	hashpass varchar(150) not null,
	enable bool,
	locked bool,
	created_at time without time zone NOT NULL DEFAULT NOW(),
    updated_at time without time zone NOT NULL DEFAULT NOW(),
	UNIQUE(email)
);

-- public.codigos_promocionais_codigos_promocionais definition

CREATE TABLE public.empresa (
	id uuid DEFAULT uuid_generate_v4 () NOT NULL,
	cnpj varchar(30) NULL,
	nome varchar(50) NULL,
	CEP   varchar(50) NULL
	Endereco     varchar(50) NULL
	Complemento   varchar(50) NULL
	Numero   varchar(50) NULL
	Pais     varchar(50) NULL
	Estado   varchar(50) NULL
	Cidade    varchar(50) NULL
	Bairro     varchar(50) NULL
	CONSTRAINT codigos_promocionais_codigos_promocionais_pkey PRIMARY KEY (id)
);


-- public.codigos_promocionais_generatepromocionalcode definition

CREATE TABLE public.codigos_promocionais_generatepromocionalcode (
	id serial4 NOT NULL,
	quantidade int2 NOT NULL,
	data_inicio date NOT NULL,
	valor numeric(11, 2) NOT NULL,
	data_fim date NOT NULL,
	nome varchar(255) NULL,
	cnpjcpf varchar(18) NULL,
	email varchar(255) NULL,
	tipo varchar(1) NULL,
	"generated" bool NOT NULL,
	id_servico int4 NULL,
	CONSTRAINT codigos_promocionais_generatepromocionalcode_id_servico_check CHECK ((id_servico >= 0)),
	CONSTRAINT codigos_promocionais_generatepromocionalcode_pkey PRIMARY KEY (id),
	CONSTRAINT codigos_promocionais_generatepromocionalcode_quantidade_check CHECK ((quantidade >= 0))
);


-- public.extensoes_extensoes definition

CREATE TABLE public.extensoes_extensoes (
	id serial4 NOT NULL,
	codextensao int4 NULL,
	nome varchar(5) NOT NULL,
	CONSTRAINT extensoes_extensoes_codextensao_check CHECK ((codextensao >= 0)),
	CONSTRAINT extensoes_extensoes_nome_6b49190f_uniq UNIQUE (nome),
	CONSTRAINT extensoes_extensoes_pkey PRIMARY KEY (id)
);
CREATE INDEX extensoes_extensoes_nome_6b49190f_like ON public.extensoes_extensoes USING btree (nome varchar_pattern_ops);


-- public.indicacoes_indicacaocliente definition

CREATE TABLE public.indicacoes_indicacaocliente (
	id serial4 NOT NULL,
	code_pai varchar(200) NOT NULL,
	code_filho varchar(200) NOT NULL,
	create_at timestamptz NOT NULL,
	email_filho varchar(254) NOT NULL,
	email_pai varchar(254) NOT NULL,
	CONSTRAINT indicacoes_indicacaocliente_pkey PRIMARY KEY (id)
);


-- public.indicacoes_indicacoes definition

CREATE TABLE public.indicacoes_indicacoes (
	id serial4 NOT NULL,
	codindicacao int4 NULL,
	nome varchar(50) NOT NULL,
	percentual_promocional numeric(11, 2) NULL,
	CONSTRAINT indicacoes_indicacoes_codindicacao_check CHECK ((codindicacao >= 0)),
	CONSTRAINT indicacoes_indicacoes_pkey PRIMARY KEY (id)
);


-- public.registros_coautores_registros_coautores definition

CREATE TABLE public.registros_coautores_registros_coautores (
	id serial4 NOT NULL,
	codcoautor int4 NULL,
	codregistro int4 NOT NULL,
	nome varchar(255) NOT NULL,
	documento varchar(100) NOT NULL,
	percentual_obra numeric(5, 2) NOT NULL,
	CONSTRAINT registros_coautores_registros_coautores_codcoautor_check CHECK ((codcoautor >= 0)),
	CONSTRAINT registros_coautores_registros_coautores_codregistro_check CHECK ((codregistro >= 0)),
	CONSTRAINT registros_coautores_registros_coautores_pkey PRIMARY KEY (id)
);


-- public.servicos_servicos definition

CREATE TABLE public.servicos_servicos (
	id serial4 NOT NULL,
	codservico int4 NULL,
	nome varchar(100) NOT NULL,
	preco numeric(11, 2) NOT NULL,
	tamanho int4 NOT NULL,
	medida varchar(5) NULL,
	servico_digitalizacao bool NOT NULL,
	CONSTRAINT servicos_servicos_codservico_check CHECK ((codservico >= 0)),
	CONSTRAINT servicos_servicos_pkey PRIMARY KEY (id),
	CONSTRAINT servicos_servicos_tamanho_check CHECK ((tamanho >= 0))
);


-- public.usuarios_confuguracao definition

CREATE TABLE public.usuarios_confuguracao (
	id serial4 NOT NULL,
	dolar numeric(5, 2) NOT NULL,
	euro numeric(5, 2) NOT NULL,
	credito_inicial numeric(5, 2) NOT NULL,
	atualização timestamptz NOT NULL,
	valor_file numeric(5, 2) NOT NULL,
	pontuacao numeric(5, 2) NOT NULL,
	percentual int2 NOT NULL,
	CONSTRAINT usuarios_confuguracao_percentual_check CHECK ((percentual >= 0)),
	CONSTRAINT usuarios_confuguracao_pkey PRIMARY KEY (id)
);


-- public.clientes_clientes definition

CREATE TABLE public.clientes_clientes (
	id serial4 NOT NULL,
	id_cliente int4 NULL,
	nome varchar(255) NOT NULL,
	valor_credito numeric(11, 2) NOT NULL,
	telefone varchar(16) NULL,
	celular varchar(16) NOT NULL,
	data_nascimento date NOT NULL,
	sexo varchar(1) NULL,
	tipo_pessoa varchar(1) NOT NULL,
	nome_mae varchar(100) NULL,
	nome_pai varchar(100) NULL,
	cnpjcpf varchar(18) NOT NULL,
	cep varchar(9) NOT NULL,
	endereco varchar(255) NOT NULL,
	complemento varchar(255) NULL,
	numero varchar(8) NOT NULL,
	pais varchar(2) NOT NULL,
	estado varchar(50) NOT NULL,
	cidade varchar(50) NOT NULL,
	bairro varchar(50) NOT NULL,
	documento_identidade varchar(50) NULL,
	documento_tipo varchar(20) NULL,
	passaporte varchar(50) NULL,
	nacionalidade varchar(20) NOT NULL,
	estadocivil varchar(1) NULL,
	biografia text NULL,
	nif varchar(100) NULL,
	facebook varchar(100) NULL,
	twitter varchar(100) NULL,
	homepage varchar(100) NULL,
	data_cadastro timestamptz NOT NULL,
	confirmation_key varchar(80) NULL,
	id_usuario_id int4 NOT NULL,
	codindicacao_id int4 NOT NULL,
	codigo_promocional varchar(200) NULL,
	link_indicacao varchar(200) NULL,
	meu_link_indicacao varchar(200) NULL,
	pontuacao numeric(11, 2) NULL,
	permitir_enviar_informativo bool NULL DEFAULT false,
	CONSTRAINT clientes_clientes_codusuario_id_key UNIQUE (id_usuario_id),
	CONSTRAINT clientes_clientes_id_cliente_108b1869_check CHECK ((id_cliente >= 0)),
	CONSTRAINT clientes_clientes_pkey PRIMARY KEY (id),
	CONSTRAINT clientes_clientes_codindicacao_id_d64d533d_fk_indicacoe FOREIGN KEY (codindicacao_id) REFERENCES public.indicacoes_indicacoes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT clientes_clientes_id_usuario_id_94d171d4_fk_usuarios_user_id FOREIGN KEY (id_usuario_id) REFERENCES public.usuarios_user(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX clientes_clientes_codindicacao_id_d64d533d ON public.clientes_clientes USING btree (codindicacao_id);


-- public.clientes_clientes_atuacao definition

CREATE TABLE public.clientes_clientes_atuacao (
	id serial4 NOT NULL,
	clientes_id int4 NOT NULL,
	servicos_id int4 NOT NULL,
	CONSTRAINT clientes_clientes_atuacao_clientes_id_servicos_id_0e1c047d_uniq UNIQUE (clientes_id, servicos_id),
	CONSTRAINT clientes_clientes_atuacao_pkey PRIMARY KEY (id),
	CONSTRAINT clientes_clientes_at_clientes_id_eee8dafa_fk_clientes_ FOREIGN KEY (clientes_id) REFERENCES public.clientes_clientes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT clientes_clientes_at_servicos_id_69591ed3_fk_servicos_ FOREIGN KEY (servicos_id) REFERENCES public.servicos_servicos(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX clientes_clientes_atuacao_clientes_id_eee8dafa ON public.clientes_clientes_atuacao USING btree (clientes_id);
CREATE INDEX clientes_clientes_atuacao_servicos_id_69591ed3 ON public.clientes_clientes_atuacao USING btree (servicos_id);


-- public.compras_compras definition

CREATE TABLE public.compras_compras (
	id serial4 NOT NULL,
	valor numeric(11, 2) NOT NULL,
	"data" timestamptz NOT NULL,
	forma_pagamento varchar(30) NOT NULL,
	autorizado bool NOT NULL,
	codigo_compra_cielo varchar(80) NULL,
	transacao_cielo varchar(80) NULL,
	msg_cielo varchar(150) NULL,
	id_cliente_id int4 NOT NULL,
	id_usuario_id int4 NOT NULL,
	CONSTRAINT compras_compras_pkey PRIMARY KEY (id),
	CONSTRAINT compras_compras_id_cliente_id_400d1005_fk_clientes_clientes_id FOREIGN KEY (id_cliente_id) REFERENCES public.clientes_clientes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT compras_compras_id_usuario_id_cb5a7879_fk_usuarios_user_id FOREIGN KEY (id_usuario_id) REFERENCES public.usuarios_user(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX compras_compras_codcliente_id_73bda84a ON public.compras_compras USING btree (id_cliente_id);
CREATE INDEX compras_compras_codusuario_id_9fe2cf24 ON public.compras_compras USING btree (id_usuario_id);


-- public.precos_precos definition

CREATE TABLE public.precos_precos (
	id serial4 NOT NULL,
	valor numeric(11, 2) NOT NULL,
	tipo_servico varchar(1) NULL,
	codservico_id int4 NULL,
	CONSTRAINT precos_precos_codservico_id_key UNIQUE (codservico_id),
	CONSTRAINT precos_precos_pkey PRIMARY KEY (id),
	CONSTRAINT precos_precos_codservico_id_e3c8ba97_fk_servicos_servicos_id FOREIGN KEY (codservico_id) REFERENCES public.servicos_servicos(id) DEFERRABLE INITIALLY DEFERRED
);


-- public.registros_arquivoregistro definition

CREATE TABLE public.registros_arquivoregistro (
	id serial4 NOT NULL,
	file varchar(255) NULL,
	"name" varchar(250) NOT NULL,
	"size" int4 NOT NULL,
	signature varchar(255) NULL,
	"version" numeric(9, 2) NOT NULL,
	id_usuario_id int4 NOT NULL,
	create_at timestamptz NOT NULL,
	paid bool NOT NULL,
	update_at timestamptz NOT NULL,
	resume text NULL,
	value numeric(9, 2) NOT NULL,
	b64 text NULL,
	code varchar(100) NULL,
	content_type varchar(50) NULL,
	codregistro int4 NULL,
	CONSTRAINT registros_arquivoregistro_codregistro_check CHECK ((codregistro >= 0)),
	CONSTRAINT registros_arquivoregistro_pkey PRIMARY KEY (id),
	CONSTRAINT registros_arquivoregistro_size_3d38ebb2_check CHECK ((size >= 0)),
	CONSTRAINT registros_arquivoreg_id_usuario_id_10199dc8_fk_usuarios_ FOREIGN KEY (id_usuario_id) REFERENCES public.usuarios_user(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX registros_arquivoregistro_id_usuario_id_10199dc8 ON public.registros_arquivoregistro USING btree (id_usuario_id);


-- public.registros_coautores_coautores definition

CREATE TABLE public.registros_coautores_coautores (
	id serial4 NOT NULL,
	nome varchar(255) NOT NULL,
	percentual_obra numeric(5, 2) NOT NULL,
	documento varchar(50) NOT NULL,
	arquivo_id int4 NOT NULL,
	coautor_id int4 NULL,
	CONSTRAINT registros_coautores_coautores_pkey PRIMARY KEY (id),
	CONSTRAINT registros_coautores__arquivo_id_d05b90a6_fk_registros FOREIGN KEY (arquivo_id) REFERENCES public.registros_arquivoregistro(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX registros_coautores_coautores_arquivo_id_d05b90a6 ON public.registros_coautores_coautores USING btree (arquivo_id);


-- public.registros_registros definition

CREATE TABLE public.registros_registros (
	id serial4 NOT NULL,
	codregistro int4 NULL,
	valor numeric(9, 2) NOT NULL,
	"data" timestamptz NOT NULL,
	arquivo_id int4 NULL,
	descricao varchar(255) NULL,
	desconto numeric(9, 2) NOT NULL,
	id_cliente_id int4 NULL,
	codservico_id int4 NOT NULL,
	id_usuario_id int4 NOT NULL,
	codqrcode int4 NULL,
	codigo_promocional varchar(200) NULL,
	manter_arquivo bool NOT NULL,
	excluido bool NOT NULL,
	CONSTRAINT registros_registros_arquivo_id_c4e4ac90_uniq UNIQUE (arquivo_id),
	CONSTRAINT registros_registros_codqrcode_check CHECK ((codqrcode >= 0)),
	CONSTRAINT registros_registros_codregistro_check CHECK ((codregistro >= 0)),
	CONSTRAINT registros_registros_pkey PRIMARY KEY (id),
	CONSTRAINT registros_registros_arquivo_id_c4e4ac90_fk_registros FOREIGN KEY (arquivo_id) REFERENCES public.registros_arquivoregistro(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT registros_registros_codservico_id_6fbae033_fk_servicos_ FOREIGN KEY (codservico_id) REFERENCES public.servicos_servicos(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT registros_registros_id_cliente_id_6f9ce4f9_fk_clientes_ FOREIGN KEY (id_cliente_id) REFERENCES public.clientes_clientes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT registros_registros_id_usuario_id_e881ddd6_fk_usuarios_user_id FOREIGN KEY (id_usuario_id) REFERENCES public.usuarios_user(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX registros_registros_codcliente_id_795d50a8 ON public.registros_registros USING btree (id_cliente_id);
CREATE INDEX registros_registros_codservico_id_6fbae033 ON public.registros_registros USING btree (codservico_id);
CREATE INDEX registros_registros_codusuario_id_d77c0c50 ON public.registros_registros USING btree (id_usuario_id);


-- public.servicos_extensoes_servicos_extensoes definition

CREATE TABLE public.servicos_extensoes_servicos_extensoes (
	id serial4 NOT NULL,
	codextensao_id int4 NOT NULL,
	codservico_id int4 NOT NULL,
	CONSTRAINT servicos_extensoes_servicos_extensoes_pkey PRIMARY KEY (id),
	CONSTRAINT servicos_extensoes_s_codextensao_id_abb02c62_fk_extensoes FOREIGN KEY (codextensao_id) REFERENCES public.extensoes_extensoes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT servicos_extensoes_s_codservico_id_15194aa7_fk_servicos_ FOREIGN KEY (codservico_id) REFERENCES public.servicos_servicos(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX servicos_extensoes_servicos_extensoes_codextensao_id_abb02c62 ON public.servicos_extensoes_servicos_extensoes USING btree (codextensao_id);
CREATE INDEX servicos_extensoes_servicos_extensoes_codservico_id_15194aa7 ON public.servicos_extensoes_servicos_extensoes USING btree (codservico_id);


-- public.servicos_servicos_extensoes definition

CREATE TABLE public.servicos_servicos_extensoes (
	id serial4 NOT NULL,
	servicos_id int4 NOT NULL,
	extensoes_id int4 NOT NULL,
	CONSTRAINT servicos_servicos_extens_servicos_id_extensoes_id_60f1a7fa_uniq UNIQUE (servicos_id, extensoes_id),
	CONSTRAINT servicos_servicos_extensoes_pkey PRIMARY KEY (id),
	CONSTRAINT servicos_servicos_ex_extensoes_id_737c4577_fk_extensoes FOREIGN KEY (extensoes_id) REFERENCES public.extensoes_extensoes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT servicos_servicos_ex_servicos_id_d0a4c666_fk_servicos_ FOREIGN KEY (servicos_id) REFERENCES public.servicos_servicos(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX servicos_servicos_extensoes_extensoes_id_737c4577 ON public.servicos_servicos_extensoes USING btree (extensoes_id);
CREATE INDEX servicos_servicos_extensoes_servicos_id_d0a4c666 ON public.servicos_servicos_extensoes USING btree (servicos_id);


-- public.clientes_atuacao_clientes_atuacao definition

CREATE TABLE public.clientes_atuacao_clientes_atuacao (
	id serial4 NOT NULL,
	id_cliente_id int4 NOT NULL,
	codservico_id int4 NOT NULL,
	id_usuario_id int4 NOT NULL,
	CONSTRAINT clientes_atuacao_clientes_atuacao_pkey PRIMARY KEY (id),
	CONSTRAINT clientes_atuacao_cli_codservico_id_c7e6f5b0_fk_servicos_ FOREIGN KEY (codservico_id) REFERENCES public.servicos_servicos(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT clientes_atuacao_cli_id_cliente_id_1778bd7e_fk_clientes_ FOREIGN KEY (id_cliente_id) REFERENCES public.clientes_clientes(id) DEFERRABLE INITIALLY DEFERRED,
	CONSTRAINT clientes_atuacao_cli_id_usuario_id_061284c1_fk_usuarios_ FOREIGN KEY (id_usuario_id) REFERENCES public.usuarios_user(id) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX clientes_atuacao_clientes_atuacao_codcliente_id_e5eb8a59 ON public.clientes_atuacao_clientes_atuacao USING btree (id_cliente_id);
CREATE INDEX clientes_atuacao_clientes_atuacao_codservico_id_c7e6f5b0 ON public.clientes_atuacao_clientes_atuacao USING btree (codservico_id);
CREATE INDEX clientes_atuacao_clientes_atuacao_codusuario_id_9f8e05c3 ON public.clientes_atuacao_clientes_atuacao USING btree (id_usuario_id);
