from enum import Enum
from typing import Dict, Optional

class Shift(Enum):
    MORNING = "Morning"
    AFTERNOON = "Afternoon"
    EVENING = "Evening"

    def _get_shift_icon(self) -> str:
        """Returns an emoji icon for a given shift"""
        if self == Shift.MORNING:
            return "🌅"
        elif self == Shift.AFTERNOON:
            return "☀️"
        elif self == Shift.EVENING:
            return "🌙"
        else:
            return "❓"

# Scheduling constraints
MIN_EMPLOYEES_PER_SHIFT = 2
MAX_EMPLOYEES_PER_SHIFT = 8
MAX_WORK_DAYS_PER_WEEK = 5

class Employee:
    def __init__(self, name: str):
        self.name = name
        self.preferences: Dict[int, Shift] = {}  # weekday -> shift
        self.schedule: Dict[int, Shift] = {}     # weekday -> shift
        self.days_worked = 0

    def set_preference(self, day: int, shift: Shift):
        """Sets the preferred shift for a specific day (0=Monday, 6=Sunday)"""
        self.preferences[day] = shift

    def get_preference(self, day: int) -> Optional[Shift]:
        """Returns the preferred shift for a given day"""
        return self.preferences.get(day)

    def can_work_day(self, day: int) -> bool:
        """Checks if employee is available to work on a given day"""
        # Can't work if already assigned a shift that day
        if day in self.schedule:
            return False
        # Can't work if already at maximum days per week
        return self.days_worked < MAX_WORK_DAYS_PER_WEEK

    def assign_shift(self, day: int, shift: Shift) -> bool:
        """Assigns a shift to the employee for a specific day"""
        if not self.can_work_day(day):
            return False
        self.schedule[day] = shift
        self.days_worked += 1
        return True

    def get_assigned_shift(self, day: int) -> Optional[Shift]:
        """Returns the assigned shift for a given day"""
        return self.schedule.get(day)

    def has_preference_match(self, day: int) -> bool:
        """Checks if assigned shift matches preference for a day"""
        assigned_shift = self.get_assigned_shift(day)
        preferred_shift = self.get_preference(day)
        return (assigned_shift is not None and
                preferred_shift is not None and
                assigned_shift == preferred_shift)

    def reset_schedule(self):
        """Clears all shift assignments"""
        self.schedule = {}
        self.days_worked = 0

    def is_duplicate_name(self, name: str) -> bool:
        """Checks if the name matches (case-insensitive)"""
        return self.name.lower() == name.lower()

    def get_work_summary(self) -> str:
        """Returns a beautifully formatted string of the employee's work schedule"""
        workload_icon = self._get_workload_icon()
        summary = f"\n{workload_icon} {self.name} ({self.days_worked}/5 days)\n"

        days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]

        has_assignments = False
        for day_num, day_name in enumerate(days):
            assigned_shift = self.get_assigned_shift(day_num)
            if assigned_shift is not None:
                has_assignments = True
                shift_icon = assigned_shift._get_shift_icon()

                if self.has_preference_match(day_num):
                    summary += f"   {shift_icon} {day_name:<10} → {assigned_shift.value} ✨ PREFERRED\n"
                elif self.get_preference(day_num) is not None:
                    preferred = self.get_preference(day_num).value
                    summary += f"   {shift_icon} {day_name:<10} → {assigned_shift.value} (wanted {preferred})\n"
                else:
                    summary += f"   {shift_icon} {day_name:<10} → {assigned_shift.value}\n"

        if not has_assignments:
            summary += "   💤 No shifts assigned\n"

        return summary

    def _get_workload_icon(self) -> str:
        """Returns an emoji indicating employee work load"""
        if self.days_worked == 0:
            return "😴"
        elif self.days_worked <= 2:
            return "😌"
        elif self.days_worked <= 4:
            return "😊"
        elif self.days_worked == 5:
            return "💪"
        else:
            return "🤯"