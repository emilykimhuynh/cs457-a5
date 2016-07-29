/*****************************************************************
*
*   Amusement_Park commands
*
*****************************************************************/
CREATE TABLE Amusement_Park (
    Name VARCHAR (30),
    Id   INT          PRIMARY KEY
);


insert into Amusement_Park values ('Disney World', 10101);
insert into Amusement_Park values ('Universal Studios', 20202);
insert into Amusement_Park values ('Hogwarts', 30303);

/*****************************************************************
*
*   Attraction commands
*
*****************************************************************/
CREATE TABLE Attraction (
    Name    VARCHAR (30) PRIMARY KEY,
    Mgr_Id  INT,
    Park_Id INT          CONSTRAINT fk_Attraction_Park_Id REFERENCES Amusement_Park (Id) ON UPDATE RESTRICT
);


insert into Attraction values ('The Krusty Krab', 12345, 20202);
insert into Attraction values ('Quidditch', 56789, 30303);
insert into Attraction values ('Teacups', 45678, 10101);
insert into Attraction values ('Tower of Terror', 33322, 20202);
insert into Attraction values ('Whomping Willow', 23456, 30303);
insert into Attraction values ('Magic Carpet', 67890, 10101);

/*****************************************************************
*
*   Location commands
*
*****************************************************************/
CREATE TABLE Location (
    Address VARCHAR (255) PRIMARY KEY,
    Park_Id INT           CONSTRAINT fk_Location_Id REFERENCES Amusement_Park (Id) ON UPDATE RESTRICT
);


insert into Location values ('Disney World Resort, Orlando, FL 32830', 10101);
insert into Location values ('6000 Universal Blvd, Orlando, FL 32819', 20202);
insert into Location values ('1 Wizard Ave, London, England', 30303);

/*****************************************************************
*
*   Visitor commands
*
*****************************************************************/
CREATE TABLE Visitor (
    Name      VARCHAR (30),
    Fast_Pass CHAR,
    Ticket_Id INT          PRIMARY KEY
);


insert into Visitor values ('Harry Potter', 'T', 11111);
insert into Visitor values ('Hermione Granger', 'F', 22222);
insert into Visitor values ('Ron Weasley', 'F', 33333);
insert into Visitor values ('Aladdin', 'F', 44444);
insert into Visitor values ('Jasmine', 'T', 55555);
insert into Visitor values ('Batman', 'T', 66666);
insert into Visitor values ('Spongebob Squarepants', 'T', 77777);
insert into Visitor values ('Donald Duck', 'F', 88888);
insert into Visitor values ('Optimus Prime', 'T', 99999);

/*****************************************************************
*
*   Employee commands
*
*****************************************************************/
CREATE TABLE Employee (
    Fname                    VARCHAR (30),
    Minit                    CHAR,
    Lname                    VARCHAR (30),
    Id                       INT          PRIMARY KEY,
    Birth_Date               CHAR (10),
    Salary                   INT,
    Hire_Date                CHAR (10),
    Address                  VARCHAR (30),
    Supervisor_Id            INT          CONSTRAINT fk_Employee_Supervisor_Id REFERENCES Employee (Id) ON UPDATE RESTRICT,
    Attraction_Name          VARCHAR (30) CONSTRAINT fk_Employee_Attraction_Name REFERENCES Attraction (Name),
    Attraction_Date_Assigned CHAR (10)
);


insert into Employee values ('Eugene', 'H', 'Krabs', 12345, '1965-01-02', 20000,
     '2001-10-20', '10 Bikini Bottom', null, 'The Krusty Krab', '2001-10-20');
insert into Employee values ('Minerva', 'M', 'McGonagall', 23456, '1960-04-11', 60000,
    '1995-05-31', '11 Diagon Alley St',	null, 'Whomping Willow', '1995-06-03');
insert into Employee values ('Rubeus', 'K', 'Hagrid', 11122, '1961-03-22', 10000,
    '1994-11-13', '22 Hogwarts Hut', 23456, 'Whomping Willow', '1994-11-13');
