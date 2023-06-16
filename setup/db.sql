CREATE DATABASE aukdcdom;

\c aukdcdom;


CREATE TABLE Users
(
    "ID" int NOT NULL,
    "Name" character varying NOT NULL,
    "PhoneNumber" bigint NOT NULL UNIQUE CHECK ("PhoneNumber" BETWEEN 6000000000 AND 9999999999),
    "Email" character varying NOT NULL UNIQUE,
    "HashedPassword" character varying NOT NULL,
    "RoleID" int NOT NULL,
     PRIMARY KEY ("ID")
);
CREATE TABLE Role
(
    "RoleID" int NOT NULL CHECK ("RoleID" BETWEEN 1 and 3),
    "Role" character varying NOT NULL CHECK ("Role" in ('Admin','Faculty','Both')),
     PRIMARY KEY ("RoleID")
);
CREATE TABLE Faculty
(
    "FacultyID" int NOT NULL REFERENCES Users("ID") ON DELETE CASCADE,
    "FacultyType" character varying NOT NULL CHECK ("FacultyType" in ('Permanent','Contract/Guest','Visiting')),
    "Department" character varying NOT NULL,
    "Designation" character varying NOT NULL CHECK ("Designation" in ('Professor','Assistant Professor', 'Associate Professor','Teaching Fellow', 'Emeritus Professor', 'Assistant Professor (SG)')),
    "PanID" character varying(10) NOT NULL UNIQUE,
    "PanPicture" character varying NOT NULL,
    "ExtensionNumber" bigint NOT NULL CHECK ("ExtensionNumber" BETWEEN 20000000 AND 99999999) UNIQUE,
    "Esign" character varying NOT NULL,
    "TDS" real NOT NULL,
     PRIMARY KEY ("FacultyID")
);

CREATE TABLE Account
(
    "BankName" character varying NOT NULL CHECK ("BankName" in ('Canara Bank', 'Indian Bank','State Bank of India')),
    "FacultyID" int NOT NULL REFERENCES Faculty("FacultyID") ON DELETE CASCADE,
    "AccountNumber" bigint NOT NULL,
    "IFSCCode" character varying NOT NULL,
    "Passbook" character varying NOT NULL,
     PRIMARY KEY ("FacultyID")
);

CREATE TABLE Course
(
    "CourseCode" character varying NOT NULL,
    "Title" character varying NOT NULL,
    "Regulation" character varying NOT NULL,
     PRIMARY KEY ("CourseCode")
);

CREATE TABLE Department
(
    "Degree" character varying NOT NULL,
    "CourseCode" character varying REFERENCES Course("CourseCode") ON DELETE CASCADE,
    "Branch" character varying NOT NULL,
    "DegreeType" character varying NOT NULL CHECK ("DegreeType" in ('Regular','Part-time')),
    "Department" character varying NOT NULL,
     PRIMARY KEY ("DegreeType", "Branch")
);

CREATE TABLE HonorariumType
(
    "TypeID" integer NOT NULL CHECK ("TypeID" BETWEEN 1 and 2),
    "Type" character varying NOT NULL CHECK ("Type" in ('Paper Valuation','Question Paper/Key')),
     PRIMARY KEY ("TypeID")
);

CREATE TABLE Honorarium
(
    "TransactionID" character varying NOT NULL,
    "FacultyID" int NOT NULL REFERENCES Faculty("FacultyID") ON DELETE CASCADE,
    "CourseCode" character varying NOT NULL REFERENCES Course("CourseCode") ON DELETE CASCADE,
    "InitialAmount" integer NOT NULL,
    "FinalAmount" integer NOT NULL,
    "TypeID" integer NOT NULL REFERENCES HonorariumType("TypeID") ON DELETE CASCADE,
    "CreatedTime" TIMESTAMP NOT NULL,
     PRIMARY KEY ("TransactionID")
);

CREATE TABLE "Paper Valuation"
(
    "TransactionID" character varying NOT NULL REFERENCES Honorarium("TransactionID") ON DELETE CASCADE,
    "TypeID" integer NOT NULL REFERENCES HonorariumType("TypeID") ON DELETE CASCADE,
    "AnswerScriptRate" integer NOT NULL,
    "AnswerScriptCount" integer NOT NULL,
     PRIMARY KEY ("TransactionID","TypeID")
);

CREATE TABLE "Question Paper/Key"
(
    "TransactionID" character varying NOT NULL REFERENCES Honorarium("TransactionID") ON DELETE CASCADE,
    "TypeID" integer NOT NULL REFERENCES HonorariumType("TypeID") ON DELETE CASCADE,
    "QuestionPaperCount" integer NOT NULL,
    "KeyCount" integer NOT NULL,
    "KeyRate" integer NOT NULL,
    "QuestionPaperRate" integer NOT NULL,
     PRIMARY KEY ("TransactionID","TypeID")
);

CREATE TABLE TimeTable
(
    "Date" date NOT NULL,
    "CourseCode" character varying NOT NULL REFERENCES Course("CourseCode") ON DELETE CASCADE,
    "PaperSetter" int NOT NULL REFERENCES Faculty("FacultyID") ON DELETE CASCADE,
    "Invigilator" character varying NOT NULL,
    "AnnualSession" character varying NOT NULL,
    "ExamType" character varying NOT NULL CHECK ("ExamType" in ('Regular','Re-Appear(RA)')),
    "DailySession" character varying NOT NULL CHECK ("DailySession" in ('FN','AN')),
     PRIMARY KEY ("Date", "CourseCode")
);

CREATE TABLE Admin
(
    "ID" int NOT NULL,
    "Name" character varying NOT NULL,
    "Password" character varying NOT NULL,
    "PhoneNumber" bigint NOT NULL,
    "Type" character varying NOT NULL CHECK ("Type" in ('Faculty','Superintendent')),
    "Email" character varying NOT NULL,
    "Session" character varying NOT NULL ,
     PRIMARY KEY ("ID")
);

CREATE TABLE sessions (
     token CHAR(43) PRIMARY KEY,
     data BYTEA NOT NULL,
     expiry TIMESTAMPTZ NOT NULL
);

INSERT INTO HonorariumType("TypeID","Type") VALUES(2,'Paper Valuation'), (1, 'Question Paper/Key');
INSERT INTO Role("RoleID","Role") VALUES(1,'Admin'), (2, 'Faculty'), (3, 'Both');
INSERT INTO Users("ID","Name","PhoneNumber","Email","HashedPassword","RoleID") VALUES(12345,'test',9876543210,'fac@gmail.com','$2a$12$qIIvAsFFmf979hkMXZhsbuTAhBGmr8oQFbqXY4fO/bCYTXItyaD92',1);
CREATE USER webaukdcdom;
ALTER USER webaukdcdom WITH PASSWORD 'neodom';
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Account TO  webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.sessions TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Faculty TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Course TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public."Paper Valuation" TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public."Question Paper/Key" TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.TimeTable TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Admin TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Honorarium TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Department TO webaukdcdom;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.Users TO webaukdcdom;
