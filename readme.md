# 📅 Employee Scheduling System

**Advanced Programming Languages(MSCS-632) - Assignment 4: Implementing Control Structures**

This application demonstrates mastery of control structures (conditionals, loops, branching) across multiple programming languages. It manages employee schedules with conflict resolution, implementing complex scheduling logic that ensures proper shift coverage while respecting employee preferences and business constraints.

## 📁 File Structure

```
employee-scheduler/
├── README.md
├── go/
│   ├── employee.go     # Employee struct and methods
│   ├── scheduler.go    # Scheduling logic
│   └── main.go         # Main application
└── python/
    ├── employee.py     # Employee class
    ├── scheduler.py    # Scheduling logic
    └── main.py         # Main application
```

## 🚀 How to Run

### Go Version
```bash
cd go/
go run main.go scheduler.go employee.go
```

### Python Version
```bash
cd python/
python main.py
```

## 🎮 How to Use

### 1. **Add Employees**
```
👤 Enter employee name: Sarah
📅 Setting up preferences for Sarah
🌅 Shifts: 0=Morning | ☀️  1=Afternoon | 🌙 2=Evening
💡 Press Enter to skip a day
Monday    : 0
   ✅ Set 🌅 Morning preference
Tuesday   : 0
   ✅ Set 🌅 Morning preference
```

### 2. **Generate Schedule**
```
🔄 Generating schedule...
✅ Step 1: Reset all schedules
🔄 Conflict resolved: TestEmp9 → Monday Afternoon (preferred shift full)
✅ Step 2: Assigned preferred shifts
✅ Step 3: Ensured minimum staffing
🎉 Schedule generated successfully!
```

### 3. **View Schedule**
```
📅                         WEEKLY EMPLOYEE SCHEDULE                          📅

📅 MONDAY
   🌅 Morning    | Sarah, Mike, Jennifer, Tom, Nicole        ✅ STAFFED (5/8)
   ☀️  Afternoon  | David, Lisa, Rachel, Daniel               ✅ STAFFED (4/8)
   🌙 Evening     | Amanda, Michelle                          ✅ STAFFED (2/8)
```

### Menu Options
1. **Add Employee** - Add employees with shift preferences
2. **Generate Schedule** - Create weekly schedule with conflict resolution
3. **View Schedule** - Display beautiful formatted schedule
4. **Exit**

### Business Rules
- Max 5 days per employee per week
- Min 2 employees per shift
- Max 8 employees per shift
- One shift per day per employee

That's it! Choose your preferred language and start scheduling! 🎉