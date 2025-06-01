import random
from typing import Dict, List
from collections import defaultdict
from employee import Employee, Shift, MIN_EMPLOYEES_PER_SHIFT, MAX_EMPLOYEES_PER_SHIFT, MAX_WORK_DAYS_PER_WEEK

class Scheduler:
    def __init__(self):
        self.employees: List[Employee] = []
        self.schedule: Dict[int, Dict[Shift, List[str]]] = defaultdict(lambda: defaultdict(list))
        self.days = list(range(7))  # 0=Monday, 6=Sunday
        self.day_names = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
        self.shifts = [Shift.MORNING, Shift.AFTERNOON, Shift.EVENING]

    def add_employee(self):
        """Interactively adds a new employee with their shift preferences"""
        # Get and validate employee name
        name = input("\n👤 Enter employee name: ").strip()

        if not name:
            print("❌ Name cannot be empty!")
            return

        # Check for duplicate names
        for emp in self.employees:
            if emp.is_duplicate_name(name):
                print(f"❌ Employee '{name}' already exists. Please use a different name.")
                return

        employee = Employee(name)

        # Collect shift preferences
        print(f"\n📅 Setting up preferences for {employee.name}")
        print("━" * 50)
        print("🌅 Shifts: 0=Morning | ☀️  1=Afternoon | 🌙 2=Evening")
        print("💡 Press Enter to skip a day")

        for day_num, day_name in enumerate(self.day_names):
            while True:
                user_input = input(f"{day_name:<10}: ").strip()

                if user_input == "":
                    break  # No preference for this day

                try:
                    preference = int(user_input)
                    if 0 <= preference < len(self.shifts):
                        shift = self.shifts[preference]
                        employee.set_preference(day_num, shift)
                        shift_icon = shift._get_shift_icon()
                        print(f"   ✅ Set {shift_icon} {shift.value} preference")
                        break
                    else:
                        print("   ❌ Invalid input. Please enter 0, 1, or 2 (or press Enter to skip)")
                except ValueError:
                    print("   ❌ Invalid input. Please enter 0, 1, or 2 (or press Enter to skip)")

        self.employees.append(employee)
        print(f"\n🎉 Employee {name} added successfully!")
        print(f"📊 Total employees: {len(self.employees)}\n")

    def assign_shifts(self):
        """Executes the complete scheduling algorithm"""
        print("\n🔄 Generating schedule...")
        print("━" * 50)

        self._reset_schedules()
        print("✅ Step 1: Reset all schedules")

        self._assign_preferred_shifts()
        print("✅ Step 2: Assigned preferred shifts")

        self._ensure_minimum_staffing()
        print("✅ Step 3: Ensured minimum staffing")

        print("━" * 50)
        print("🎉 Schedule generated successfully!")
        print(f"📋 {len(self.employees)} employees scheduled across {len(self.days)} days")
        print("💡 Use 'View Schedule' to see the complete weekly schedule.\n")

    def _reset_schedules(self):
        """Clears all employee schedules and the main schedule grid"""
        for employee in self.employees:
            employee.reset_schedule()

        self.schedule = defaultdict(lambda: defaultdict(list))

    def _assign_preferred_shifts(self):
        """Assigns employees to their preferred shifts with conflict resolution"""
        for employee in self.employees:
            for day in self.days:
                if employee.days_worked >= MAX_WORK_DAYS_PER_WEEK:
                    break

                preferred_shift = employee.get_preference(day)
                if preferred_shift is not None:
                    if self._can_assign(employee, day, preferred_shift):
                        self._assign(employee, day, preferred_shift)
                    else:
                        self._resolve_conflict(employee, day, preferred_shift)

    def _resolve_conflict(self, employee: Employee, preferred_day: int, preferred_shift: Shift):
        """Attempts to find alternative assignments when preferred shift is unavailable"""
        # Strategy 1: Try alternative shifts on the same day
        for other_shift in self.shifts:
            if other_shift != preferred_shift:
                if self._can_assign(employee, preferred_day, other_shift):
                    self._assign(employee, preferred_day, other_shift)
                    day_name = self.day_names[preferred_day]
                    print(f"🔄 Conflict resolved: {employee.name} → {day_name} {other_shift.value} (preferred shift full)")
                    return

        # Strategy 2: Try other days
        for other_day in self.days:
            if other_day == preferred_day:
                continue

            # Try preferred shift on this alternative day
            if self._can_assign(employee, other_day, preferred_shift):
                self._assign(employee, other_day, preferred_shift)
                day_name = self.day_names[other_day]
                print(f"🔄 Conflict resolved: {employee.name} → {day_name} {preferred_shift.value} (moved to different day)")
                return

            # Try any available shift on this alternative day
            for other_shift in self.shifts:
                if self._can_assign(employee, other_day, other_shift):
                    self._assign(employee, other_day, other_shift)
                    day_name = self.day_names[other_day]
                    print(f"🔄 Conflict resolved: {employee.name} → {day_name} {other_shift.value} (alternative assignment)")
                    return

        # No resolution found
        print(f"⚠️  Warning: Could not assign {employee.name} anywhere (schedule full)")

    def _can_assign(self, employee: Employee, day: int, shift: Shift) -> bool:
        """Checks if an employee can be assigned to a specific day and shift"""
        if not employee.can_work_day(day):
            return False

        if len(self.schedule[day][shift]) >= MAX_EMPLOYEES_PER_SHIFT:
            return False

        return True

    def _assign(self, employee: Employee, day: int, shift: Shift):
        """Assigns employee to shift and updates both schedules"""
        if employee.assign_shift(day, shift):
            self.schedule[day][shift].append(employee.name)

    def _ensure_minimum_staffing(self):
        """Fills understaffed shifts to meet minimum requirements"""
        for day in self.days:
            for shift in self.shifts:
                current_staff = len(self.schedule[day][shift])

                if current_staff < MIN_EMPLOYEES_PER_SHIFT:
                    needed = MIN_EMPLOYEES_PER_SHIFT - current_staff
                    available_employees = self._get_available_employees(day)

                    # Randomly shuffle for fair distribution
                    random.shuffle(available_employees)

                    assigned = 0
                    for employee in available_employees:
                        if assigned >= needed:
                            break

                        if self._can_assign(employee, day, shift):
                            self._assign(employee, day, shift)
                            assigned += 1

                    if assigned < needed:
                        day_name = self.day_names[day]
                        print(f"⚠️  Understaffed: {day_name} {shift.value} needs {needed-assigned} more employees (have {current_staff+assigned}/{MIN_EMPLOYEES_PER_SHIFT})")

    def _get_available_employees(self, day: int) -> List[Employee]:
        """Returns employees who can work on the specified day"""
        return [emp for emp in self.employees if emp.can_work_day(day)]

    def print_schedule(self):
        """Displays the complete weekly schedule and employee summaries"""
        if not self.employees:
            print("\n📝 No employees added yet! Please add some employees first.")
            return

        self._print_schedule_header()
        self._print_weekly_grid()
        self._print_schedule_footer()
        self._print_employee_summaries()

    def _print_schedule_header(self):
        """Displays the main schedule title"""
        print("\n" + "═" * 90)
        print("📅                         WEEKLY EMPLOYEE SCHEDULE                          📅")
        print("═" * 90)

    def _print_weekly_grid(self):
        """Displays the main schedule in a clean tabular format"""
        for day_num, day_name in enumerate(self.day_names):
            print(f"\n📅 {day_name.upper()}")
            print("─" * 85)

            for shift in self.shifts:
                employees = self.schedule[day_num][shift]
                staff_count = len(employees)

                shift_icon = self._get_shift_icon(shift)
                print(f"   {shift_icon} {shift.value:<10} | ", end="")

                if staff_count == 0:
                    print(f"{'No employees assigned':<50}", end="")
                    if staff_count < MIN_EMPLOYEES_PER_SHIFT:
                        print(f" 🚨 UNDERSTAFFED (need {MIN_EMPLOYEES_PER_SHIFT})", end="")
                else:
                    employee_list = ", ".join(employees)
                    if len(employee_list) > 45:
                        employee_list = employee_list[:42] + "..."
                    print(f"{employee_list:<50}", end="")

                    # Status indicator
                    if staff_count < MIN_EMPLOYEES_PER_SHIFT:
                        print(f" 🚨 UNDERSTAFFED ({staff_count}/{MIN_EMPLOYEES_PER_SHIFT})", end="")
                    elif staff_count < MAX_EMPLOYEES_PER_SHIFT:
                        print(f" ✅ STAFFED ({staff_count}/{MAX_EMPLOYEES_PER_SHIFT})", end="")
                    else:
                        print(f" 🏆 FULL ({staff_count}/{MAX_EMPLOYEES_PER_SHIFT})", end="")
                print()

    def _print_schedule_footer(self):
        """Displays summary statistics"""
        print("\n" + "═" * 90)

        total_shifts = len(self.days) * len(self.shifts)
        staffed_shifts = 0
        full_shifts = 0
        total_assignments = 0

        for day in self.days:
            for shift in self.shifts:
                count = len(self.schedule[day][shift])
                total_assignments += count
                if count >= MIN_EMPLOYEES_PER_SHIFT:
                    staffed_shifts += 1
                if count == MAX_EMPLOYEES_PER_SHIFT:
                    full_shifts += 1

        print(f"📊 SCHEDULE STATS: {staffed_shifts}/{total_shifts} shifts properly staffed │ {full_shifts} full shifts │ {total_assignments} total assignments")
        print("═" * 90)

    def _print_employee_summaries(self):
        """Displays individual employee work schedules"""
        print("\n👥 EMPLOYEE WORK SUMMARIES")
        print("─" * 60)

        for i, employee in enumerate(self.employees):
            print(employee.get_work_summary(), end="")
            if i < len(self.employees) - 1:
                print()

        print("─" * 60)

    def _get_shift_icon(self, shift: Shift) -> str:
        """Returns an emoji icon for each shift type"""
        if shift == Shift.MORNING:
            return "🌅"
        elif shift == Shift.AFTERNOON:
            return "☀️"
        elif shift == Shift.EVENING:
            return "🌙"
        else:
            return "⏰"

    def get_employee_count(self) -> int:
        """Returns the total number of employees"""
        return len(self.employees)