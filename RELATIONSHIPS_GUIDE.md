# GORM Relationships Guide

## Overview

All models now include proper GORM relationships using concrete types. This allows you to use GORM's powerful `Preload` feature to fetch related data efficiently.

---

## Model Relationships

### Department

```go
type Department struct {
    gorm.Model
    Name        string
    Description string

    // Relationships
    Teachers []Teacher `gorm:"foreignKey:DepartmentID"`
    Courses  []Course  `gorm:"foreignKey:DepartmentID"`
}
```

**Has Many:**

- Teachers (via `DepartmentID`)
- Courses (via `DepartmentID`)

---

### Teacher

```go
type Teacher struct {
    gorm.Model
    FirstName    string
    LastName     string
    Email        string
    Phone        string
    DepartmentID uint

    // Relationships
    Department Department `gorm:"foreignKey:DepartmentID"`
    Courses    []Course   `gorm:"foreignKey:TeacherID"`
}
```

**Belongs To:** Department  
**Has Many:** Courses

---

### Student

```go
type Student struct {
    gorm.Model
    FirstName      string
    LastName       string
    Email          string
    Phone          string
    DateOfBirth    time.Time
    EnrollmentDate time.Time

    // Relationships
    StudentCourses  []StudentCourse  `gorm:"foreignKey:StudentID"`
    StudentHomework []StudentHomework `gorm:"foreignKey:StudentID"`
    Attendances     []Attendance     `gorm:"foreignKey:StudentID"`
    Grades          []Grade          `gorm:"foreignKey:StudentID"`
}
```

**Has Many:**

- StudentCourses (enrollments)
- StudentHomework (submissions)
- Attendances
- Grades

---

### Course

```go
type Course struct {
    gorm.Model
    Name         string
    Code         string
    Description  string
    Credits      int
    DepartmentID uint
    TeacherID    uint

    // Relationships
    Department     Department     `gorm:"foreignKey:DepartmentID"`
    Teacher        Teacher        `gorm:"foreignKey:TeacherID"`
    StudentCourses []StudentCourse `gorm:"foreignKey:CourseID"`
    Attendances    []Attendance   `gorm:"foreignKey:CourseID"`
    Homework       []Homework     `gorm:"foreignKey:CourseID"`
    Exams          []Exam         `gorm:"foreignKey:CourseID"`
}
```

**Belongs To:** Department, Teacher  
**Has Many:** StudentCourses, Attendances, Homework, Exams

---

### Attendance

```go
type Attendance struct {
    gorm.Model
    StudentID uint
    CourseID  uint
    Date      time.Time
    Status    AttendanceStatus

    // Relationships
    Student Student `gorm:"foreignKey:StudentID"`
    Course  Course  `gorm:"foreignKey:CourseID"`
}
```

**Belongs To:** Student, Course

---

### Homework

```go
type Homework struct {
    gorm.Model
    Title       string
    Description string
    CourseID    uint
    DueDate     time.Time
    MaxScore    float64

    // Relationships
    Course          Course           `gorm:"foreignKey:CourseID"`
    StudentHomework []StudentHomework `gorm:"foreignKey:HomeworkID"`
}
```

**Belongs To:** Course  
**Has Many:** StudentHomework (submissions)

---

### Exam

```go
type Exam struct {
    gorm.Model
    Title    string
    CourseID uint
    ExamDate time.Time
    Duration int
    MaxScore float64

    // Relationships
    Course Course  `gorm:"foreignKey:CourseID"`
    Grades []Grade `gorm:"foreignKey:ExamID"`
}
```

**Belongs To:** Course  
**Has Many:** Grades

---

### Grade

```go
type Grade struct {
    gorm.Model
    StudentID uint
    ExamID    uint
    Score     float64

    // Relationships
    Student Student `gorm:"foreignKey:StudentID"`
    Exam    Exam    `gorm:"foreignKey:ExamID"`
}
```

**Belongs To:** Student, Exam

---

### StudentCourse (Junction Table)

```go
type StudentCourse struct {
    gorm.Model
    StudentID      uint
    CourseID       uint
    EnrollmentDate time.Time

    // Relationships
    Student Student `gorm:"foreignKey:StudentID"`
    Course  Course  `gorm:"foreignKey:CourseID"`
}
```

**Belongs To:** Student, Course

---

### StudentHomework (Junction Table)

