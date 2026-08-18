package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string

	SupabaseS3Endpoint        string
	SupabaseS3Region          string
	SupabaseS3AccessKeyID     string
	SupabaseS3SecretAccessKey string
	SupabaseStorageBucket     string
}

func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://kayam:kayam@localhost:5432/kayam?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),

		SupabaseS3Endpoint:        getEnv("SUPABASE_S3_ENDPOINT", ""),
		SupabaseS3Region:          getEnv("SUPABASE_S3_REGION", ""),
		SupabaseS3AccessKeyID:     getEnv("SUPABASE_S3_ACCESS_KEY_ID", ""),
		SupabaseS3SecretAccessKey: getEnv("SUPABASE_S3_SECRET_ACCESS_KEY", ""),
		SupabaseStorageBucket:     getEnv("SUPABASE_STORAGE_BUCKET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
