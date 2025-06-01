package main

import (
	"fmt"
	"time"
)

// Shift represents the time periods employees can work
type Shift string

const (
	Morning   Shift = "Morning"
	Afternoon Shift = "Afternoon"
	Evening   Shift = "Evening"
)

// Scheduling constraints
const (
	MinEmployeesPerShift = 2 // Minimum employees required per shift
	MaxEmployeesPerShift = 8 // Maximum employees allowed per shift
	MaxWorkDaysPerWeek   = 5 // Maximum days an employee can work per week
)

// Employee represents an individual employee with their preferences and assigned schedule
type Employee struct {
	Name       string                 // Employee's name
	Preference map[time.Weekday]Shift // Preferred shifts for each day
	Schedule   map[time.Weekday]Shift // Assigned shifts for each day
	DaysWorked int                    // Number of days currently assigned
}

// NewEmployee creates a new employee with empty preferences and schedule
func NewEmployee(name string) *Employee {
	return &Employee{
		Name:       name,
		Preference: make(map[time.Weekday]Shift),
		Schedule:   make(map[time.Weekday]Shift),
		DaysWorked: 0,
	}
}

// SetPreference sets the preferred shift for a specific day
func (e *Employee) SetPreference(day time.Weekday, shift Shift) {
	e.Preference[day] = shift
}

// GetPreference returns the preferred shift for a given day
func (e *Employee) GetPreference(day time.Weekday) (Shift, bool) {
	shift, exists := e.Preference[day]
	return shift, exists
}

// CanWorkDay checks if employee is available to work on a given day
func (e *Employee) CanWorkDay(day time.Weekday) bool {
	// Can't work if already assigned a shift that day
	if _, alreadyScheduled := e.Schedule[day]; alreadyScheduled {
		return false
	}
	// Can't work if already at maximum days per week
	return e.DaysWorked < MaxWorkDaysPerWeek
}

// AssignShift assigns a shift to the employee for a specific day
func (e *Employee) AssignShift(day time.Weekday, shift Shift) bool {
	if !e.CanWorkDay(day) {
		return false
	}
	e.Schedule[day] = shift
	e.DaysWorked++
	return true
}

// RemoveShift removes a shift assignment for a specific day
func (e *Employee) RemoveShift(day time.Weekday) bool {
	if _, exists := e.Schedule[day]; exists {
		delete(e.Schedule, day)
		e.DaysWorked--
		return true
	}
	return false
}

// GetAssignedShift returns the assigned shift for a given day
func (e *Employee) GetAssignedShift(day time.Weekday) (Shift, bool) {
	shift, exists := e.Schedule[day]
	return shift, exists
}

// HasPreferenceMatch checks if assigned shift matches preference for a day
func (e *Employee) HasPreferenceMatch(day time.Weekday) bool {
	assignedShift, hasAssignment := e.GetAssignedShift(day)
	preferredShift, hasPreference := e.GetPreference(day)
	return hasAssignment && hasPreference && assignedShift == preferredShift
}

// ResetSchedule clears all shift assignments
func (e *Employee) ResetSchedule() {
	e.Schedule = make(map[time.Weekday]Shift)
	e.DaysWorked = 0
}

// GetWorkSummary returns a beautifully formatted string of the employee's work schedule
func (e *Employee) GetWorkSummary() string {
	// Employee header with work load indicator
	workLoadIcon := e.getWorkLoadIcon()
	summary := fmt.Sprintf("\n%s %s (%d/5 days)\n", workLoadIcon, e.Name, e.DaysWorked)

	days := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	}

	hasAssignments := false
	// Show each day's assignment with preference matching indicators
	for _, day := range days {
		if assignedShift, hasAssignment := e.GetAssignedShift(day); hasAssignment {
			hasAssignments = true
			shiftIcon := getShiftIcon(assignedShift)

			// Check if preference was matched
			if e.HasPreferenceMatch(day) {
				summary += fmt.Sprintf("   %s %-10s → %s ✨ PREFERRED\n", shiftIcon, day, assignedShift)
			} else if preferredShift, hasPreference := e.GetPreference(day); hasPreference {
				summary += fmt.Sprintf("   %s %-10s → %s (wanted %s)\n", shiftIcon, day, assignedShift, preferredShift)
			} else {
				summary += fmt.Sprintf("   %s %-10s → %s\n", shiftIcon, day, assignedShift)
			}
		}
	}

	if !hasAssignments {
		summary += "   💤 No shifts assigned\n"
	}

	return summary
}

// getWorkLoadIcon returns an emoji indicating employee work load
func (e *Employee) getWorkLoadIcon() string {
	switch {
	case e.DaysWorked == 0:
		return "😴"
	case e.DaysWorked <= 2:
		return "😌"
	case e.DaysWorked <= 4:
		return "😊"
	case e.DaysWorked == 5:
		return "💪"
	default:
		return "🤯"
	}
}