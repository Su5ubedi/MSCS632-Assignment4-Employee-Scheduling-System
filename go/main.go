package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// main function runs the employee scheduling application
func main() {
	scheduler := NewScheduler()
	scanner := bufio.NewScanner(os.Stdin)

	// Display application header and constraints
	fmt.Println("Employee Scheduling System")
	fmt.Println("=========================")
	fmt.Printf("Constraints: Min %d, Max %d employees per shift | Max %d days per employee\n\n",
		MinEmployeesPerShift, MaxEmployeesPerShift, MaxWorkDaysPerWeek)

	// Main application loop
	for {
		// Display menu options
		fmt.Println("1. Add Employee")
		fmt.Println("2. Generate Schedule")
		fmt.Println("3. View Schedule")
		fmt.Println("4. Exit")
		fmt.Print("Choose an option: ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		// Handle user menu selection
		switch choice {
		case "1":
			scheduler.AddEmployee()
		case "2":
			// Ensure employees exist before generating schedule
			if scheduler.GetEmployeeCount() == 0 {
				fmt.Println("Please add employees first!")
				fmt.Println()
				continue
			}
			scheduler.AssignShifts()
		case "3":
			scheduler.PrintSchedule()
		case "4":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
