-- Drop tables if they exist
IF OBJECT_ID('employees', 'U') IS NOT NULL 
  DROP TABLE employees;
IF OBJECT_ID('projects', 'U') IS NOT NULL 
  DROP TABLE projects;
IF OBJECT_ID('departments', 'U') IS NOT NULL 
  DROP TABLE departments;

-- Create departments table
CREATE TABLE departments (
  id INT PRIMARY KEY,
  name NVARCHAR(100) NOT NULL,
  location NVARCHAR(100)
);

-- Create employees table
CREATE TABLE employees (
  id INT PRIMARY KEY,
  first_name NVARCHAR(50) NOT NULL,
  last_name NVARCHAR(50) NOT NULL,
  email NVARCHAR(100) UNIQUE,
  hire_date DATE,
  department_id INT,
  salary DECIMAL(10,2),
  CONSTRAINT fk_department FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- Create projects table
CREATE TABLE projects (
  id INT PRIMARY KEY,
  name NVARCHAR(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget DECIMAL(12,2)
);

-- Insert data into departments
INSERT INTO departments (id, name, location) VALUES (1, 'IT', 'Building A');
INSERT INTO departments (id, name, location) VALUES (2, 'Finance', 'Building B');
INSERT INTO departments (id, name, location) VALUES (3, 'HR', 'Building A');

-- Insert data into employees
INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (1, 'John', 'Smith', 'john.smith@example.com', '2020-01-15', 1, 85000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (2, 'Jane', 'Doe', 'jane.doe@example.com', '2019-06-01', 2, 75000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (3, 'Robert', 'Johnson', 'robert.johnson@example.com', '2021-03-10', 1, 65000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (4, 'Sarah', 'Williams', 'sarah.williams@example.com', '2018-11-15', 3, 60000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (5, 'Michael', 'Brown', 'michael.brown@example.com', '2020-09-01', 2, 55000);

-- Insert data into projects
INSERT INTO projects (id, name, start_date, end_date, budget) 
VALUES (1, 'ERP Implementation', '2021-01-01', '2021-12-31', 500000);

INSERT INTO projects (id, name, start_date, end_date, budget) 
VALUES (2, 'Website Redesign', '2021-03-15', '2021-08-15', 150000);

INSERT INTO projects (id, name, start_date, end_date, budget) 
VALUES (3, 'Mobile App Development', '2021-06-01', NULL, 300000); 