package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Scheduler manages the overall employee scheduling system
type Scheduler struct {
	Employees []*Employee                         // List of all employees
	Schedule  map[time.Weekday]map[Shift][]string // day -> shift -> employee names
	Days      []time.Weekday                      // Days of the week
	Shifts    []Shift                             // Available shifts
}

// NewScheduler creates a new scheduler with empty schedule
func NewScheduler() *Scheduler {
	days := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	}
	shifts := []Shift{Morning, Afternoon, Evening}

	// Initialize empty schedule for all days and shifts
	schedule := make(map[time.Weekday]map[Shift][]string)
	for _, day := range days {
		schedule[day] = make(map[Shift][]string)
		for _, shift := range shifts {
			schedule[day][shift] = []string{}
		}
	}

	return &Scheduler{
		Employees: []*Employee{},
		Schedule:  schedule,
		Days:      days,
		Shifts:    shifts,
	}
}

// AddEmployee interactively adds a new employee with their shift preferences
func (s *Scheduler) AddEmployee() {
	scanner := bufio.NewScanner(os.Stdin)

	// Get and validate employee name
	fmt.Print("\n👤 Enter employee name: ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())

	if name == "" {
		fmt.Println("❌ Name cannot be empty!")
		return
	}

	// Check for duplicate names (case-insensitive)
	for _, emp := range s.Employees {
		if strings.EqualFold(emp.Name, name) {
			fmt.Printf("❌ Employee '%s' already exists. Please use a different name.\n", name)
			return
		}
	}

	employee := NewEmployee(name)

	// Collect shift preferences for each day of the week
	fmt.Printf("\n📅 Setting up preferences for %s\n", employee.Name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🌅 Shifts: 0=Morning | ☀️  1=Afternoon | 🌙 2=Evening")
	fmt.Println("💡 Press Enter to skip a day")

	for _, day := range s.Days {
		for {
			fmt.Printf("%-10s: ", day)
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())

			if input == "" {
				break // No preference for this day
			}

			// Parse and validate shift preference
			preference, err := strconv.Atoi(input)
			if err != nil || preference < 0 || preference >= len(s.Shifts) {
				fmt.Println("   ❌ Invalid input. Please enter 0, 1, or 2 (or press Enter to skip)")
				continue
			}

			shiftIcon := getShiftIcon(s.Shifts[preference])
			employee.SetPreference(day, s.Shifts[preference])
			fmt.Printf("   ✅ Set %s %s preference\n", shiftIcon, s.Shifts[preference])
			break
		}
	}

	s.Employees = append(s.Employees, employee)
	fmt.Printf("\n🎉 Employee %s added successfully!\n", name)
	fmt.Printf("📊 Total employees: %d\n\n", len(s.Employees))
}

// PrintSchedule displays the complete weekly schedule and individual employee summaries
func (s *Scheduler) PrintSchedule() {
	if len(s.Employees) == 0 {
		fmt.Println("\n📝 No employees added yet! Please add some employees first.")
		return
	}

	s.printScheduleHeader()
	s.printWeeklyGrid()
	s.printScheduleFooter()
	s.printEmployeeSummaries()
}

// printScheduleHeader displays the main schedule title with decorative formatting
func (s *Scheduler) printScheduleHeader() {
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("📅                         WEEKLY EMPLOYEE SCHEDULE                          📅")
	fmt.Println(strings.Repeat("═", 90))
}

// printWeeklyGrid displays the main schedule in a clean tabular format
func (s *Scheduler) printWeeklyGrid() {
	for _, day := range s.Days {
		fmt.Printf("\n📅 %s\n", day)
		fmt.Println(strings.Repeat("─", 85))

		for _, shift := range s.Shifts {
			employees := s.Schedule[day][shift]
			staffCount := len(employees)

			// Format shift name with emoji
			shiftIcon := getShiftIcon(shift)
			fmt.Printf("   %s %-10s │ ", shiftIcon, shift)

			if staffCount == 0 {
				fmt.Printf("%-50s", "No employees assigned")
				if staffCount < MinEmployeesPerShift {
					fmt.Printf(" 🚨 UNDERSTAFFED (need %d)", MinEmployeesPerShift)
				}
			} else {
				employeeList := strings.Join(employees, ", ")
				if len(employeeList) > 45 {
					employeeList = employeeList[:42] + "..."
				}
				fmt.Printf("%-50s", employeeList)

				// Status indicator
				if staffCount < MinEmployeesPerShift {
					fmt.Printf(" 🚨 UNDERSTAFFED (%d/%d)", staffCount, MinEmployeesPerShift)
				} else if staffCount >= MinEmployeesPerShift && staffCount < MaxEmployeesPerShift {
					fmt.Printf(" ✅ STAFFED (%d/%d)", staffCount, MaxEmployeesPerShift)
				} else {
					fmt.Printf(" 🏆 FULL (%d/%d)", staffCount, MaxEmployeesPerShift)
				}
			}
			fmt.Println()
		}
	}
}

// printScheduleFooter displays summary statistics
func (s *Scheduler) printScheduleFooter() {
	fmt.Println("\n" + strings.Repeat("═", 90))

	totalShifts := len(s.Days) * len(s.Shifts)
	staffedShifts := 0
	fullShifts := 0
	totalAssignments := 0

	for _, day := range s.Days {
		for _, shift := range s.Shifts {
			count := len(s.Schedule[day][shift])
			totalAssignments += count
			if count >= MinEmployeesPerShift {
				staffedShifts++
			}
			if count == MaxEmployeesPerShift {
				fullShifts++
			}
		}
	}

	fmt.Printf("📊 SCHEDULE STATS: %d/%d shifts properly staffed │ %d full shifts │ %d total assignments\n",
		staffedShifts, totalShifts, fullShifts, totalAssignments)
	fmt.Println(strings.Repeat("═", 90))
}

// printEmployeeSummaries displays individual employee work schedules in a clean format
func (s *Scheduler) printEmployeeSummaries() {
	fmt.Println("\n👥 EMPLOYEE WORK SUMMARIES")
	fmt.Println(strings.Repeat("─", 60))

	for i, employee := range s.Employees {
		fmt.Print(employee.GetWorkSummary())

		// Add spacing between employees except for the last one
		if i < len(s.Employees)-1 {
			fmt.Println()
		}
	}
	fmt.Printf("\n" + strings.Repeat("─", 60) + "\n")
}

// getShiftIcon returns an emoji icon for each shift type
func getShiftIcon(shift Shift) string {
	switch shift {
	case Morning:
		return "🌅"
	case Afternoon:
		return "☀️"
	case Evening:
		return "🌙"
	default:
		return "⏰"
	}
}

// GetEmployeeCount returns the total number of employees in the system
func (s *Scheduler) GetEmployeeCount() int {
	return len(s.Employees)
}

// AssignShifts executes the complete scheduling algorithm in three phases:
// 1. Reset all schedules, 2. Assign preferences with conflict resolution, 3. Fill minimum staffing
func (s *Scheduler) AssignShifts() {
	fmt.Println("\n🔄 Generating schedule...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	s.resetSchedules()
	fmt.Println("✅ Step 1: Reset all schedules")

	s.assignPreferredShifts()
	fmt.Println("✅ Step 2: Assigned preferred shifts")

	s.ensureMinimumStaffing()
	fmt.Println("✅ Step 3: Ensured minimum staffing")

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 Schedule generated successfully!")
	fmt.Printf("📋 %d employees scheduled across %d days\n", len(s.Employees), len(s.Days))
	fmt.Println("💡 Use 'View Schedule' to see the complete weekly schedule.")
}

// resetSchedules clears all employee schedules and the main schedule grid to start fresh
func (s *Scheduler) resetSchedules() {
	// Reset each employee's individual schedule and work day counter
	for _, employee := range s.Employees {
		employee.ResetSchedule()
	}

	// Clear all shifts in the main schedule
	for _, day := range s.Days {
		for _, shift := range s.Shifts {
			s.Schedule[day][shift] = []string{}
		}
	}
}

// assignPreferredShifts attempts to assign each employee to their preferred shifts,
// with automatic conflict resolution when preferred shifts are full
func (s *Scheduler) assignPreferredShifts() {
	for _, employee := range s.Employees {
		for _, day := range s.Days {
			// Stop assigning if employee has reached maximum work days (5 per week)
			if employee.DaysWorked >= MaxWorkDaysPerWeek {
				break
			}

			// Check if employee has a preference for this day
			if preferredShift, hasPreference := employee.GetPreference(day); hasPreference {
				if s.canAssign(employee, day, preferredShift) {
					s.assign(employee, day, preferredShift)
				} else {
					// Preferred shift is full or employee can't work - try to resolve conflict
					s.resolveConflict(employee, day, preferredShift)
				}
			}
		}
	}
}

// resolveConflict attempts to find alternative assignments when an employee's preferred shift is unavailable.
// Strategy: 1) Try other shifts same day, 2) Try preferred shift other days, 3) Try any shift other days
func (s *Scheduler) resolveConflict(employee *Employee, preferredDay time.Weekday, preferredShift Shift) {
	// Strategy 1: Try alternative shifts on the same day
	for _, otherShift := range s.Shifts {
		if otherShift != preferredShift {
			if s.canAssign(employee, preferredDay, otherShift) {
				s.assign(employee, preferredDay, otherShift)
				fmt.Printf("🔄 Conflict resolved: %s → %s %s (preferred shift full)\n",
					employee.Name, preferredDay, otherShift)
				return
			}
		}
	}

	// Strategy 2: Try other days (preferred shift first, then any shift)
	for _, otherDay := range s.Days {
		if otherDay == preferredDay {
			continue // Skip the day we already tried
		}

		// Try preferred shift on this alternative day
		if s.canAssign(employee, otherDay, preferredShift) {
			s.assign(employee, otherDay, preferredShift)
			fmt.Printf("🔄 Conflict resolved: %s → %s %s (moved to different day)\n",
				employee.Name, otherDay, preferredShift)
			return
		}

		// Try any available shift on this alternative day
		for _, otherShift := range s.Shifts {
			if s.canAssign(employee, otherDay, otherShift) {
				s.assign(employee, otherDay, otherShift)
				fmt.Printf("🔄 Conflict resolved: %s → %s %s (alternative assignment)\n",
					employee.Name, otherDay, otherShift)
				return
			}
		}
	}

	// No resolution found - employee cannot be assigned anywhere
	fmt.Printf("⚠️  Warning: Could not assign %s anywhere (schedule full)\n", employee.Name)
}

// canAssign checks if an employee can be assigned to a specific day and shift
// Validates: employee availability, work day limits, and shift capacity limits
func (s *Scheduler) canAssign(employee *Employee, day time.Weekday, shift Shift) bool {
	// Check if employee can work this day (not already scheduled, under 5-day limit)
	if !employee.CanWorkDay(day) {
		return false
	}

	// Check if shift has reached maximum capacity (8 employees per shift)
	if len(s.Schedule[day][shift]) >= MaxEmployeesPerShift {
		return false
	}

	return true
}

// assign adds an employee to a specific shift and updates both the employee's schedule
// and the main schedule grid
func (s *Scheduler) assign(employee *Employee, day time.Weekday, shift Shift) {
	employee.AssignShift(day, shift)
	s.Schedule[day][shift] = append(s.Schedule[day][shift], employee.Name)
}

// ensureMinimumStaffing fills any shifts that don't meet the minimum staffing requirement (2 employees)
// by randomly assigning available employees who haven't reached their work day limit
func (s *Scheduler) ensureMinimumStaffing() {
	for _, day := range s.Days {
		for _, shift := range s.Shifts {
			currentStaff := len(s.Schedule[day][shift])

			// Check if this shift needs more employees to meet minimum requirement
			if currentStaff < MinEmployeesPerShift {
				needed := MinEmployeesPerShift - currentStaff
				availableEmployees := s.getAvailableEmployees(day)

				// Randomly shuffle available employees for fair distribution
				rand.Seed(time.Now().UnixNano())
				rand.Shuffle(len(availableEmployees), func(i, j int) {
					availableEmployees[i], availableEmployees[j] = availableEmployees[j], availableEmployees[i]
				})

				// Assign employees until minimum is met or no more employees available
				assigned := 0
				for _, employee := range availableEmployees {
					if assigned >= needed {
						break
					}

					if s.canAssign(employee, day, shift) {
						s.assign(employee, day, shift)
						assigned++
					}
				}

				// Warn if unable to meet minimum staffing despite trying
				if assigned < needed {
					fmt.Printf("⚠️  Understaffed: %s %s needs %d more employees (have %d/%d)\n",
						day, shift, needed-assigned, currentStaff+assigned, MinEmployeesPerShift)
				}
			}
		}
	}
}

// getAvailableEmployees returns a list of employees who can work on the specified day.
// Available means: not already scheduled for that day and under the 5-day weekly limit
func (s *Scheduler) getAvailableEmployees(day time.Weekday) []*Employee {
	var available []*Employee
	for _, employee := range s.Employees {
		if employee.CanWorkDay(day) {
			available = append(available, employee)
		}
	}
	return available
}

// // GenerateRealisticTest creates 20 employees with diverse, realistic shift preferences
// func (s *Scheduler) GenerateRealisticTest() {
// 	s.Employees = []*Employee{}

// 	// Define realistic employee profiles with varied preferences
// 	employees := []struct {
// 		name  string
// 		prefs map[time.Weekday]Shift
// 	}{
// 		// Early birds who prefer morning shifts
// 		{"Sarah", map[time.Weekday]Shift{
// 			time.Monday: Morning, time.Tuesday: Morning, time.Wednesday: Morning,
// 			time.Thursday: Morning, time.Friday: Morning}},
// 		{"Mike", map[time.Weekday]Shift{
// 			time.Monday: Morning, time.Wednesday: Morning, time.Friday: Morning, time.Saturday: Morning}},
// 		{"Jennifer", map[time.Weekday]Shift{
// 			time.Tuesday: Morning, time.Thursday: Morning, time.Saturday: Morning, time.Sunday: Morning}},

// 		// Afternoon preference workers
// 		{"David", map[time.Weekday]Shift{
// 			time.Monday: Afternoon, time.Tuesday: Afternoon, time.Wednesday: Afternoon,
// 			time.Thursday: Afternoon, time.Friday: Afternoon}},
// 		{"Lisa", map[time.Weekday]Shift{
// 			time.Monday: Afternoon, time.Wednesday: Afternoon, time.Friday: Afternoon}},
// 		{"Robert", map[time.Weekday]Shift{
// 			time.Tuesday: Afternoon, time.Thursday: Afternoon, time.Saturday: Afternoon, time.Sunday: Afternoon}},

// 		// Night owls who prefer evening shifts
// 		{"Amanda", map[time.Weekday]Shift{
// 			time.Monday: Evening, time.Tuesday: Evening, time.Wednesday: Evening,
// 			time.Thursday: Evening, time.Friday: Evening}},
// 		{"James", map[time.Weekday]Shift{
// 			time.Friday: Evening, time.Saturday: Evening, time.Sunday: Evening}},
// 		{"Michelle", map[time.Weekday]Shift{
// 			time.Monday: Evening, time.Wednesday: Evening, time.Saturday: Evening}},

// 		// Weekend specialists
// 		{"Carlos", map[time.Weekday]Shift{
// 			time.Saturday: Morning, time.Sunday: Morning}},
// 		{"Emma", map[time.Weekday]Shift{
// 			time.Saturday: Afternoon, time.Sunday: Afternoon, time.Friday: Evening}},
// 		{"Alex", map[time.Weekday]Shift{
// 			time.Saturday: Evening, time.Sunday: Evening, time.Friday: Afternoon}},

// 		// Flexible workers with mixed preferences
// 		{"Jessica", map[time.Weekday]Shift{
// 			time.Monday: Morning, time.Wednesday: Afternoon, time.Friday: Evening, time.Sunday: Morning}},
// 		{"Kevin", map[time.Weekday]Shift{
// 			time.Tuesday: Morning, time.Thursday: Afternoon, time.Saturday: Evening}},
// 		{"Rachel", map[time.Weekday]Shift{
// 			time.Monday: Afternoon, time.Thursday: Morning, time.Sunday: Evening}},

// 		// Part-time workers with minimal preferences
// 		{"Tom", map[time.Weekday]Shift{time.Monday: Morning, time.Friday: Morning}},
// 		{"Anna", map[time.Weekday]Shift{time.Wednesday: Afternoon, time.Sunday: Afternoon}},
// 		{"Chris", map[time.Weekday]Shift{time.Tuesday: Evening, time.Saturday: Evening}},

// 		// Highly available workers
// 		{"Nicole", map[time.Weekday]Shift{
// 			time.Monday: Morning, time.Tuesday: Afternoon, time.Wednesday: Evening,
// 			time.Thursday: Morning, time.Friday: Afternoon}},
// 		{"Daniel", map[time.Weekday]Shift{
// 			time.Monday: Afternoon, time.Tuesday: Morning, time.Wednesday: Morning,
// 			time.Thursday: Evening, time.Friday: Morning}},
// 	}

// 	// Create employees with their preferences
// 	for _, emp := range employees {
// 		employee := NewEmployee(emp.name)
// 		for day, shift := range emp.prefs {
// 			employee.SetPreference(day, shift)
// 		}
// 		s.Employees = append(s.Employees, employee)
// 	}

// 	fmt.Println("🎭 Generated 20 employees with realistic preferences:")
// 	fmt.Println("   • 3 Early Birds (morning lovers)")
// 	fmt.Println("   • 3 Day Shifters (afternoon preferred)")
// 	fmt.Println("   • 3 Night Owls (evening specialists)")
// 	fmt.Println("   • 3 Weekend Warriors (Saturday/Sunday focus)")
// 	fmt.Println("   • 3 Flexible Workers (mixed shifts)")
// 	fmt.Println("   • 3 Part-timers (minimal preferences)")
// 	fmt.Println("   • 2 Workaholics (maximum availability)")
// 	fmt.Println("")
// 	fmt.Println("🎯 Expected results:")
// 	fmt.Println("   ✅ Most shifts should be properly staffed")
// 	fmt.Println("   ✅ Good mix of satisfied and conflict-resolved preferences")
// 	fmt.Println("   ✅ Realistic distribution across all days and shifts")
// 	fmt.Println("   ✅ Minimal understaffing warnings")
// }
