CREATE EXTENSION "uuid-ossp";

CREATE TYPE "genders" AS ENUM (
  'ชาย',
  'หญิง',
  'อื่นๆ'
);

CREATE TYPE "blood_types" AS ENUM (
  'A',
  'B',
  'AB',
  'O'
);

CREATE TYPE "honor_types" AS ENUM (
  'ไม่มี',
  'เกียรตินิยมอันดับ 1',
  'เกียรตินิยมอันดับ 2'
);



CREATE TABLE "prefixes" (
  "id" SERIAL PRIMARY KEY,
  "name_th" VARCHAR(50) UNIQUE NOT NULL,
  "name_en" VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE "roles" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "is_government" boolean DEFAULT 'false',
  "description" varchar
);

CREATE TABLE "appointment_types" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "person_data" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "prefix_id" int,
  "role_id" uuid,
  "first_name_th" varchar,
  "middle_name_th" varchar,
  "last_name_th" varchar,
  "first_name_en" varchar,
  "middle_name_en" varchar,
  "last_name_en" varchar,
  "national_id" varchar(13),
  "passport_id" varchar(13),
  "gender_id" genders,
  "religion_id" varchar,
  "nationality" varchar,
  "ethnicity" varchar,
  "blood_type" blood_types,
  "physical_status" varchar,
  "phone" varchar,
  "personal_email" varchar,
  "birth_date" date,
  "appointment_type_id" int,
  "marital_status" VARCHAR(30),
  "spouse_name" varchar(50),
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "person_addresses" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "house_registry_no" varchar,
  "house_no" varchar,
  "moo" varchar,
  "building" varchar,
  "soi" varchar,
  "road" varchar,
  "subdistrict_code" varchar,
  "subdistrict" varchar,
  "district_code" varchar,
  "district" varchar,
  "province_code" varchar,
  "province" varchar,
  "postal_code" varchar,
  "is_house_registered" boolean DEFAULT false,
  "is_current" boolean DEFAULT false,
  "is_retired_residence" boolean DEFAULT false,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "relation_types" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "family_members" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "relation_type_id" int NOT NULL,
  "prefix_id" int,
  "first_name_th" varchar NOT NULL,
  "middle_name_th" varchar,
  "last_name_th" varchar NOT NULL,
  "former_last_name_th" varchar,
  "national_id" varchar(13),
  "gender" genders,
  "is_alive" boolean,
  "birth_date" date,
  "phone" varchar,
  "email" varchar,
  "marital_status" varchar,
  "is_emergency_contact" boolean DEFAULT false,
  "house_no" varchar,
  "moo" varchar,
  "building" varchar,
  "soi" varchar,
  "road" varchar,
  "subdistrict_code" varchar,
  "subdistrict" varchar,
  "district_code" varchar,
  "district" varchar,
  "province_code" varchar,
  "province" varchar,
  "postal_code" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "education_levels" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "qualifications" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "initial_name" varchar NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "subjects_groups" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "education_records" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "education_level_id" int NOT NULL,
  "qualification_id" int NOT NULL,
  "major_group_id" int,
  "minor_id" int,
  "institution" varchar NOT NULL,
  "country" varchar NOT NULL,
  "start_date" date,
  "graduation_date" date,
  "honor_type" honor_types,
  "is_used_for_first_appointment" boolean DEFAULT false,
  "is_highest_degree" boolean DEFAULT false,
  "is_meets_position_standard" boolean DEFAULT false,
  "is_recognized_by_kksa" boolean DEFAULT false,
  "is_recognized_by_ocsc" boolean DEFAULT false,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "license_types" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "professional_licenses" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "license_type_id" int NOT NULL,
  "license_no" varchar NOT NULL,
  "issue_date" date NOT NULL,
  "expiry_date" date NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "executive_group" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "delete_at" TIMESTAMP
);

CREATE TABLE "departments" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "delete_at" TIMESTAMP
);

CREATE TABLE "academic_rank_type" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "delete_at" TIMESTAMP
);

CREATE TABLE "position_ranks" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "academic_rank" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "government_id" uuid NOT NULL,
  "ranking_type_id" int NOT NULL,
  "position_rank_id" int NOT NULL,
  "criteria" varchar,
  "field_of_expertise" varchar,
  "special_reward" decimal(12,2),
  "position_reward" decimal(12,2),
  "awarded_date" date,
  "approval_date" date,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "delete_at" TIMESTAMP
);