insert into Employee values ('Perry', 'W', 'White', 33322, '1970-11-29', 80000,
    '1997-10-02', '12 Metropolis Ave', null, 'Tower of Terror', '1998-06-21');
insert into Employee values ('Clark', 'J', 'Kent', 34567, '1978-12-19', 70000,
    '1998-08-04', '555 Daily Planet Dr', 33322, 'Tower of Terror', '1998-08-04');
insert into Employee values ('Minnie', 'M', 'Mouse', 45678, '1953-02-07', 90000,
    '1980-02-05', '12 Old School Dr', null, 'Teacups',	'1980-02-05');
insert into Employee values ('Mickey', 'S', 'Mouse', 44433, '1952-10-24', 60000,
    '1981-01-14', '12 Old School Dr', 45678, 'Teacups', '1982-12-15');
insert into Employee values ('Rolanda', 'T', 'Hooch', 56789, '1968-07-17', 50000,
    '1983-05-13', '101 Hogsmead Blvd', null, 'Quidditch', '1986-01-12');
insert into Employee values ('Oliver', 'L',	'Wood', 55544, '1990-02-14', 30000,
    '2010-03-10', '101 Magic Rd', 56789, 'Quidditch', '2010-03-10');
insert into Employee values ('Rajah', 'A', 'Tiger', 67890, '1989-07-10', 40000,
    '2009-08-09', '9 Genie Dr', null, 'Magic Carpet', '2009-08-09');

/*****************************************************************
*
*   Goes_To commands
*
*****************************************************************/
CREATE TABLE Goes_To (
    Park_Id           INT,
    Visitor_Ticket_Id INT CONSTRAINT fk_Goes_To_Ticket_Id REFERENCES Visitor (Ticket_Id) ON UPDATE RESTRICT,
    CONSTRAINT fk_Goes_To_Park_Id FOREIGN KEY (
        Park_Id
    )
    REFERENCES Amusement_Park (Id) ON UPDATE RESTRICT
);


insert into Goes_To values (30303, 11111);
insert into Goes_To values (30303, 22222);
insert into Goes_To values (30303, 33333);
insert into Goes_To values (10101, 44444);
insert into Goes_To values (10101, 55555);
insert into Goes_To values (20202, 66666);
insert into Goes_To values (20202, 77777);
insert into Goes_To values (10101, 88888);
insert into Goes_To values (20202, 99999);

/*****************************************************************
*
*   Interacts_With commands
*
*****************************************************************/
CREATE TABLE Interacts_With (
    Visitor_Ticket_Id INT          CONSTRAINT fk_Interacts_With_Ticket_Id REFERENCES Visitor (Ticket_Id) ON UPDATE RESTRICT,
    Attraction_Name   VARCHAR (30) CONSTRAINT fk_Interacts_With_Name REFERENCES Attraction (Name) ON UPDATE RESTRICT
);


insert into Interacts_With values (11111, 'Whomping Willow');
insert into Interacts_With values (11111, 'Quidditch');
insert into Interacts_With values (22222, 'Whomping Willow');
insert into Interacts_With values (33333, 'Quidditch');
insert into Interacts_With values (44444, 'Magic Carpet');
insert into Interacts_With values (55555, 'Magic Carpet');
insert into Interacts_With values (66666, 'Tower of Terror');
insert into Interacts_With values (77777, 'The Krusty Krab');
insert into Interacts_With values (88888, 'Magic Carpet');
insert into Interacts_With values (88888, 'Teacups');
insert into Interacts_With values (99999, 'Tower of Terror');

/*****************************************************************
*
*   Create views commands
*
*****************************************************************/
create view Visitors_At_Parks as
    select v.Name as Visitor_Name, p.Name as Park_Name
    from Visitor v, Amusement_Park p, Goes_To g
    where g.Park_Id = p.Id and g.Visitor_Ticket_Id=v.Ticket_Id;

create view Parks_and_Locations as
    select Name, Address
    from Amusement_Park, Location
    where Id = Park_Id;
