package main

import ()

type Server struct {
	user        string
	password    string
	hostname    string
	feduser     string
	fedpassword string
	fedhostname string
	database    string
	federateddb string
	servername  string
}

var shard0 = Server{"app", "JsyJ586n", "ps-prod-db02.plansource.com", "federated", "JsyJ586n", "ps-prod-federated-db01.plansource.com", "benefits_production", "benefits_shard0", "db02"}
var shard01 = Server{"app", "JsyJ586n", "ps-prod-db07.plansource.com", "federated", "JsyJ586n", "ps-prod-federated-db01.plansource.com", "benefits_production", "benefits_shard01", "db07"}
var shard02 = Server{"app", "JsyJ586n", "ps-prod-db09.plansource.com", "federated", "JsyJ586n", "ps-prod-federated-db01.plansource.com", "benefits_production", "benefits_shard02", "db09"}
var shard03 = Server{"app", "JsyJ586n", "ps-prod-db11.plansource.com", "federated", "JsyJ586n", "ps-prod-federated-db01.plansource.com", "benefits_production", "benefits_shard03", "db11"}
var archive = Server{"app", "JsyJ586n", "ps-backup-db05.plansource.com", "federated", "JsyJ586n", "ps-prod-federated-db01.plansource.com", "archive_production", "benefits_archive", "db05"}

func getConnection(request string) Server {
	if request == "shard0" {
		return shard0
	} else if request == "shard01" {
		return shard01
	} else if request == "shard02" {
		return shard02
	} else if request == "shard03" {
		return shard03
	} else if request == "archive" {
		return archive
	} else {
		return Server{"fail", "fail", "fail", "fail", "fail", "fail", "fail", "fail", "fail"}
	}
}