CREATE TABLE "salary_history" (
  "id" SERIAL PRIMARY KEY,
  "person_id" uuid NOT NULL,
  "round" varchar,
  "order_no" varchar NOT NULL,
  "percent_increase" decimal,
  "salary_after" decimal(12,2),
  "evaluation_level" varchar,
  "create_at" TIMESTAMP NOT NULL,
  "delete_at" TIMESTAMP
);

CREATE TABLE "work_statuses" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(50) UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "delete_at" TIMESTAMP
);

CREATE TABLE "governments" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "personnel_id" uuid NOT NULL,
  "work_status_id" int NOT NULL,
  "executive_group_id" int,
  "department_id" int,
  "subject_group_id" int,
  "salary" decimal(12,2),
  "retirement_date" date,
  "government_entry_date" date,
  "direct_pay_no" varchar,
  "position_number" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "decorations" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "government_id" uuid NOT NULL,
  "announce_date" date,
  "year" int,
  "class" varchar,
  "gazette_book_no" varchar,
  "gazette_section" varchar,
  "gazette_publish_date" date,
  "gazette_page" varchar,
  "gazette_order_no" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "trainings" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "academic_year" varchar,
  "course_title" varchar NOT NULL,
  "start_date" date,
  "end_date" date,
  "hours" int,
  "organizer" varchar,
  "location" varchar,
  "outcomes" text,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "innovations" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "academic_year" varchar,
  "project_title" varchar NOT NULL,
  "start_date" date,
  "end_date" date,
  "organizer" varchar,
  "location" varchar,
  "supervisor" varchar,
  "subject_group" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "innovation_student_awards" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "innovation_id" uuid NOT NULL,
  "student_name" varchar NOT NULL,
  "class_level" varchar NOT NULL,
  "class_room" varchar,
  "award" varchar,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "contract_actions" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP
);

CREATE TABLE "hiring_natures" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "budget_types" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "school_revenue_types" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "has_social_security" boolean NOT NULL DEFAULT false,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "employment_position" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "social_security_scheme" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "employment_contracts" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "contract_action_id" int NOT NULL,
  "hiring_nature_id" int NOT NULL,
  "employment_position_id" int NOT NULL,
  "department_id" int,
  "salary" decimal(12,2) NOT NULL,
  "social_security_id" int NOT NULL,
  "period_start" date NOT NULL,
  "period_end" date NOT NULL,
  "order_no" varchar,
  "contract_sequence" int,
  "budget_type_id" int NOT NULL,
  "school_revenue_type_id" int,
  "first_contract_start" date,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "leave_records" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "person_id" uuid NOT NULL,
  "leave_type" varchar(50) NOT NULL,
  "reason" text,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "total_days" decimal(5,2),
  "period_kind" varchar(50),
  "attendance_issue" varchar(100),
  "created_at" timestamp NOT NULL,
  "deleted_at" TIMESTAMP
);

