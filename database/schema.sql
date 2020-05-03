--
-- PostgreSQL database dump
--

-- Dumped from database version 10.12
-- Dumped by pg_dump version 10.12

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
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: db4s_download_info; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_download_info (
    download_id integer NOT NULL,
    friendly_name text
);


--
-- Name: db4s_download_info_download_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_download_info_download_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_download_info_download_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_download_info_download_id_seq OWNED BY public.db4s_download_info.download_id;


--
-- Name: db4s_downloads_daily; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_downloads_daily (
    daily_id integer NOT NULL,
    stats_date timestamp without time zone,
    db4s_download integer,
    num_downloads integer
);


--
-- Name: db4s_downloads_daily_daily_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_downloads_daily_daily_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_downloads_daily_daily_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_downloads_daily_daily_id_seq OWNED BY public.db4s_downloads_daily.daily_id;


--
-- Name: db4s_downloads_monthly; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_downloads_monthly (
    monthly_id integer NOT NULL,
    stats_date timestamp without time zone,
    db4s_download integer,
    num_downloads integer
);


--
-- Name: db4s_downloads_monthly_monthly_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_downloads_monthly_monthly_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_downloads_monthly_monthly_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_downloads_monthly_monthly_id_seq OWNED BY public.db4s_downloads_monthly.monthly_id;


--
-- Name: db4s_downloads_weekly; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_downloads_weekly (
    weekly_id integer NOT NULL,
    stats_date timestamp without time zone,
    db4s_download integer,
    num_downloads integer
);


--
-- Name: db4s_downloads_weekly_weekly_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_downloads_weekly_weekly_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_downloads_weekly_weekly_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_downloads_weekly_weekly_id_seq OWNED BY public.db4s_downloads_weekly.weekly_id;


--
-- Name: db4s_release_info; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_release_info (
    release_id integer NOT NULL,
    version_number text,
    friendly_name text
);


--
-- Name: db4s_release_info_release_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_release_info_release_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_release_info_release_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_release_info_release_id_seq OWNED BY public.db4s_release_info.release_id;


--
-- Name: db4s_users_daily; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_users_daily (
    daily_id integer NOT NULL,
    stats_date timestamp without time zone,
    db4s_release integer,
    unique_ips integer
);


--
-- Name: db4s_users_daily_daily_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_users_daily_daily_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_users_daily_daily_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_users_daily_daily_id_seq OWNED BY public.db4s_users_daily.daily_id;


--
-- Name: db4s_users_monthly; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_users_monthly (
    monthly_id integer NOT NULL,
    stats_date timestamp without time zone,
    db4s_release integer,
    unique_ips integer
);


--
-- Name: db4s_users_monthly_monthly_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_users_monthly_monthly_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_users_monthly_monthly_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_users_monthly_monthly_id_seq OWNED BY public.db4s_users_monthly.monthly_id;


--
-- Name: db4s_users_weekly; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.db4s_users_weekly (
    weekly_id integer NOT NULL,
    stats_date timestamp without time zone,
    db4s_release integer,
    unique_ips integer
);


--
-- Name: db4s_users_weekly_weekly_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.db4s_users_weekly_weekly_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: db4s_users_weekly_weekly_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.db4s_users_weekly_weekly_id_seq OWNED BY public.db4s_users_weekly.weekly_id;


--
-- Name: download_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.download_log (
    remote_addr text,
    remote_user text,
    request_time timestamp with time zone,
    request_type text,
    request text,
    protocol text,
    status integer,
    body_bytes_sent bigint,
    http_referer text,
    http_user_agent text,
    download_id bigint NOT NULL,
    client_ipv4 text,
    client_ipv6 text,
    client_ip_strange text,
    client_port integer
);


--
-- Name: download_log_download_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.download_log_download_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: download_log_download_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.download_log_download_id_seq OWNED BY public.download_log.download_id;


--
-- Name: github_download_counts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_download_counts (
    count_id integer NOT NULL,
    asset integer,
    download_count integer NOT NULL,
    count_timestamp integer NOT NULL
);


--
-- Name: github_download_counts_count_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.github_download_counts_count_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: github_download_counts_count_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.github_download_counts_count_id_seq OWNED BY public.github_download_counts.count_id;


--
-- Name: github_download_timestamps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_download_timestamps (
    timestamp_id integer NOT NULL,
    count_timestamp timestamp with time zone NOT NULL
);


