DELETE FROM permissions WHERE name LIKE 'attendance.%';
DROP TABLE IF EXISTS attendance_devices;
DROP TABLE IF EXISTS attendance_locations;
DROP TABLE IF EXISTS attendance_corrections;
DROP TABLE IF EXISTS attendance_punches;
DROP TABLE IF EXISTS attendance_records;
DROP TABLE IF EXISTS attendance_policies;