CREATE TABLE "document_links" (
  "document_key" varchar PRIMARY KEY,
  "entity_table" varchar(64) NOT NULL,
  "entity_id" uuid NOT NULL,
  "purpose" varchar(50),
  "created_by" uuid,
  "created_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

COMMENT ON COLUMN "prefixes"."name_th" IS 'คำนำหน้า เช่น นาย , นาง , นางสาว';

COMMENT ON COLUMN "prefixes"."name_en" IS 'คำนำหน้า ENG เช่น MR , Miss , Mrs';

COMMENT ON COLUMN "roles"."name" IS 'ครู/ครูผู้บริหาร,ครูผู้ช่วย,แม่บ้าน';

COMMENT ON COLUMN "roles"."description" IS 'อธิบาย permission ไม่ต้องใส่ก็ได้';

COMMENT ON COLUMN "appointment_types"."name" IS 'ประเภทการบรรจุ';

COMMENT ON COLUMN "person_data"."national_id" IS 'เลขบัตรประชาชน';

COMMENT ON COLUMN "person_data"."passport_id" IS 'เลขpassport';

COMMENT ON COLUMN "person_data"."nationality" IS 'สัญชาติ';

COMMENT ON COLUMN "person_data"."ethnicity" IS 'เชื้อชาติ';

COMMENT ON COLUMN "person_data"."blood_type" IS 'group เลือด';

COMMENT ON COLUMN "person_data"."physical_status" IS 'สถานภาพทางร่างกาย';

COMMENT ON COLUMN "person_data"."personal_email" IS 'อีเมลส่วนตัว';

COMMENT ON COLUMN "person_data"."appointment_type_id" IS 'ประเภทการบรรจุ';

COMMENT ON COLUMN "person_data"."marital_status" IS 'สถานะภาพ เช่น สมรส โสด หย่า ร้าง และ หม้าย';

COMMENT ON COLUMN "person_data"."spouse_name" IS 'ชื่อคู่สมรส';

COMMENT ON COLUMN "person_addresses"."house_registry_no" IS 'รหัสทะเบียนบ้าน';

COMMENT ON COLUMN "person_addresses"."house_no" IS 'บ้านเลขที่';

COMMENT ON COLUMN "person_addresses"."moo" IS 'หมู่ที่';

COMMENT ON COLUMN "person_addresses"."building" IS 'หมู่บ้าน/อาคาร';

COMMENT ON COLUMN "person_addresses"."soi" IS 'ตรอก/ซอย';

COMMENT ON COLUMN "person_addresses"."road" IS 'ถนน';

COMMENT ON COLUMN "person_addresses"."subdistrict" IS 'ตำบล/แขวง';

COMMENT ON COLUMN "person_addresses"."district" IS 'อำเภอ/เขต';

COMMENT ON COLUMN "person_addresses"."province" IS 'จังหวัด';

COMMENT ON COLUMN "person_addresses"."postal_code" IS 'รหัสไปรษณีย์';

COMMENT ON COLUMN "person_addresses"."is_house_registered" IS 'ที่อยู่ตามทะเบียนบ้าน';

COMMENT ON COLUMN "person_addresses"."is_current" IS 'ที่อยู่ปัจจุบัน';

COMMENT ON COLUMN "person_addresses"."is_retired_residence" IS 'ที่อยู่หลังเกษียณ';

COMMENT ON COLUMN "relation_types"."name" IS 'บิดา/มารดา/คู่สมรส/บุตร/ปู่/ย่า/ตา/ยาย';

COMMENT ON COLUMN "family_members"."former_last_name_th" IS 'นามสกุลเดิม';

COMMENT ON COLUMN "family_members"."marital_status" IS 'สถานะภาพ เช่น สมรส โสด หย่า ร้าง และ หม้าย';

COMMENT ON COLUMN "family_members"."is_emergency_contact" IS 'เลือกได้ 1 คน/ผู้ใช้';

COMMENT ON COLUMN "family_members"."house_no" IS 'บ้านเลขที่';

COMMENT ON COLUMN "family_members"."moo" IS 'หมู่ที่';

COMMENT ON COLUMN "family_members"."building" IS 'หมู่บ้าน/อาคาร';

COMMENT ON COLUMN "family_members"."soi" IS 'ตรอก/ซอย';

COMMENT ON COLUMN "family_members"."road" IS 'ถนน';

COMMENT ON COLUMN "family_members"."subdistrict" IS 'ตำบล/แขวง';

COMMENT ON COLUMN "family_members"."district" IS 'อำเภอ/เขต';

COMMENT ON COLUMN "family_members"."province" IS 'จังหวัด';

COMMENT ON COLUMN "family_members"."postal_code" IS 'รหัสไปรษณีย์';

COMMENT ON COLUMN "education_levels"."name" IS 'ปริญญาตรี/โท/เอก/เทียบเท่า';

COMMENT ON COLUMN "qualifications"."name" IS 'วุฒิการศึกษา เช่น วิทยาศาสตร์บัณฑิต';

COMMENT ON COLUMN "qualifications"."initial_name" IS 'ชื่อย่อ';

COMMENT ON COLUMN "subjects_groups"."name" IS 'ชื่อสาขาวิชาเอก';

COMMENT ON COLUMN "education_records"."major_group_id" IS 'สาขาวิชาเอก';

COMMENT ON COLUMN "education_records"."minor_id" IS 'สาขาวิชาโท';

COMMENT ON COLUMN "education_records"."institution" IS 'สถานศึกษา';

COMMENT ON COLUMN "education_records"."country" IS 'ประเทศ';

COMMENT ON COLUMN "education_records"."start_date" IS 'วันที่เข้าศึกษา';

COMMENT ON COLUMN "education_records"."graduation_date" IS 'วันทีสำเร็จการศึกษา';

COMMENT ON COLUMN "education_records"."honor_type" IS 'เกียรตินิยม';

COMMENT ON COLUMN "education_records"."is_used_for_first_appointment" IS 'ใช้บรรจุครั้งแรก';

COMMENT ON COLUMN "education_records"."is_highest_degree" IS 'วุฒิการศึกษาสูงสุด';

COMMENT ON COLUMN "education_records"."is_meets_position_standard" IS 'วุฒิตรงตามมาตรฐานตำแหน่ง';

COMMENT ON COLUMN "education_records"."is_recognized_by_kksa" IS 'วุฒิ กคศ./คุรุสภา รับรอง';

COMMENT ON COLUMN "education_records"."is_recognized_by_ocsc" IS 'วุฒิ ก.พ. รับรอง';

COMMENT ON COLUMN "license_types"."name" IS 'P-License/B-License/A-License/ใบอนุญาตสอน';

COMMENT ON COLUMN "professional_licenses"."license_no" IS 'เลขที่ใบประกอบวิชาชีพ';

COMMENT ON COLUMN "professional_licenses"."issue_date" IS 'วันออกบัตร';

COMMENT ON COLUMN "professional_licenses"."expiry_date" IS 'วันหมดอายุ';

COMMENT ON COLUMN "academic_rank_type"."id" IS 'ประเภทวิทยฐานะ เช่น ชำนาญการ / ชำนาญการพิเศษ / เชี่ยวชาญ / เชี่ยวชาญพิเศษ';

COMMENT ON COLUMN "position_ranks"."name" IS 'ครูผู้ช่วย / คศ.1 / คศ.2 / คศ.3 / คศ.4 / คศ.5';

COMMENT ON COLUMN "academic_rank"."id" IS 'วิทยฐานะ เก็บรวมถึงประวัติด้วย';

COMMENT ON COLUMN "academic_rank"."ranking_type_id" IS 'ชนิดวิทยฐานะ';

COMMENT ON COLUMN "academic_rank"."criteria" IS 'หลักเกณฑ์การได้รับวิทยฐานะ (เกณฑ์ ว13 / ว17 / ว21 เป็นต้น)';

COMMENT ON COLUMN "academic_rank"."field_of_expertise" IS 'สาขาที่ได้รับวิทยฐานะ เช่น วิทยาศาสตร์, คณิตศาสตร์';

COMMENT ON COLUMN "academic_rank"."special_reward" IS 'เงินค่าตอบแทนพิเศษกรณีเต็มขั้น';

COMMENT ON COLUMN "academic_rank"."position_reward" IS 'เงินประจำตำแหน่ง';

COMMENT ON COLUMN "academic_rank"."awarded_date" IS 'วันที่ได้รับวิทยฐานะ';

COMMENT ON COLUMN "academic_rank"."approval_date" IS 'วันที่อนุมัติ/ประกาศในคำสั่ง';

COMMENT ON COLUMN "salary_history"."id" IS 'history จะไม่มีการupdate หรือ delete จะสร้างใหม่อย่างเดียว';

COMMENT ON COLUMN "salary_history"."round" IS 'รอบที่ขึ้นเงินเดือน เช่น รอบ 1';

COMMENT ON COLUMN "salary_history"."order_no" IS 'เลขที่คำสั่งเลื่อนเงินเดือน';

COMMENT ON COLUMN "salary_history"."percent_increase" IS 'ร้อยละการเลื่อน';

COMMENT ON COLUMN "salary_history"."salary_after" IS 'เงินเดือนหลังเลื่อน';

COMMENT ON COLUMN "salary_history"."evaluation_level" IS 'ระดับผลการประเมิน';

COMMENT ON COLUMN "work_statuses"."name" IS 'ทำงาน/ย้าย/เกษียณอายุราชการ/ช่วยราชการ/ลาออก/เสียชีวิต';

COMMENT ON COLUMN "governments"."subject_group_id" IS 'ใช้เป็นกลุ่มสาระถ้าเก็บในตารางนี้';

COMMENT ON COLUMN "governments"."salary" IS 'เงินเดือน';

COMMENT ON COLUMN "governments"."retirement_date" IS 'วันเกษียณ';

COMMENT ON COLUMN "governments"."government_entry_date" IS 'วันเดือนปีเข้ารับราชการ';

COMMENT ON COLUMN "governments"."direct_pay_no" IS 'เลขที่จ่ายตรง';

COMMENT ON COLUMN "governments"."position_number" IS 'เลขที่ตำแหน่ง (เลขประจำตำแหน่ง)';

COMMENT ON COLUMN "decorations"."id" IS 'id การรับเครื่องราชอิสริยาภรณ์';

COMMENT ON COLUMN "decorations"."announce_date" IS 'วันที่ประกาศในราชกิจจา';

COMMENT ON COLUMN "decorations"."year" IS 'ประจำปี พ.ศ.';

COMMENT ON COLUMN "decorations"."class" IS 'ชั้นตรา';

COMMENT ON COLUMN "decorations"."gazette_book_no" IS 'เล่มที่';

COMMENT ON COLUMN "decorations"."gazette_section" IS 'ตอนที่';

COMMENT ON COLUMN "decorations"."gazette_publish_date" IS 'ลงวันที่';

COMMENT ON COLUMN "decorations"."gazette_page" IS 'หน้า';

COMMENT ON COLUMN "decorations"."gazette_order_no" IS 'ลำดับที่ / เลขคำสั่ง (ถ้ามี)';

COMMENT ON COLUMN "trainings"."id" IS 'อบรบดูงาน';

COMMENT ON COLUMN "trainings"."academic_year" IS 'ปีการศึกษา เช่น 2568';

COMMENT ON COLUMN "trainings"."course_title" IS 'หลักสูตร/เรื่อง/หัวข้อ';

COMMENT ON COLUMN "trainings"."start_date" IS 'ระหว่างวันที่ (เริ่ม)';

COMMENT ON COLUMN "trainings"."end_date" IS 'ระหว่างวันที่ (สิ้นสุด)';

COMMENT ON COLUMN "trainings"."hours" IS 'จำนวนชั่วโมง';

COMMENT ON COLUMN "trainings"."organizer" IS 'หน่วยงานที่จัด';

COMMENT ON COLUMN "trainings"."location" IS 'สถานที่จัด';

COMMENT ON COLUMN "trainings"."outcomes" IS 'สิ่งที่ได้รับจากการอบรม';

COMMENT ON COLUMN "innovations"."academic_year" IS 'ปีการศึกษา';

COMMENT ON COLUMN "innovations"."project_title" IS 'หลักสูตร/เรื่อง/หัวข้อ';

COMMENT ON COLUMN "innovations"."start_date" IS 'ระหว่างวันที่ (เริ่ม)';

COMMENT ON COLUMN "innovations"."end_date" IS 'ระหว่างวันที่ (สิ้นสุด)';

COMMENT ON COLUMN "innovations"."organizer" IS 'หน่วยงานที่จัด';

COMMENT ON COLUMN "innovations"."location" IS 'สถานที่จัด';

COMMENT ON COLUMN "innovations"."supervisor" IS 'ครูผู้ดูแล';

COMMENT ON COLUMN "innovations"."subject_group" IS 'กลุ่มสาระการเรียนรู้';

COMMENT ON COLUMN "innovation_student_awards"."class_level" IS 'ระดับชั้น เช่น ม.1, ม.2';

COMMENT ON COLUMN "innovation_student_awards"."class_room" IS 'ห้อง เช่น 1, 2';

COMMENT ON COLUMN "innovation_student_awards"."award" IS 'ชื่อรางวัล/ระดับรางวัล';

COMMENT ON COLUMN "contract_actions"."name" IS 'ทำสัญญาครั้งแรก, ต่อสัญญา';

COMMENT ON COLUMN "hiring_natures"."name" IS 'จ้างใหม่, ต่อเนื่อง';

COMMENT ON COLUMN "budget_types"."name" IS 'เงินงบประมาณ,เงินนอกงบประมาณ,เงินรายได้สถานศึกษา,เงินอุดหนุนค่าใช้จ่ายรายหัว';

COMMENT ON COLUMN "school_revenue_types"."name" IS 'ลูกจ้างชั่วคราว,เหมาบริการ';

COMMENT ON COLUMN "school_revenue_types"."has_social_security" IS 'มีประกันสังคมให้ไหม';

COMMENT ON COLUMN "social_security_scheme"."name" IS 'มาตรา 33,39,40,อื่นๆ';

COMMENT ON COLUMN "employment_contracts"."salary" IS 'เงินค่าจ้าง/เดือน';

COMMENT ON COLUMN "leave_records"."leave_type" IS 'ลาป่วย/ลากิจ/ลาคลอด/ลาช่วยภริยา/ลาบวช';

COMMENT ON COLUMN "leave_records"."reason" IS 'เหตุผลการลา';

COMMENT ON COLUMN "leave_records"."start_date" IS 'ลาตั้งแต่วันที่';

COMMENT ON COLUMN "leave_records"."end_date" IS 'ลาสิ้นสุดวันที่';

COMMENT ON COLUMN "leave_records"."total_days" IS 'จำนวนวันลา เช่น 0.5, 1, 3';

COMMENT ON COLUMN "leave_records"."period_kind" IS 'ลาครึ่งวันเช้า,บ่าย,เต็มวัน';

COMMENT ON COLUMN "leave_records"."attendance_issue" IS 'ไม่ได้สแกนหน้า/ลืมสแกน/เครื่องไม่บันทึก/ไปราชการ/สาย';

COMMENT ON COLUMN "document_links"."document_key" IS 'id เชื่อมกับ resouce.id';

COMMENT ON COLUMN "document_links"."entity_table" IS 'เช่น ''decorations'',''trainings'',''employment_contracts'', ...';

COMMENT ON COLUMN "document_links"."entity_id" IS 'PK ของแถวนั้น ๆ';

COMMENT ON COLUMN "document_links"."purpose" IS 'เช่น ''order_scan'',''certificate'',''photo'',''attachment''';

ALTER TABLE "person_data" ADD FOREIGN KEY ("prefix_id") REFERENCES "prefixes" ("id");

ALTER TABLE "person_data" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "person_data" ADD FOREIGN KEY ("appointment_type_id") REFERENCES "appointment_types" ("id");

ALTER TABLE "person_addresses" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "family_members" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "family_members" ADD FOREIGN KEY ("relation_type_id") REFERENCES "relation_types" ("id");

ALTER TABLE "family_members" ADD FOREIGN KEY ("prefix_id") REFERENCES "prefixes" ("id");

ALTER TABLE "education_records" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "education_records" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");

ALTER TABLE "education_records" ADD FOREIGN KEY ("qualification_id") REFERENCES "qualifications" ("id");

ALTER TABLE "education_records" ADD FOREIGN KEY ("major_group_id") REFERENCES "subjects_groups" ("id");

ALTER TABLE "education_records" ADD FOREIGN KEY ("minor_id") REFERENCES "subjects_groups" ("id");

ALTER TABLE "professional_licenses" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "professional_licenses" ADD FOREIGN KEY ("license_type_id") REFERENCES "license_types" ("id");

ALTER TABLE "academic_rank" ADD FOREIGN KEY ("government_id") REFERENCES "governments" ("id");

ALTER TABLE "academic_rank" ADD FOREIGN KEY ("ranking_type_id") REFERENCES "academic_rank_type" ("id");

ALTER TABLE "academic_rank" ADD FOREIGN KEY ("position_rank_id") REFERENCES "position_ranks" ("id");

ALTER TABLE "salary_history" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "governments" ADD FOREIGN KEY ("personnel_id") REFERENCES "person_data" ("id");

ALTER TABLE "governments" ADD FOREIGN KEY ("work_status_id") REFERENCES "work_statuses" ("id");

ALTER TABLE "governments" ADD FOREIGN KEY ("executive_group_id") REFERENCES "executive_group" ("id");

ALTER TABLE "governments" ADD FOREIGN KEY ("department_id") REFERENCES "departments" ("id");

ALTER TABLE "governments" ADD FOREIGN KEY ("subject_group_id") REFERENCES "departments" ("id");

ALTER TABLE "decorations" ADD FOREIGN KEY ("government_id") REFERENCES "governments" ("id");

ALTER TABLE "trainings" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "innovations" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "innovation_student_awards" ADD FOREIGN KEY ("innovation_id") REFERENCES "innovations" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("contract_action_id") REFERENCES "contract_actions" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("hiring_nature_id") REFERENCES "hiring_natures" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("employment_position_id") REFERENCES "employment_position" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("department_id") REFERENCES "departments" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("social_security_id") REFERENCES "social_security_scheme" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("budget_type_id") REFERENCES "budget_types" ("id");

ALTER TABLE "employment_contracts" ADD FOREIGN KEY ("school_revenue_type_id") REFERENCES "school_revenue_types" ("id");

ALTER TABLE "leave_records" ADD FOREIGN KEY ("person_id") REFERENCES "person_data" ("id");

ALTER TABLE "document_links" ADD FOREIGN KEY ("created_by") REFERENCES "person_data" ("id");
