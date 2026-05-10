# สรุป API ระบบจัดการข้อมูล (Data Service)

Base URL: `{base_url}/v1`

เอกสารนี้สรุปรายการ API หลักของระบบจัดการข้อมูล โดยอธิบายหน้าที่ของแต่ละเส้นแบบย่อ และระบุ Parameters ที่ใช้เรียก API

## ระบบตรวจสอบสถานะ (Health Check)

ใช้ตรวจสอบว่าระบบยังทำงานและตอบสนองได้ตามปกติ

**GET `/health-check`**  
ใช้ตรวจสอบสถานะการทำงานของระบบ

## ระบบตรวจสอบโครงสร้างแหล่งข้อมูล (Introspection)

ใช้จัดการและตรวจสอบแหล่งข้อมูล schema ตาราง และคอลัมน์จากฐานข้อมูลจริง

**GET `/introspect/sources`**  
ใช้ดึงรายการแหล่งข้อมูลทั้งหมดที่ลงทะเบียนไว้ในระบบ  
Parameters: `page`, `per_page`

**GET `/introspect/sources/{id}`**  
ใช้ดึงรายละเอียดแหล่งข้อมูลตามรหัสแหล่งข้อมูล  
Parameters: `id`

**POST `/introspect/sources`**  
ใช้สร้างแหล่งข้อมูลใหม่ เช่น ฐานข้อมูล PostgreSQL หรือ MySQL

**PUT `/introspect/sources/{id}`**  
ใช้แก้ไขข้อมูลแหล่งข้อมูลที่มีอยู่  
Parameters: `id`

**DELETE `/introspect/sources/{id}`**  
ใช้ลบแหล่งข้อมูลตามรหัสแหล่งข้อมูล  
Parameters: `id`

**PATCH `/introspect/sources/{id}/activate`**  
ใช้เปิดใช้งานแหล่งข้อมูล  
Parameters: `id`

**PATCH `/introspect/sources/{id}/deactivate`**  
ใช้ปิดใช้งานแหล่งข้อมูล  
Parameters: `id`

**GET `/introspect/schemas`**  
ใช้ดึงรายการ schema ของแหล่งข้อมูล  
Parameters: `source_id`, `page`, `per_page`

**GET `/introspect/tables`**  
ใช้ดึงรายการตารางภายใน schema  
Parameters: `source_id`, `schema`, `page`, `per_page`

**GET `/introspect/columns`**  
ใช้ดึงรายการคอลัมน์ของตาราง  
Parameters: `source_id`, `schema`, `table`, `page`, `per_page`

## ระบบจัดการข้อมูลตารางโดยตรง (Direct Table Data)

ใช้ดึง เพิ่ม แก้ไข และลบข้อมูลจากตารางจริงผ่านแหล่งข้อมูลที่ลงทะเบียนไว้

**GET `/introspect/sources/{id}/schemas/{schema}/tables/{table}/data`**  
ใช้ดึงข้อมูลจากตารางจริงในแหล่งข้อมูลโดยตรง  
Parameters: `id`, `schema`, `table`, `page`, `per_page`, `where`, `where_logical_operator`, `sort_by`, `sort_order`

**GET `/introspect/sources/{id}/schemas/{schema}/tables/{table}/data/key/{key}`**  
ใช้ดึงข้อมูล 1 รายการจากตารางจริงตาม key  
Parameters: `id`, `schema`, `table`, `key`, `key_field`

**POST `/introspect/sources/{id}/schemas/{schema}/tables/{table}/data`**  
ใช้เพิ่มข้อมูลใหม่ลงในตารางจริง  
Parameters: `id`, `schema`, `table`

**PUT `/introspect/sources/{id}/schemas/{schema}/tables/{table}/data/key/{key}`**  
ใช้แก้ไขข้อมูลในตารางจริงตาม key  
Parameters: `id`, `schema`, `table`, `key`, `key_field`

**DELETE `/introspect/sources/{id}/schemas/{schema}/tables/{table}/data/key/{key}`**  
ใช้ลบข้อมูลในตารางจริงตาม key  
Parameters: `id`, `schema`, `table`, `key`, `key_field`

## ระบบทะเบียนชุดข้อมูล (Dataset Catalog)

ใช้จัดการข้อมูลพื้นฐานของ dataset เช่น ชื่อ เจ้าของ domain tag และระดับความอ่อนไหว

