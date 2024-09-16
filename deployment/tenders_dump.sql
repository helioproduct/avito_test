--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1 (Debian 16.1-1.pgdg120+1)
-- Dumped by pg_dump version 16.1 (Debian 16.1-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: organization_type; Type: TYPE; Schema: public; Owner: helio
--

CREATE TYPE public.organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);


ALTER TYPE public.organization_type OWNER TO helio;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bids; Type: TABLE; Schema: public; Owner: helio
--

CREATE TABLE public.bids (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    tender_id uuid NOT NULL,
    author_type text NOT NULL,
    author_id uuid NOT NULL,
    status text NOT NULL,
    decision text DEFAULT 'Pending'::text NOT NULL,
    version integer DEFAULT 1,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.bids OWNER TO helio;

--
-- Name: employee; Type: TABLE; Schema: public; Owner: helio
--

CREATE TABLE public.employee (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(50) NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.employee OWNER TO helio;

--
-- Name: organization; Type: TABLE; Schema: public; Owner: helio
--

CREATE TABLE public.organization (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    type public.organization_type,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.organization OWNER TO helio;

--
-- Name: organization_responsible; Type: TABLE; Schema: public; Owner: helio
--

CREATE TABLE public.organization_responsible (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    organization_id uuid,
    user_id uuid
);


ALTER TABLE public.organization_responsible OWNER TO helio;

--
-- Name: tender; Type: TABLE; Schema: public; Owner: helio
--

CREATE TABLE public.tender (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    servicetype text NOT NULL,
    organization_id uuid,
    creator_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    current_version integer DEFAULT 1,
    status text DEFAULT 'Created'::text NOT NULL,
    CONSTRAINT status_check CHECK ((status = ANY (ARRAY['Created'::text, 'Published'::text, 'Closed'::text]))),
    CONSTRAINT tender_servicetype_check CHECK ((servicetype = ANY (ARRAY['Construction'::text, 'Delivery'::text, 'Manufacture'::text])))
);


ALTER TABLE public.tender OWNER TO helio;

--
-- Name: tender_versions; Type: TABLE; Schema: public; Owner: helio
--

CREATE TABLE public.tender_versions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    tender_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    servicetype text NOT NULL,
    organization_id uuid,
    creator_id uuid NOT NULL,
    version integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT tender_versions_servicetype_check CHECK ((servicetype = ANY (ARRAY['Construction'::text, 'Delivery'::text, 'Manufacture'::text])))
);


CREATE TABLE public.decisions (
    -- decision_id SERIAL PRIMARY KEY,
    bid_id UUID NOT NULL,
    username TEXT NOT NULL,
    decision_value TEXT NOT NULL CHECK (decision_value IN ('Approved', 'Rejected')),
    UNIQUE (bid_id, username, decision_value) 
);


ALTER TABLE public.tender_versions OWNER TO helio;

--
-- Data for Name: bids; Type: TABLE DATA; Schema: public; Owner: helio
--

COPY public.bids (id, name, description, tender_id, author_type, author_id, status, decision, version, created_at, updated_at) FROM stdin;
13fb3e95-e14f-410a-8c7e-370d8d55803b	сделать сальто	просто пойти и ебнуть	d2b5d591-180b-43ac-9086-c498c23d8245	Organization	e23e620a-d28e-45ce-b7a4-a9983fe3b218	Created	Approved	1	2024-09-13 19:11:29.514002	2024-09-13 16:46:45.368676
d214b5b0-0b33-43e5-b37c-2f3d7666b69f	Самое лучшее предложение	описание	c86505fe-cbb0-4e26-9e88-ac4cb3a49ef8	Organization	b570ea0f-4d18-4e56-bbec-a2a2c0f65145	Created	Pending	1	2024-09-16 20:05:25.192918	2024-09-16 20:05:25.192919
8e661fd9-c4ff-41d0-9610-cd198aed8c9f	Самое лучшее предложение 3	описание 3	c86505fe-cbb0-4e26-9e88-ac4cb3a49ef8	User	60dec5e9-9d1e-4574-9e47-f26d22157ada	Published	Pending	1	2024-09-16 20:10:23.846269	2024-09-16 17:55:46.54797
0c2982fc-c198-4509-b089-12165575cc2b	предлжоение 4	ало да	c86505fe-cbb0-4e26-9e88-ac4cb3a49ef8	Organization	b570ea0f-4d18-4e56-bbec-a2a2c0f65145	Created	Pending	1	2024-09-16 21:00:04.993883	2024-09-16 21:00:04.993883
18f09d63-cf43-42f3-8c08-925b302c2316	Самое лучшее предложение 2	самое лучшее предложение на диком fsdfsd	c86505fe-cbb0-4e26-9e88-ac4cb3a49ef8	Organization	161fcb6d-fcca-4f63-9342-86562acec2df	Created	Pending	6	2024-09-16 20:07:11.371071	2024-09-16 19:41:55.916894
0f0f15de-3c92-4eff-9c5e-4ea4d9e2ed02	самое лучшее предложение на Ближнем Востоке	самое лучшее описание	d2b5d591-180b-43ac-9086-c498c23d8245	Organization	e23e620a-d28e-45ce-b7a4-a9983fe3b218	Created	Pending	6	2024-09-16 19:47:44.512622	2024-09-16 19:48:38.36254
e7942813-e2e7-460d-9e23-bf36167a677d	доствим за 50 минут	безопасно и дешево	d2b5d591-180b-43ac-9086-c498c23d8245	Organization	b570ea0f-4d18-4e56-bbec-a2a2c0f65145	Published	Approved	1	2024-09-16 20:34:52.789017	2024-09-16 20:37:49.938132
\.


--
-- Data for Name: employee; Type: TABLE DATA; Schema: public; Owner: helio
--

COPY public.employee (id, username, first_name, last_name, created_at, updated_at) FROM stdin;
60dec5e9-9d1e-4574-9e47-f26d22157ada	helio	Nikolay	Popov	2024-09-13 09:33:32.964307	2024-09-13 09:33:32.964307
b5cce1d2-d9b2-4eb6-91e7-ff707a61f3fe	happy	Arina	Noskova	2024-09-13 09:46:15.874524	2024-09-13 09:46:15.874524
8d7eafb0-c219-4ec2-bf3a-278bf107d145	johndoe	John	Doe	2024-09-13 23:01:43.746665	2024-09-13 23:01:43.746665
0d4d9be7-6d31-4752-b1ea-7d8fb5d4508a	janesmith	Jane	Smith	2024-09-13 23:01:43.763768	2024-09-13 23:01:43.763768
1244a699-6d6c-49c4-b097-b5406ef58a5c	alicebrown	Alice	Brown	2024-09-13 23:01:43.764763	2024-09-13 23:01:43.764763
cb55cd6c-b5e3-475b-a4b4-7f76c0bd2516	bobjones	Bob	Jones	2024-09-13 23:01:43.765804	2024-09-13 23:01:43.765804
8dd95bcd-11fe-4754-8105-e37a75896408	charliedavis	Charlie	Davis	2024-09-13 23:01:43.766512	2024-09-13 23:01:43.766512
4a7bf99b-ecd0-4166-beec-25e39b7adf06	eveadams	Eve	Adams	2024-09-13 23:01:43.767205	2024-09-13 23:01:43.767205
6da4d1ff-1611-4230-a238-d7d1ecb98e45	davidclark	David	Clark	2024-09-13 23:01:43.767703	2024-09-13 23:01:43.767703
375ee308-1633-4eb7-b03c-5287582c2cca	gracemiller	Grace	Miller	2024-09-13 23:01:43.76845	2024-09-13 23:01:43.76845
25f61718-0cc6-49f6-a4b2-d27629940e69	hankwilson	Hank	Wilson	2024-09-13 23:01:43.769459	2024-09-13 23:01:43.769459
be3ffbcd-add7-4595-9606-43c013b2a78a	ivytaylor	Ivy	Taylor	2024-09-13 23:01:43.770128	2024-09-13 23:01:43.770128
\.


--
-- Data for Name: organization; Type: TABLE DATA; Schema: public; Owner: helio
--

COPY public.organization (id, name, description, type, created_at, updated_at) FROM stdin;
aef8f743-6bbe-4bb5-8d0e-b41f72b4d9cb	Tech Corp	Technology Company	LLC	2024-09-13 23:01:43.770666	2024-09-13 23:01:43.770666
b570ea0f-4d18-4e56-bbec-a2a2c0f65145	Health Corp	Healthcare Services	IE	2024-09-13 23:01:43.772871	2024-09-13 23:01:43.772871
2fa63d8c-667b-47f3-9a42-d21091182519	Edu Ltd	Education Services	IE	2024-09-13 23:01:43.773578	2024-09-13 23:01:43.773578
fef44152-9b1b-436f-9ae7-c65575a68813	AgriCo	Agriculture Company	JSC	2024-09-13 23:01:43.774147	2024-09-13 23:01:43.774147
93e06fae-d7ee-4f91-acf0-f55db86ac2c3	FinTech	Financial Technology	LLC	2024-09-13 23:01:43.775075	2024-09-13 23:01:43.775075
e23e620a-d28e-45ce-b7a4-a9983fe3b218	North Pools	pools from north	IE	2024-09-13 18:56:37.501662	2024-09-13 18:56:37.501662
161fcb6d-fcca-4f63-9342-86562acec2df	safen	safe, safety and safen	LLC	2024-09-13 12:01:45.935644	2024-09-13 12:01:45.935644
\.


--
-- Data for Name: organization_responsible; Type: TABLE DATA; Schema: public; Owner: helio
--

COPY public.organization_responsible (id, organization_id, user_id) FROM stdin;
2b2ab657-7746-4ccf-9edb-e52a599bf344	161fcb6d-fcca-4f63-9342-86562acec2df	60dec5e9-9d1e-4574-9e47-f26d22157ada
cfee07b5-18c2-44ef-8db6-df3219d4ccd7	e23e620a-d28e-45ce-b7a4-a9983fe3b218	b5cce1d2-d9b2-4eb6-91e7-ff707a61f3fe
bbacc00b-dedc-4075-9605-ca5492a815e6	b570ea0f-4d18-4e56-bbec-a2a2c0f65145	8d7eafb0-c219-4ec2-bf3a-278bf107d145
676c069e-e46b-491b-b0e3-21c521360286	161fcb6d-fcca-4f63-9342-86562acec2df	4a7bf99b-ecd0-4166-beec-25e39b7adf06
\.


--
-- Data for Name: tender; Type: TABLE DATA; Schema: public; Owner: helio
--

COPY public.tender (id, name, description, servicetype, organization_id, creator_id, created_at, updated_at, current_version, status) FROM stdin;
c86505fe-cbb0-4e26-9e88-ac4cb3a49ef8	Тендер на создание сервсиса для управления тендарами	описание	Construction	161fcb6d-fcca-4f63-9342-86562acec2df	4a7bf99b-ecd0-4166-beec-25e39b7adf06	2024-09-15 21:42:16.477331	2024-09-16 19:59:57.771013	9	Published
4d1962e5-2c77-4cac-a22a-9a3b1e4b72ad	построить стену	бывает и такое	Construction	e23e620a-d28e-45ce-b7a4-a9983fe3b218	b5cce1d2-d9b2-4eb6-91e7-ff707a61f3fe	2024-09-14 18:08:02.916855	2024-09-16 20:06:41.349528	3	Created
27e4edb5-f6c7-4958-ac20-ea2b69a61a35	another tender for safen	my first tender	Delivery	161fcb6d-fcca-4f63-9342-86562acec2df	60dec5e9-9d1e-4574-9e47-f26d22157ada	2024-09-14 18:12:07.291538	2024-09-14 18:12:07.291538	1	Closed
d2b5d591-180b-43ac-9086-c498c23d8245	доставить груз из точки A в точку A	не придумали	Delivery	161fcb6d-fcca-4f63-9342-86562acec2df	60dec5e9-9d1e-4574-9e47-f26d22157ada	2024-09-13 15:09:33.44043	2024-09-16 20:39:49.035478	36	Closed
\.


--
-- Data for Name: tender_versions; Type: TABLE DATA; Schema: public; Owner: helio
--

COPY public.tender_versions (id, tender_id, name, description, servicetype, organization_id, creator_id, version, created_at, updated_at) FROM stdin;
\.


--
-- Name: bids bids_pkey; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.bids
    ADD CONSTRAINT bids_pkey PRIMARY KEY (id);


--
-- Name: employee employee_pkey; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT employee_pkey PRIMARY KEY (id);


--
-- Name: employee employee_username_key; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT employee_username_key UNIQUE (username);


--
-- Name: organization organization_pkey; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);


--
-- Name: organization_responsible organization_responsible_pkey; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.organization_responsible
    ADD CONSTRAINT organization_responsible_pkey PRIMARY KEY (id);


--
-- Name: tender tender_pkey; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.tender
    ADD CONSTRAINT tender_pkey PRIMARY KEY (id);


--
-- Name: tender_versions tender_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.tender_versions
    ADD CONSTRAINT tender_versions_pkey PRIMARY KEY (id);


--
-- Name: bids bids_tender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.bids
    ADD CONSTRAINT bids_tender_id_fkey FOREIGN KEY (tender_id) REFERENCES public.tender(id);


--
-- Name: organization_responsible organization_responsible_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.organization_responsible
    ADD CONSTRAINT organization_responsible_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- Name: organization_responsible organization_responsible_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.organization_responsible
    ADD CONSTRAINT organization_responsible_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.employee(id) ON DELETE CASCADE;


--
-- Name: tender tender_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.tender
    ADD CONSTRAINT tender_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.employee(id) ON DELETE SET NULL;


--
-- Name: tender tender_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.tender
    ADD CONSTRAINT tender_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- Name: tender_versions tender_versions_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.tender_versions
    ADD CONSTRAINT tender_versions_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- Name: tender_versions tender_versions_tender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: helio
--

ALTER TABLE ONLY public.tender_versions
    ADD CONSTRAINT tender_versions_tender_id_fkey FOREIGN KEY (tender_id) REFERENCES public.tender(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

