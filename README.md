# sanarbo
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_sanarbo&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_sanarbo)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_sanarbo&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_sanarbo)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_sanarbo&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_sanarbo)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=lao-tseu-is-alive_sanarbo&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=lao-tseu-is-alive_sanarbo)
[![cve-trivy-scan](https://github.com/lao-tseu-is-alive/sanarbo/actions/workflows/cve-trivy-scan.yml/badge.svg)](https://github.com/lao-tseu-is-alive/sanarbo/actions/workflows/cve-trivy-scan.yml)

**sanarbo** means literally "health of tree" in esperanto. It's a web app to manage the health care of trees in the City of Lausanne (Switzerland).

## Dependencies
[Echo: high performance, extensible, minimalist Go web framework](https://echo.labstack.com/)

[deepmap/oapi-codegen: OpenAPI Client and Server Code Generator](https://github.com/deepmap/oapi-codegen)

[pgx: PostgreSQL Driver and Toolkit](https://pkg.go.dev/github.com/jackc/pgx)

[Json Web Token for Go (RFC 7519)](https://github.com/cristalhq/jwt)

[OpenLayers](https://openlayers.org/)

[Common libraries package to other Micro Service in our goeland team](https://github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs)


## Project Layout and conventions
This project uses the Standard Go Project Layout : https://github.com/golang-standards/project-layout

## Configuration
Add in the root directory of the project an .env file based on .env_sample with your own values.

The following goeland tables are needed to use this application:
* in the public schema:
    * go_metadata_db_schema
    * go_group
    * go_orgunit
    * go_user
    * thi_arbre_chk
    * thi_arbre_entourage
    * thi_arbre_etat_sanitaire
    * thi_arbre_etat_sanitaire_remarque
    * thi_arbre_note
    * thi_arbre_rev_surface
    * thi_arbre_to_be_checked
    * thi_arbre_validation
    * thi_street
    * thi_street_building_address

* in the geodata_gestion_com schema:
    * spadom_secteurs
    * spadom_surfaces

The sanarbo role must have at least SELECT privileges on these tables