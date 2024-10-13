package db

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// Open SQLite database
func openDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./ble.db")
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    return db, nil
}

// Find UUID with is_active set to 0
func findInactiveUUID(db *sql.DB) (string, error) {
    var uuid string
    query := `SELECT uuid FROM devices WHERE is_active = 0 LIMIT 1`
    err := db.QueryRow(query).Scan(&uuid)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("There are not enough devices available.")
        }
        return "", fmt.Errorf("%v", err)
    }
    return uuid, nil
}

// Update is_active value of the UUID to 1
func updateUUIDStatusToActive(db *sql.DB, uuid string) error {
    query := `UPDATE devices SET is_active = 1 WHERE uuid = ?`
    _, err := db.Exec(query, uuid)
    if err != nil {
        return fmt.Errorf("failed to update UUID status: %v", err)
    }
    return nil
}

// GetAndActivateUUID: Function to find and activate a UUID
func GetAndActivateUUID() (string, error) {
    db, err := openDB()
    if err != nil {
        return "", fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    // Find UUID with is_active set to 0
    uuid, err := findInactiveUUID(db)
    if err != nil {
        return "", fmt.Errorf("%v", err)
    }

    // Update is_active value of the UUID to 1
    err = updateUUIDStatusToActive(db, uuid)
    if err != nil {
        return "", fmt.Errorf("Failed to update UUID status to active: %v", err)
    }

    return uuid, nil
}