--
-- Name: github_download_timestamps_timestamp_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.github_download_timestamps_timestamp_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: github_download_timestamps_timestamp_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.github_download_timestamps_timestamp_id_seq OWNED BY public.github_download_timestamps.timestamp_id;


--
-- Name: github_release_assets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_release_assets (
    asset_id integer NOT NULL,
    asset_name text NOT NULL
);


--
-- Name: github_release_assets_asset_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.github_release_assets_asset_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: github_release_assets_asset_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.github_release_assets_asset_id_seq OWNED BY public.github_release_assets.asset_id;


--
-- Name: db4s_download_info download_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_download_info ALTER COLUMN download_id SET DEFAULT nextval('public.db4s_download_info_download_id_seq'::regclass);


--
-- Name: db4s_downloads_daily daily_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_daily ALTER COLUMN daily_id SET DEFAULT nextval('public.db4s_downloads_daily_daily_id_seq'::regclass);


--
-- Name: db4s_downloads_monthly monthly_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_monthly ALTER COLUMN monthly_id SET DEFAULT nextval('public.db4s_downloads_monthly_monthly_id_seq'::regclass);


--
-- Name: db4s_downloads_weekly weekly_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_weekly ALTER COLUMN weekly_id SET DEFAULT nextval('public.db4s_downloads_weekly_weekly_id_seq'::regclass);


--
-- Name: db4s_release_info release_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_release_info ALTER COLUMN release_id SET DEFAULT nextval('public.db4s_release_info_release_id_seq'::regclass);


--
-- Name: db4s_users_daily daily_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_users_daily ALTER COLUMN daily_id SET DEFAULT nextval('public.db4s_users_daily_daily_id_seq'::regclass);


--
-- Name: db4s_users_monthly monthly_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_users_monthly ALTER COLUMN monthly_id SET DEFAULT nextval('public.db4s_users_monthly_monthly_id_seq'::regclass);


--
-- Name: db4s_users_weekly weekly_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_users_weekly ALTER COLUMN weekly_id SET DEFAULT nextval('public.db4s_users_weekly_weekly_id_seq'::regclass);


--
-- Name: download_log download_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.download_log ALTER COLUMN download_id SET DEFAULT nextval('public.download_log_download_id_seq'::regclass);


--
-- Name: github_download_counts count_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_download_counts ALTER COLUMN count_id SET DEFAULT nextval('public.github_download_counts_count_id_seq'::regclass);


--
-- Name: github_download_timestamps timestamp_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_download_timestamps ALTER COLUMN timestamp_id SET DEFAULT nextval('public.github_download_timestamps_timestamp_id_seq'::regclass);


--
-- Name: github_release_assets asset_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_release_assets ALTER COLUMN asset_id SET DEFAULT nextval('public.github_release_assets_asset_id_seq'::regclass);


--
-- Name: db4s_download_info db4s_download_info_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_download_info
    ADD CONSTRAINT db4s_download_info_pk PRIMARY KEY (download_id);


--
-- Name: db4s_downloads_daily db4s_downloads_daily_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_daily
    ADD CONSTRAINT db4s_downloads_daily_pk PRIMARY KEY (daily_id);


--
-- Name: db4s_downloads_monthly db4s_downloads_monthly_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_monthly
    ADD CONSTRAINT db4s_downloads_monthly_pk PRIMARY KEY (monthly_id);


--
-- Name: db4s_downloads_weekly db4s_downloads_weekly_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_weekly
    ADD CONSTRAINT db4s_downloads_weekly_pk PRIMARY KEY (weekly_id);


--
-- Name: download_log download_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.download_log
    ADD CONSTRAINT download_log_pkey PRIMARY KEY (download_id);


--
-- Name: github_download_timestamps github_download_timestamps_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_download_timestamps
    ADD CONSTRAINT github_download_timestamps_pk PRIMARY KEY (timestamp_id);


--
-- Name: github_release_assets github_release_assets_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_release_assets
    ADD CONSTRAINT github_release_assets_pk PRIMARY KEY (asset_id);


--
-- Name: db4s_downloads_daily_stats_date_db4s_download_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_downloads_daily_stats_date_db4s_download_uindex ON public.db4s_downloads_daily USING btree (stats_date, db4s_download);


--
-- Name: db4s_downloads_monthly_stats_date_db4s_download_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_downloads_monthly_stats_date_db4s_download_uindex ON public.db4s_downloads_monthly USING btree (stats_date, db4s_download);