```go
type StudentHomework struct {
    gorm.Model
    StudentID      uint
    HomeworkID     uint
    SubmissionDate *time.Time
    Score          *float64
    Status         HomeworkStatus

    // Relationships
    Student  Student  `gorm:"foreignKey:StudentID"`
    Homework Homework `gorm:"foreignKey:HomeworkID"`
}
```

**Belongs To:** Student, Homework

---

## Using Relationships in Queries

### Basic Preload

Load a student with all their courses:

```go
var student student.Student
db.Preload("StudentCourses").First(&student, 1)
```

### Nested Preload

Load a student with courses and each course's details:

```go
var student student.Student
db.Preload("StudentCourses.Course").First(&student, 1)
```

### Multiple Preloads

Load a student with all related data:

```go
var student student.Student
db.Preload("StudentCourses").
   Preload("StudentHomework").
   Preload("Attendances").
   Preload("Grades").
   First(&student, 1)
```

### Preload with Conditions

Load only active courses for a student:

```go
var student student.Student
db.Preload("StudentCourses", "enrollment_date > ?", time.Now().AddDate(0, -6, 0)).
   First(&student, 1)
```

### Deep Nested Preload

Load a course with department, teacher, and all enrolled students:

```go
var course course.Course
db.Preload("Department").
   Preload("Teacher").
   Preload("StudentCourses.Student").
   First(&course, 1)
```

---

## Example Repository Methods

### Get Student with All Enrollments

```go
func (r *StudentRepository) GetStudentWithCourses(studentID uint) (*student.Student, error) {
    var s student.Student
    err := r.db.Preload("StudentCourses.Course").First(&s, studentID).Error
    return &s, err
}
```

### Get Course with Teacher and Department

```go
func (r *CourseRepository) GetCourseDetails(courseID uint) (*course.Course, error) {
    var c course.Course
    err := r.db.
        Preload("Department").
        Preload("Teacher").
        First(&c, courseID).Error
    return &c, err
}
```

### Get Student's Grades with Exam Details

```go
func (r *GradeRepository) GetStudentGrades(studentID uint) ([]grade.Grade, error) {
    var grades []grade.Grade
    err := r.db.
        Preload("Exam.Course").
        Where("student_id = ?", studentID).
        Find(&grades).Error
    return grades, err
}
```

### Get Department with All Teachers and Courses

```go
func (r *DepartmentRepository) GetDepartmentFull(deptID uint) (*department.Department, error) {
    var dept department.Department
    err := r.db.
        Preload("Teachers").
        Preload("Courses").
        First(&dept, deptID).Error
    return &dept, err
}
```

---

## Important Notes

### Placeholder Types

Each model file contains placeholder type definitions to avoid circular import issues:

```go
// Placeholder types
type Student struct{}
type Course struct{}
```

**These are intentional!** They allow the model to compile without importing other packages. When you query with `Preload`, GORM will use the actual types from the database.

### Avoiding N+1 Queries

Always use `Preload` to fetch related data in a single query:

❌ **Bad (N+1 queries):**

```go
var students []student.Student
db.Find(&students)
for _, s := range students {
    // This makes a separate query for each student!
    db.Model(&s).Association("StudentCourses").Find(&s.StudentCourses)
}
```

✅ **Good (2 queries total):**

```go
var students []student.Student
db.Preload("StudentCourses").Find(&students)
```

### Omitempty in JSON

All relationship fields use `json:"...,omitempty"` which means:

- They won't appear in JSON if not loaded
- You must explicitly `Preload` them to include in API responses

### Circular References

Be careful with circular references in JSON serialization:

```go
// This could cause infinite loop in JSON marshaling
var course course.Course
db.Preload("StudentCourses.Student.StudentCourses").First(&course, 1)
```

**Solution:** Use DTOs to control what gets serialized.

---

## Testing Relationships

### Verify Foreign Keys Exist

```sql
-- Check foreign keys in PostgreSQL
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
  ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
  ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY';
```

### Test Preload in Code

```go
// Test that preload works
var student student.Student
err := db.Preload("StudentCourses").First(&student, 1).Error
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Student: %s %s\n", student.FirstName, student.LastName)
fmt.Printf("Enrolled in %d courses\n", len(student.StudentCourses))
```

---

## Summary

✅ **All models have proper GORM relationships**  
✅ **Use concrete types (not `interface{}`)**  
✅ **Placeholder types prevent circular imports**  
✅ **Use `Preload` to fetch related data**  
✅ **Avoid N+1 query problems**  
✅ **Relationships work with migrations**

**Next Steps:**

1. Implement repository methods with `Preload`
2. Create DTOs for API responses
3. Test relationship queries
4. Build service layer with business logic
