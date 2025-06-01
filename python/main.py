from scheduler import Scheduler
from employee import MIN_EMPLOYEES_PER_SHIFT, MAX_WORK_DAYS_PER_WEEK

def main():
    scheduler = Scheduler()

    print("Employee Scheduling System")
    print("=" * 25)
    print(f"Constraints: Min {MIN_EMPLOYEES_PER_SHIFT} employees per shift | Max {MAX_WORK_DAYS_PER_WEEK} days per employee\n")

    while True:
        print("1. Add Employee")
        print("2. Generate Schedule")
        print("3. View Schedule")
        print("4. Exit")

        choice = input("Choose an option: ").strip()

        if choice == "1":
            scheduler.add_employee()
        elif choice == "2":
            if scheduler.get_employee_count() == 0:
                print("Please add employees first!")
                print()
                continue
            scheduler.assign_shifts()
        elif choice == "3":
            scheduler.print_schedule()
        elif choice == "4":
            print("Goodbye!")
            break
        else:
            print("Invalid option. Please try again.")

if __name__ == "__main__":
    main()