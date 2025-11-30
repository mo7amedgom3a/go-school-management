#!/bin/bash

# === Configuration ===
PROJECT_NAME="school_management"

echo "📦 Initializing Go Module"
go mod init $MODULE_NAME

echo "📁 Creating folder structure..."

# === Create CMD Folder ===
mkdir -p cmd/server

# === Application Core Structure ===
mkdir -p internal/{config,database,server}

# === Modules ===
MODULES=(
        "student"
        "teacher"
        "course"
        "department"
        "attendance"
        "homework"
        "students_homework"
        "student_courses"
        "exam"
        "grade"
)

mkdir -p internal/modules

for module in "${MODULES[@]}"; do
        echo "📦 Creating module: $module"
        mkdir -p internal/modules/$module

        touch internal/modules/$module/${module}_model.go
        touch internal/modules/$module/${module}_dto.go
        touch internal/modules/$module/${module}_repository.go
        touch internal/modules/$module/${module}_service.go
        touch internal/modules/$module/${module}_controller.go
done

# === pkg Folder (shared utils) ===
mkdir -p pkg/{logger,response,validation,middleware}
