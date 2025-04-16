package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	Username             string     `json:"username"`
	PasswordPlain        string     `json:"password_hash"`
	ExternalID           *int       `json:"external_id"`
	OrgUnitID            *int       `json:"orgunit_id"`
	GroupsIDStr          string     `json:"groups_id"`
	Phone                *string    `json:"phone"`
	IsLockedStr          string     `json:"is_locked"`
	IsAdminStr           string     `json:"is_admin"`
	CreateTime           time.Time  `json:"create_time"`
	Creator              int        `json:"creator"`
	LastModificationTime *time.Time `json:"last_modification_time"`
	LastModificationUser *int       `json:"last_modification_user"`
	IsActiveStr          string     `json:"is_active"`
	InactivationTime     *time.Time `json:"inactivation_time"`
	InactivationReason   *string    `json:"inactivation_reason"`
	Comment              *string    `json:"comment"`
	BadPasswordCount     int        `json:"bad_password_count"`
}

func parsePGArray(pgArray string) ([]int, error) {
	cleaned := strings.Trim(pgArray, "{}")
	if cleaned == "" {
		return []int{}, nil
	}
	parts := strings.Split(cleaned, ",")
	result := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		result[i] = n
	}
	return result, nil
}

func parseBool(str string) bool {
	return str == "t" || str == "true"
}

func main() {
	// Paramètres ligne de commande
	dbName := flag.String("db", "", "Nom de la base de données")
	dbUser := flag.String("user", "", "Nom d'utilisateur")
	dbPass := flag.String("pass", "", "Mot de passe")
	dbHost := flag.String("host", "localhost", "Hôte PostgreSQL")
	dbPort := flag.String("port", "5432", "Port PostgreSQL")
	jsonInputArg := flag.String("json", "", "Chaîne JSON à importer (alternative à stdin)")
	flag.Parse()

	if *dbName == "" || *dbUser == "" || *dbPass == "" {
		log.Fatal("❌ Vous devez spécifier -db, -user et -pass")
	}

	var jsonInput []byte
	var err error

	if *jsonInputArg != "" {
		jsonInput = []byte(*jsonInputArg)
	} else {
		jsonInput, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Erreur lecture JSON depuis stdin:", err)
		}
	}

	/*
	Requête utilisée pour générer le JSON sur sqlserver: :
	
		SELECT 
			(E.Prenom + ' ' + E.Nom) AS name,
			E.email AS email,
			EM.Login AS username, 
			LOWER(CONVERT(VARCHAR(64), HASHBYTES('SHA2_256', EM.Alias), 2)) AS password_hash, 
			E.IdEmploye AS external_id,
			NULL AS orgunit_id,
			'{3}' AS groups_id,
			NULL AS phone,
			'f' AS is_locked,
			'f' AS is_admin,
			FORMAT(GETUTCDATE(), 'yyyy-MM-ddTHH:mm:ss.fff') + 'Z' AS create_time,
			10958 AS creator,
			NULL AS last_modification_time,
			NULL AS last_modification_user,
			't' AS is_active,
			NULL AS inactivation_time,
			NULL AS inactivation_reason,
			NULL AS comment,
			0 AS bad_password_count
		FROM EmployeMdP EM
		INNER JOIN EmployeMdPApplicationAcces MAA ON MAA.IdEmploye = EM.IdEmploye
		INNER JOIN Employe E ON E.IdEmploye = EM.IdEmploye AND E.IsActive = 1
		WHERE MAA.CodeApplication = 'arbre'
		ORDER BY EM.IdEmploye
		FOR JSON PATH, INCLUDE_NULL_VALUES;
	*/

	var users []User
	if err := json.Unmarshal(jsonInput, &users); err != nil {
		log.Fatalf("Erreur parsing JSON: %v", err)
	}

	// Connexion PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPass, *dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur de connexion à PostgreSQL:", err)
	}
	defer db.Close()

	// Préparation requête INSERT
	stmt := `
        INSERT INTO public.go_user (
            name, email, username, password_hash,
            external_id, orgunit_id, groups_id, phone,
            is_locked, is_admin, create_time, creator,
            last_modification_time, last_modification_user,
            is_active, inactivation_time, inactivation_reason,
            comment, bad_password_count
        )
        VALUES (
            $1, $2, $3, $4,
            $5, $6, $7, $8,
            $9, $10, $11, $12,
            $13, $14,
            $15, $16, $17,
            $18, $19
        )
        ON CONFLICT (username) DO NOTHING;
    `

	for _, u := range users {
		groupsID, err := parsePGArray(u.GroupsIDStr)
		if err != nil {
			log.Println("Erreur parsing groups_id pour", u.Username, ":", err)
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordPlain), 4)
		if err != nil {
			log.Println("Erreur Bcrypt pour", u.Username, ":", err)
			continue
		}

		_, err = db.Exec(stmt,
			u.Name,
			u.Email,
			u.Username,
			string(hashedPassword),
			u.ExternalID,
			u.OrgUnitID,
			pq.Array(groupsID),
			u.Phone,
			parseBool(u.IsLockedStr),
			parseBool(u.IsAdminStr),
			u.CreateTime,
			u.Creator,
			u.LastModificationTime,
			u.LastModificationUser,
			parseBool(u.IsActiveStr),
			u.InactivationTime,
			u.InactivationReason,
			u.Comment,
			u.BadPasswordCount,
		)

		if err != nil {
			log.Println("❌ Erreur insertion:", u.Username, "->", err)
		} else {
			fmt.Println("✅ Utilisateur inséré:", u.Username)
		}
	}

	fmt.Println("🎉 Import terminé avec hachage cost=4.")
}