--
-- Name: db4s_downloads_weekly_stats_date_db4s_download_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_downloads_weekly_stats_date_db4s_download_uindex ON public.db4s_downloads_weekly USING btree (stats_date, db4s_download);


--
-- Name: db4s_release_info_release_id_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_release_info_release_id_uindex ON public.db4s_release_info USING btree (release_id);


--
-- Name: db4s_release_info_version_number_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_release_info_version_number_uindex ON public.db4s_release_info USING btree (version_number);


--
-- Name: db4s_users_daily_stats_date_db4s_release_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_users_daily_stats_date_db4s_release_uindex ON public.db4s_users_daily USING btree (stats_date, db4s_release);


--
-- Name: db4s_users_monthly_stats_date_db4s_release_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_users_monthly_stats_date_db4s_release_uindex ON public.db4s_users_monthly USING btree (stats_date, db4s_release);


--
-- Name: db4s_users_weekly_stats_date_db4s_release_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX db4s_users_weekly_stats_date_db4s_release_uindex ON public.db4s_users_weekly USING btree (stats_date, db4s_release);


--
-- Name: download_log_request_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX download_log_request_index ON public.download_log USING btree (request);


--
-- Name: download_log_request_time_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX download_log_request_time_index ON public.download_log USING btree (request_time);


--
-- Name: github_download_counts_count_id_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX github_download_counts_count_id_uindex ON public.github_download_counts USING btree (count_id);


--
-- Name: github_download_timestamps_download_timestamp_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX github_download_timestamps_download_timestamp_uindex ON public.github_download_timestamps USING btree (count_timestamp);


--
-- Name: github_download_timestamps_timestamp_id_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX github_download_timestamps_timestamp_id_uindex ON public.github_download_timestamps USING btree (timestamp_id);


--
-- Name: github_release_assets_asset_id_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX github_release_assets_asset_id_uindex ON public.github_release_assets USING btree (asset_id);


--
-- Name: github_release_assets_asset_name_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX github_release_assets_asset_name_uindex ON public.github_release_assets USING btree (asset_name);


--
-- Name: db4s_downloads_daily db4s_downloads_daily_db4s_download_info_download_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_daily
    ADD CONSTRAINT db4s_downloads_daily_db4s_download_info_download_id_fk FOREIGN KEY (db4s_download) REFERENCES public.db4s_download_info(download_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: db4s_downloads_monthly db4s_downloads_monthly_db4s_download_info_download_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_monthly
    ADD CONSTRAINT db4s_downloads_monthly_db4s_download_info_download_id_fk FOREIGN KEY (db4s_download) REFERENCES public.db4s_download_info(download_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: db4s_downloads_weekly db4s_downloads_weekly_db4s_download_info_download_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_downloads_weekly
    ADD CONSTRAINT db4s_downloads_weekly_db4s_download_info_download_id_fk FOREIGN KEY (db4s_download) REFERENCES public.db4s_download_info(download_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: db4s_users_daily db4s_users_daily_db4s_release_info_release_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_users_daily
    ADD CONSTRAINT db4s_users_daily_db4s_release_info_release_id_fk FOREIGN KEY (db4s_release) REFERENCES public.db4s_release_info(release_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: db4s_users_monthly db4s_users_monthly_db4s_release_info_release_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_users_monthly
    ADD CONSTRAINT db4s_users_monthly_db4s_release_info_release_id_fk FOREIGN KEY (db4s_release) REFERENCES public.db4s_release_info(release_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: db4s_users_weekly db4s_users_weekly_db4s_release_info_release_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.db4s_users_weekly
    ADD CONSTRAINT db4s_users_weekly_db4s_release_info_release_id_fk FOREIGN KEY (db4s_release) REFERENCES public.db4s_release_info(release_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: github_download_counts github_download_counts_github_download_timestamps_timestamp_id_; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_download_counts
    ADD CONSTRAINT github_download_counts_github_download_timestamps_timestamp_id_ FOREIGN KEY (count_timestamp) REFERENCES public.github_download_timestamps(timestamp_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: github_download_counts github_download_counts_github_release_assets_asset_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_download_counts
    ADD CONSTRAINT github_download_counts_github_release_assets_asset_id_fk FOREIGN KEY (asset) REFERENCES public.github_release_assets(asset_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