**GET `/datasets`**  
ใช้ค้นหาและดึงรายการ dataset ทั้งหมดในระบบ  
Parameters: `search_word`, `domain`, `tag`, `has_pii`, `owner`, `page`, `per_page`, `sort_by`, `sort_order`

**GET `/datasets/{id}`**  
ใช้ดึงรายละเอียด dataset ตามรหัส dataset  
Parameters: `id`

**POST `/datasets`**  
ใช้สร้าง dataset ใหม่ หรือแก้ไข dataset ที่มีอยู่

**DELETE `/datasets/{id}`**  
ใช้ลบ dataset ตามรหัส dataset  
Parameters: `id`

## ระบบเวอร์ชันและสัญญาข้อมูล (Dataset Version and Contract)

ใช้จัดการ version ของ dataset รวมถึง schema, access policy และ policy สำหรับการให้บริการข้อมูล

**GET `/datasets/{id}/versions`**  
ใช้ดึงรายการ version ของ dataset  
Parameters: `id`, `page`, `per_page`, `source_id`, `search_word`, `status`

**GET `/datasets/{id}/versions/{version}`**  
ใช้ดึงรายละเอียด contract ของ dataset version  
Parameters: `id`, `version`

**POST `/datasets/{id}/versions`**  
ใช้สร้าง version ใหม่ให้กับ dataset  
Parameters: `id`

**PUT `/datasets/{id}/versions/{version}`**  
ใช้แก้ไขข้อมูล contract ของ dataset version  
Parameters: `id`, `version`

**PATCH `/datasets/{id}/versions/{version}`**  
ใช้แก้ไขสถานะของ dataset version เช่น `active`, `preview`, `deprecated`  
Parameters: `id`, `version`

## ระบบให้บริการข้อมูลชุดข้อมูล (Dataset Serving)

ใช้ให้บริการข้อมูลจาก dataset version ตามสิทธิ์และ policy ที่กำหนดไว้

**GET `/datasets/{id}/versions/{version}/data`**  
ใช้ดึงข้อมูลจาก dataset version ตามสิทธิ์และ policy ที่กำหนด  
Parameters: `id`, `version`, `page`, `per_page`, `view`, `where`, `where_logical_operator`, `sort_by`, `sort_order`

**GET `/datasets/{id}/versions/{version}/data/key/{key}`**  
ใช้ดึงข้อมูล 1 รายการจาก dataset version ตาม key  
Parameters: `id`, `version`, `key`, `view`

**POST `/datasets/{id}/versions/{version}/data`**  
ใช้เพิ่มข้อมูลใหม่ใน dataset version  
Parameters: `id`, `version`

**PUT `/datasets/{id}/versions/{version}/data/key/{key}`**  
ใช้แก้ไขข้อมูลใน dataset version ตาม key  
Parameters: `id`, `version`, `key`

**DELETE `/datasets/{id}/versions/{version}/data/key/{key}`**  
ใช้ลบข้อมูลใน dataset version ตาม key  
Parameters: `id`, `version`, `key`

## ระบบรายงานและนำเข้า/ส่งออกข้อมูล (Reporting and Export/Import)

ใช้สร้างงาน export/import ข้อมูล และจัดการ template สำหรับออกรายงาน PDF

**POST `/reporting/export/job`**  
ใช้สร้างงาน export ข้อมูล dataset เป็นไฟล์ เช่น CSV หรือ XLSX

**GET `/reporting/export/job/{job_id}`**  
ใช้ดึงสถานะและรายละเอียดงาน export dataset  
Parameters: `job_id`

**POST `/reporting/import/template`**  
ใช้สร้าง template สำหรับนำเข้าข้อมูล dataset

**POST `/reporting/import/job`**  
ใช้สร้างงาน import ข้อมูลจากไฟล์เข้าสู่ dataset

**GET `/reporting/import/job/{job_id}`**  
ใช้ดึงสถานะและรายละเอียดงาน import dataset  
Parameters: `job_id`

**POST `/reporting/templates/upload`**  
ใช้อัปโหลดไฟล์ PDF template สำหรับสร้างรายงาน

**POST `/reporting/templates/{reporting_template_id}/export/key/{key}`**  
ใช้สร้างงาน export รายงาน PDF จาก template และข้อมูลตาม key  
Parameters: `reporting_template_id`, `key`

**GET `/reporting/templates/export/job/{job_id}`**  
ใช้ดึงสถานะและรายละเอียดงาน export รายงาน PDF  
Parameters: `job_id`
