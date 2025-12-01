# Database Migration Guide

## Overview

This guide explains how to apply database migrations from your GORM models to PostgreSQL in the Go School Management System.

---

## How It Works

### GORM AutoMigrate

GORM provides an `AutoMigrate` function that automatically:

- Creates tables if they don't exist
- Adds missing columns
- Adds missing indexes
- **Does NOT** delete columns or change column types (safe for production)

### Migration File

**Location:** [`internal/database/migrate.go`](file:///mnt/sda2/repos/go-school-management/internal/database/migrate.go)

This file contains the `RunMigrations()` function that migrates all models in the correct order:

```go
func RunMigrations() error {
    log.Println("🔄 Running database migrations...")

    err := DB.AutoMigrate(
        // Core entities first
        &department.Department{},
        &teacher.Teacher{},
        &student.Student{},
        &course.Course{},

        // Then dependent entities
        &attendance.Attendance{},
        &homework.Homework{},
        &exam.Exam{},
        &grade.Grade{},

        // Finally junction tables
        &student_courses.StudentCourse{},
        &students_homework.StudentHomework{},
    )

    if err != nil {
        log.Printf("❌ Migration failed: %v", err)
        return err
    }

    log.Println("✅ Database migrations completed successfully!")
    return nil
}
```

---

## Automatic Migration on Startup

The migration runs automatically when you start the server:

**In [`cmd/server/main.go`](file:///mnt/sda2/repos/go-school-management/cmd/server/main.go):**

```go
func main() {
    cfg := config.LoadConfig()

    // 1. Create database if it doesn't exist
    database.CreateDatabaseIfNotExists(cfg)

    // 2. Connect to database
    database.ConnectDB(cfg)

    // 3. Run migrations (creates/updates tables)
    database.RunMigrations()

    // 4. Start server
    r := gin.Default()
    r.Run(":" + cfg.AppPort)
}
```

---

## Running Migrations

### Method 1: Automatic (Recommended)

Simply start your application:

```bash
go run cmd/server/main.go
```

**Output:**

```
✅ Connected to PostgreSQL database: school_db
🔄 Running database migrations...
✅ Database migrations completed successfully!
🚀 Server starting on port: 8080
```

### Method 2: Manual Migration Command (Optional)

You can create a separate migration command for manual control:

**Create:** `cmd/migrate/main.go`

```go
package main

import (
    "log"
    "school_management/internal/config"
    "school_management/internal/database"
)

func main() {
    cfg := config.LoadConfig()

    database.CreateDatabaseIfNotExists(cfg)
    database.ConnectDB(cfg)

    if err := database.RunMigrations(); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    log.Println("Migration completed!")
}
```

**Run:**

```bash
go run cmd/migrate/main.go
```

---

## What Gets Created

When you run migrations, GORM will create these tables:

### Tables Created

| Table Name          | Description                |
| ------------------- | -------------------------- |
| `departments`       | Academic departments       |
| `teachers`          | Teacher information        |
| `students`          | Student information        |
| `courses`           | Course catalog             |
| `attendances`       | Attendance records         |
| `homework`          | Homework assignments       |
| `exams`             | Exam definitions           |
| `grades`            | Student grades             |
| `student_courses`   | Student-course enrollments |
| `students_homework` | Homework submissions       |

### Standard Fields (from gorm.Model)

Every table includes:

- `id` - Primary key (auto-increment)
- `created_at` - Timestamp when record was created
- `updated_at` - Timestamp when record was last updated
- `deleted_at` - Soft delete timestamp (NULL if not deleted)

---

## Verifying Migrations

### Check Tables in PostgreSQL

```bash
# Connect to PostgreSQL
psql -U postgres -d school_db

# List all tables
\dt

# Describe a specific table
\d students

# Exit
\q
```

### Expected Output

```
                List of relations
 Schema |        Name         | Type  |  Owner
--------+---------------------+-------+----------
 public | attendances         | table | postgres
 public | courses             | table | postgres
 public | departments         | table | postgres
 public | exams               | table | postgres
 public | grades              | table | postgres
 public | homework            | table | postgres
 public | student_courses     | table | postgres
 public | students            | table | postgres
 public | students_homework   | table | postgres
 public | teachers            | table | postgres
```

---

## Migration Order Matters

The order in `AutoMigrate()` is important because of foreign key constraints:

```
1. departments (no dependencies)
2. teachers (depends on departments)
3. students (no dependencies)
4. courses (depends on departments, teachers)
5. attendance (depends on students, courses)
6. homework (depends on courses)
7. exams (depends on courses)
8. grades (depends on students, exams)
9. student_courses (depends on students, courses)
10. students_homework (depends on students, homework)
```

---

## Common Migration Scenarios

### Adding a New Field to a Model

1. Add the field to your model:

   ```go
   type Student struct {
       gorm.Model
       FirstName string `gorm:"not null"`
       // ... existing fields

       // New field
       MiddleName string `gorm:"size:50"` // ← Add this
   }
   ```

2. Restart the application
3. GORM will automatically add the new column

### Creating a New Table

1. Create the model file
2. Add it to `migrate.go`:
   ```go
   err := DB.AutoMigrate(
       // ... existing models
       &newmodule.NewModel{}, // ← Add this
   )
   ```
3. Restart the application

### Dropping a Column (Manual)

⚠️ **GORM AutoMigrate does NOT drop columns automatically** (safety feature)

To drop a column manually:

```go
// In migrate.go, add after AutoMigrate:
DB.Migrator().DropColumn(&student.Student{}, "MiddleName")
```

Or use SQL:

```sql
ALTER TABLE students DROP COLUMN middle_name;
```

---

## Migration Best Practices

### 1. **Always Backup Before Migration**

```bash
# Backup database
pg_dump -U postgres school_db > backup_$(date +%Y%m%d).sql

# Restore if needed
psql -U postgres school_db < backup_20231201.sql
```

### 2. **Test Migrations in Development First**

```bash
# Use a separate test database
export DB_NAME=school_db_test
go run cmd/server/main.go
```

### 3. **Check Migration Status**

Add logging to see what GORM is doing:

```go
DB.AutoMigrate(...) // Enable GORM logger in config
```

### 4. **Version Control Your Models**

Always commit model changes to git before running migrations.

---

## Troubleshooting

### Issue: "relation already exists"

**Cause:** Table already exists
**Solution:** This is normal - GORM will skip creating it

### Issue: "column does not exist"

**Cause:** Old code trying to access a column that was removed
**Solution:** Update all code references before removing columns

### Issue: "foreign key constraint violation"

**Cause:** Migration order is wrong
**Solution:** Ensure parent tables are migrated before child tables

### Issue: "permission denied"

**Cause:** Database user lacks CREATE TABLE permissions
**Solution:** Grant permissions:

```sql
GRANT ALL PRIVILEGES ON DATABASE school_db TO postgres;
```

---

## Advanced: Custom Migrations

For complex schema changes, you can add custom migration logic:

```go
func RunMigrations() error {
    log.Println("🔄 Running database migrations...")

    // Standard AutoMigrate
    err := DB.AutoMigrate(...)
    if err != nil {
        return err
    }

    // Custom migrations
    if err := addCustomIndexes(); err != nil {
        return err
    }

    if err := seedDefaultData(); err != nil {
        return err
    }

    log.Println("✅ Database migrations completed successfully!")
    return nil
}

func addCustomIndexes() error {
    // Add composite index
    return DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_student_email_active
        ON students(email) WHERE deleted_at IS NULL
    `).Error
}

func seedDefaultData() error {
    // Add default departments
    var count int64
    DB.Model(&department.Department{}).Count(&count)

    if count == 0 {
        departments := []department.Department{
            {Name: "Computer Science", Description: "CS Department"},
            {Name: "Mathematics", Description: "Math Department"},
        }
        return DB.Create(&departments).Error
    }
    return nil
}
```

---

## Summary

✅ **Migrations run automatically** when you start the server  
✅ **Safe for production** - won't delete existing data  
✅ **Idempotent** - can run multiple times safely  
✅ **Order matters** - parent tables before child tables  
✅ **No manual SQL needed** - GORM handles everything

**Next Steps:**

1. Start your server: `go run cmd/server/main.go`
2. Verify tables were created: `psql -U postgres -d school_db -c "\dt"`
3. Start building your repositories and services!
