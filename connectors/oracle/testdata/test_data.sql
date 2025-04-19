-- Drop tables if they exist
DROP TABLE employees;
DROP TABLE projects;
DROP TABLE departments;

-- Create departments table
CREATE TABLE departments (
  id NUMBER PRIMARY KEY,
  name VARCHAR2(100) NOT NULL,
  location VARCHAR2(100)
);

-- Create employees table
CREATE TABLE employees (
  id NUMBER PRIMARY KEY,
  first_name VARCHAR2(50) NOT NULL,
  last_name VARCHAR2(50) NOT NULL,
  email VARCHAR2(100) UNIQUE,
  hire_date DATE,
  department_id NUMBER,
  salary NUMBER(10,2),
  CONSTRAINT fk_department FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- Create projects table
CREATE TABLE projects (
  id NUMBER PRIMARY KEY,
  name VARCHAR2(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMBER(12,2)
);

-- Insert data into departments
INSERT INTO departments (id, name, location) VALUES (1, 'IT', 'Building A');
INSERT INTO departments (id, name, location) VALUES (2, 'Finance', 'Building B');
INSERT INTO departments (id, name, location) VALUES (3, 'HR', 'Building A');

-- Insert data into employees
INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (1, 'John', 'Smith', 'john.smith@example.com', TO_DATE('2020-01-15', 'YYYY-MM-DD'), 1, 85000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (2, 'Jane', 'Doe', 'jane.doe@example.com', TO_DATE('2019-06-01', 'YYYY-MM-DD'), 2, 75000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (3, 'Robert', 'Johnson', 'robert.johnson@example.com', TO_DATE('2021-03-10', 'YYYY-MM-DD'), 1, 65000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (4, 'Sarah', 'Williams', 'sarah.williams@example.com', TO_DATE('2018-11-15', 'YYYY-MM-DD'), 3, 60000);

INSERT INTO employees (id, first_name, last_name, email, hire_date, department_id, salary) 
VALUES (5, 'Michael', 'Brown', 'michael.brown@example.com', TO_DATE('2020-09-01', 'YYYY-MM-DD'), 2, 55000);

-- Insert data into projects
INSERT INTO projects (id, name, start_date, end_date, budget) 
VALUES (1, 'ERP Implementation', TO_DATE('2021-01-01', 'YYYY-MM-DD'), TO_DATE('2021-12-31', 'YYYY-MM-DD'), 500000);

INSERT INTO projects (id, name, start_date, end_date, budget) 
VALUES (2, 'Website Redesign', TO_DATE('2021-03-15', 'YYYY-MM-DD'), TO_DATE('2021-08-15', 'YYYY-MM-DD'), 150000);

INSERT INTO projects (id, name, start_date, end_date, budget) 
VALUES (3, 'Mobile App Development', TO_DATE('2021-06-01', 'YYYY-MM-DD'), NULL, 300000);

COMMIT; 