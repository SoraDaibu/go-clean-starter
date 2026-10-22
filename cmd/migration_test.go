package cmd

import (
	"net/url"
	"testing"
)

func TestMigrationDatabaseURL(t *testing.T) {
	for _, sslMode := range []string{"", "require"} {
		t.Run(sslMode, func(t *testing.T) {
			t.Setenv("PGSSLMODE", sslMode)
			t.Setenv("DB_HOST", "::1")
			t.Setenv("DB_PORT", "5432")
			t.Setenv("DB_NAME", "starter")
			t.Setenv("DB_USER", "test@user")
			t.Setenv("DB_PASSWORD", "p@ss:/?# word ")
			u, err := url.Parse(getMigrationDatabaseURL())
			if err != nil {
				t.Fatal(err)
			}
			password, _ := u.User.Password()
			if u.Scheme != "postgres" || u.Host != "[::1]:5432" || u.Path != "/starter" || u.User.Username() != "test@user" || password != "p@ss:/?# word " {
				t.Fatalf("incorrect PostgreSQL connection URL components")
			}
			wantSSL := sslMode
			if wantSSL == "" {
				wantSSL = "disable"
			}
			if u.Query().Get("sslmode") != wantSSL {
				t.Fatalf("sslmode = %q; want %q", u.Query().Get("sslmode"), wantSSL)
			}
		})
	}
}
